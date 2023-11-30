package israc

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/anmitsu/go-shlex"
	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

type WrappedReader struct {
	closed       bool
	actualReader io.Reader
}

func (wr *WrappedReader) Close() error {
	wr.closed = true
	return nil
}

func (wr *WrappedReader) Read(p []byte) (n int, err error) {
	if wr.closed {
		return 0, io.EOF
	}
	return wr.actualReader.Read(p)
}

func Handle(s ssh.Session) {
	workdir := "/"
	tty := term.NewTerminal(s, "iSRAC> ")

	resolvePath := func(path string) string {
		if strings.HasPrefix(path, "/") {
			return path
		} else {
			return filepath.Join(workdir, path)
		}
	}

	var commands = map[string]*cobra.Command{
		"ls": {
			Use:   "ls",
			Short: "List files",
			Run: func(cmd *cobra.Command, args []string) {
				path := workdir
				if len(args) > 0 {
					path = resolvePath(args[0])
				}

				dir, err := os.ReadDir(path)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}
				for _, entry := range dir {
					cmd.Println(entry.Name())
				}
			},
		},
		"pwd": {
			Use:   "pwd",
			Short: "Print working directory",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Println(workdir)
			},
		},
		"cd": {
			Use:   "cd",
			Short: "Change working directory",
			Run: func(cmd *cobra.Command, args []string) {
				if len(args) == 0 {
					workdir = "/"
					return
				}
				newPath := resolvePath(args[0])
				stat, err := os.Stat(newPath)
				if err != nil {
					cmd.PrintErrln(args[0]+":", err)
					return
				}
				if !stat.IsDir() {
					cmd.PrintErrln(args[0] + ": is a file")
					return
				}

				workdir = newPath
			},
		},
		"exec": {
			Use:   "exec",
			Short: "Execute external command",
			Run: func(cmd *cobra.Command, args []string) {
				if len(args) == 0 {
					return
				}
				command := exec.Command(args[0], args[1:]...)
				command.Dir = workdir
				ptyReq, winCh, isPty := s.Pty()
				if isPty {
					command.Env = append(command.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))
					f, err := pty.Start(command)
					if err != nil {
						panic(err)
					}
					go func() {
						for win := range winCh {
							setWinsize(f, win.Width, win.Height)
						}
					}()
					stdin := &WrappedReader{actualReader: s}
					go func() {
						io.Copy(f, stdin) // stdin
					}()
					io.Copy(s, f) // stdout
					stdin.Close()
					command.Wait()
				} else {
					cmd.PrintErrln("No PTY requested.")
				}
			},
		},
		"fetch": {
			Use:   "fetch <url> <filename>",
			Short: "Download a file over HTTP(S) and write it somewhere",
			Args:  cobra.ExactArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				url := args[0]
				filename := resolvePath(args[1])

				out, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0o755)
				if err != nil {
					cmd.PrintErrln("open output file:", err)
					return
				}
				defer out.Close()

				res, err := http.Get(url)
				if err != nil {
					cmd.PrintErrln("http get:", err)
					return
				}
				defer res.Body.Close()

				_, err = io.Copy(out, res.Body)
				if err != nil {
					cmd.PrintErrln("write response:", err)
				}
			},
		},
		"cat": {
			Use:   "cat [...path]",
			Short: "Concatenate files",
			Run: func(cmd *cobra.Command, args []string) {
				for _, arg := range args {
					file, err := os.Open(resolvePath(arg))
					if err != nil {
						cmd.PrintErrf("open %s: %s\n", arg, err.Error())
						continue
					}
					defer file.Close()
					io.Copy(s, file)
				}
			},
		},
		"mv": {
			Use:   "mv <src> <dst>",
			Short: "Move a single file",
			Args:  cobra.ExactArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				err := os.Rename(resolvePath(args[0]), resolvePath(args[1]))
				if err != nil {
					cmd.PrintErrln(err)
				}
			},
		},
		"cp": {
			Use:   "cp <src> <dst>",
			Short: "Copy a single file",
			Args:  cobra.ExactArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				src, err := os.Open(resolvePath(args[0]))
				if err != nil {
					cmd.PrintErrln("open source:", err)
					return
				}
				defer src.Close()

				dst, err := os.Open(resolvePath(args[1]))
				if err != nil {
					cmd.PrintErrln("open destination:", err)
					return
				}
				defer dst.Close()

				io.Copy(dst, src)
			},
		},
		"rm": {
			Use:   "rm <path>",
			Short: "Remove a single file or empty directory",
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				err := os.Remove(resolvePath(args[0]))
				if err != nil {
					cmd.PrintErrln(err)
				}
			},
		},
		"mkdir": {
			Use:   "mkdir <path>",
			Short: "Make directory",
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				err := os.MkdirAll(resolvePath(args[0]), 0o755)
				if err != nil {
					cmd.PrintErrln(err)
				}
			},
		},
	}

	// make the command more raw
	commands["exec"].SetHelpFunc(commands["exec"].Run)

	var rootCmd = &cobra.Command{
		Use: "israc",
		Run: func(cmd *cobra.Command, args []string) { cmd.Println("woof!") },
	}
	for _, command := range commands {
		rootCmd.AddCommand(command)
	}
	rootCmd.SetOut(s)
	rootCmd.SetErr(s)

	for {
		line, err := tty.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				io.WriteString(s, "cya!\n")
				break
			}
			io.WriteString(s, fmt.Sprintf("cannot read line: %s", err.Error()))
			break
		}
		splat, err := shlex.Split(line, true)
		if err != nil {
			io.WriteString(s, fmt.Sprintf("cannot split line: %s", err.Error()))
			continue
		}
		rootCmd.SetArgs(splat)
		rootCmd.Execute()
	}
}
