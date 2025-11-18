package hyrinx

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

var osCom = map[string]string{
	"windows": "cmd",
	"linux":   "xdg-open",
	"darwin":  "open",
}

// Launches the file at the specified path with the given args.
//   - path: The absolute path to the file
//   - opts: Launch args
func Launch(path string, opts ...string) {
	var cmd *exec.Cmd
	if isWindowsExecutable(path) {
		cmd = cmdWindowsExecutable(path, opts...)
	} else {
		com, err := createCommand(path, opts...)
		if err != nil {
			fmt.Printf("Failed to create command: %s\n", err)
			return
		}

		cmd = exec.Command(osCom[runtime.GOOS], com...)
	}

	if err := cmd.Run(); err != nil {
		log.Print(err)
	}
}

// Creates an execute command based on the current OS
//   - path: The absolute path to the file
//   - opts: Launch args
//
// Returns: os-specific command, err
func createCommand(path string, opts ...string) ([]string, error) {
	var args []string
	switch runtime.GOOS {
	case "windows":
		args = []string{"/C", "start", "", path}
	case "linux":
		args = []string{path}
	case "darwin":
		args = []string{path}
	default:
		fmt.Printf("Unsupported operating system: %s\n", runtime.GOOS)
		return []string{}, errors.New("unsupported operating system")
	}
	args = append(args, opts...)

	return args, nil
}

// Creates a command to launch windows executables directly
func cmdWindowsExecutable(path string, opts ...string) *exec.Cmd {
	cmd := exec.Command(path, opts...)
	return cmd
}
