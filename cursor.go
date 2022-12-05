package cursor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var Esc = "\x1b"

func escape(format string, args ...interface{}) string {
	return fmt.Sprintf("%s%s", Esc, fmt.Sprintf(format, args...))
}

func Show() string {
	return escape("[?25h")
}

// Hide returns ANSI escape sequence to hide the cursor
func Hide() string {
	return escape("[?25l")
}
func MoveTo(line, col int) string {
	return escape("[%d;%dH", line, col)
}

func ClearEntireScreen() string {
	return escape("[2J")
}

func setConsoleColors() error {
	console := windows.Stdout
	var consoleMode uint32
	windows.GetConsoleMode(console, &consoleMode)
	consoleMode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	return windows.SetConsoleMode(console, consoleMode)
}

func Color(colorcode int) string {
	setConsoleColors() // init of terminal
	code := strconv.Itoa(colorcode)
	color := "\u001b[38;5;" + code + "m" // acii code used by terminals to print colors
	return fmt.Sprintf(color)

}

func Clear() {
	cmd := exec.Command("cmd", "/c", "cls")
	if runtime.GOOS == "windows" {
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Println("Your OS is not supported")
		fmt.Scanf(" ")
		os.Exit(3)
	}
}

func SetConsoleTitle(title string) (int, error) {
	handle, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return 0, err
	}
	defer syscall.FreeLibrary(handle)
	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")
	if err != nil {
		return 0, err
	}
	r, _, err := syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	return int(r), err
}
