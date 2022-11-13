package main

import (
	"errors"
	"os"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

type Serial struct {
	f *os.File
}

func openSerialPort(deviceName string) (*Serial, error) {
	f, err := os.OpenFile(deviceName, unix.O_RDWR|unix.O_NOCTTY|unix.O_NONBLOCK, 0666)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil && f != nil {
			f.Close()
		}
	}()

	var bauds uint32 = unix.B9600

	fd := f.Fd()

	t := unix.Termios{
		Iflag:  unix.IGNPAR,
		Cflag:  unix.CREAD | unix.CLOCAL | bauds | unix.CS8, /*8 data bits*/
		Ispeed: bauds,
		Ospeed: bauds,
	}

	if _, _, errno := unix.Syscall6(
		unix.SYS_IOCTL,
		uintptr(fd),
		uintptr(unix.TCSETS),
		uintptr(unsafe.Pointer(&t)),
		0,
		0,
		0,
	); errno != 0 {
		return nil, errno
	}

	if err = unix.SetNonblock(int(fd), false); err != nil {
		return nil, err
	}

	return &Serial{f: f}, nil
}

func (p *Serial) MakeMeasurement() (Measurement, error) {
	//https://cdn-reichelt.de/documents/datenblatt/C150/MH-Z19C-PC_DATENBLATT.pdf
	_, err := p.f.Write([]byte("\xff\x01\x86\x00\x00\x00\x00\x00\x79"))
	if err != nil {
		return Measurement{}, err
	}

	buf := make([]byte, 64)
	n, err := p.f.Read(buf)
	if err != nil {
		return Measurement{}, err
	}

	if n >= 9 && buf[0] == 0xff && buf[1] == 0x86 && buf[8] == checksum(buf) {
		co2 := uint32(buf[2])*256 + uint32(buf[3])
		if co2 < 400 || co2 > 5000 {
			return Measurement{}, errors.New("sensor value out of range")
		}
		return Measurement{
			Co2:         uint32(buf[2])*256 + uint32(buf[3]),
			Temperature: uint32(buf[4]) - 40,
			Timestamp:   time.Now().Unix(),
		}, nil
	} else {
		return Measurement{}, errors.New("sensor returned invalid response")
	}
}

func checksum(buffer []byte) byte {
	var val byte = 0
	for i := 0; i < 8; i++ {
		val += buffer[i]
	}
	val = 0xff - val
	//val += 1
	return val
}
