package bar

import (
	"bytes"
	"fmt"
	"runtime"
	"syscall"
	"unsafe"

	spin "github.com/tj/go-spin"
)

const (
	TIOCGWINSZ     = 0x5413
	TIOCGWINSZ_OSX = 1074295912
)

type window struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func TerminalWidth() (int, error) {
	w := new(window)
	tio := syscall.TIOCGWINSZ
	if runtime.GOOS == "darwin" {
		tio = TIOCGWINSZ_OSX
	}
	res, _, err := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(tio),
		uintptr(unsafe.Pointer(w)),
	)
	if int(res) == -1 {
		return 0, err
	}
	return int(w.Col), nil
}

var sp *spin.Spinner = spin.New()

func Icon(used, taskNum, total int) string {
	s := bytes.NewBufferString("")
	s.WriteString("[")
	for index := 0; index < total-10; index++ {
		if index <= used*(total-10)/taskNum {
			s.WriteString("=")
		} else {
			s.WriteString(" ")
		}
	}
	s.WriteString(fmt.Sprintf("]%s", sp.Next()))
	// s.WriteString()
	return s.String()
}
