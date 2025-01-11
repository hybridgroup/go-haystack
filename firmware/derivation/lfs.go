package main

import (
	"encoding/binary"
	"io"
	"machine"
	"os"

	"tinygo.org/x/tinyfs/littlefs"
)

var lfs = littlefs.New(machine.Flash)

func getIndex() uint64 {
	f, err := lfs.Open("/haystack")
	if err != nil {
		return 0
	}
	defer f.Close()

	var buf [8]byte
	_, err = io.ReadFull(f, buf[:])
	must("read index from file", err)
	return binary.LittleEndian.Uint64(buf[:])
}

func writeIndex(i uint64) error {
	f, err := lfs.OpenFile("/haystack", os.O_CREATE)
	if err != nil {
		return err
	}
	defer f.Close()

	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], i)
	_, err = f.Write(buf[:])
	return err
}

func init() {
	config := littlefs.Config{
		CacheSize:     64,
		LookaheadSize: 32,
		BlockCycles:   512,
	}
	lfs.Configure(&config)
	if err := lfs.Mount(); err != nil {
		must("format littlefs", lfs.Format())
		must("mount littlefs", lfs.Mount())
	}
	println("littlefs mounted")
}
