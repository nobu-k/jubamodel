package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"unsafe"
)

type Model struct {
	Path       string             `json:"path"`
	Header     Header             `json:"header"`
	SystemData []SystemDataHeader `json:"-"`
	UserData   []UserDataHeader   `json:"-"`
}

type Header struct {
	Magic          string `json:"Magic"`
	FormatVersion  uint64 `json:"format_version"`
	JubatusVersion string `json:"jubatus_version"`
	CRC32          uint32 `json:"crc32"`
	SystemDataSize uint64 `json:"system_data_size"`
	UserDataSize   uint64 `json:"user_data_size"`

	Raw          []byte       `json:"-"`
	BinaryHeader BinaryHeader `json:"-"`
}

type BinaryHeader struct {
	Magic                        [8]byte
	FormatVersion                uint64
	Major, Minor, Maintenance    uint32
	CRC32                        uint32
	SystemDataSize, UserDataSize uint64
}

type SystemDataHeader struct {
	Version   string                 `json:"version"`
	Timestamp int64                  `json:"timestamp"`
	Type      string                 `json:"type"`
	ID        string                 `json:"id"`
	Config    map[string]interface{} `json:"config"`
}

type UserDataHeader struct {
	Version uint64 `json:"version"`
}

func Info(paths []string) ([]*Model, error) {
	res := []*Model{}
	for _, p := range paths {
		m, err := info(p)
		if err != nil { // TODO: Add option to ignore all files which are not jubatus models.
			fmt.Fprintln(os.Stderr, "Cannot read a model file:", err)
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func info(path string) (*Model, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	m := &Model{
		Path: absPath,
		Header: Header{
			Raw: make([]byte, unsafe.Sizeof(BinaryHeader{})),
		},
	}
	if n, err := f.Read(m.Header.Raw); err != nil {
		return nil, err
	} else if n < len(m.Header.Raw) {
		return nil, fmt.Errorf("the file is too small")
	}

	bh := &m.Header.BinaryHeader
	if err := binary.Read(bytes.NewReader(m.Header.Raw), binary.BigEndian, bh); err != nil {
		return nil, err
	}

	header := &m.Header
	header.Magic = string(bh.Magic[:])
	header.FormatVersion = bh.FormatVersion
	header.JubatusVersion = fmt.Sprint(bh.Major, ".", bh.Minor, ".", bh.Maintenance)
	header.CRC32 = bh.CRC32
	header.SystemDataSize = bh.SystemDataSize
	header.UserDataSize = bh.UserDataSize

	// TODO: read containers
	return m, nil
}

func ExecInfo(args []string) {
	if len(args) == 0 {
		InfoUsage()
	}

	ms, err := Info(args)
	if err != nil {
		os.Exit(1)
	}

	var js []byte
	if len(ms) == 1 {
		js, err = json.Marshal(ms[0])
	} else { // including 0
		js, err = json.Marshal(ms)
	}
	if err == nil {
		fmt.Print(string(js))
	}
}

func InfoUsage() {
	fmt.Println("Usage jubamodel info file [files...]")
	fmt.Println()
	fmt.Println("Show information of model files")
	os.Exit(1)
}
