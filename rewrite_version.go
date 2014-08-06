package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"strconv"
	"strings"
)

func Rewrite(path string, version string) error {
	m, err := info(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot get information of the model:", err)
		return err
	}

	bh := m.Header.BinaryHeader
	bh.Major, bh.Minor, bh.Maintenance, err = parseVersion(version)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid Jubatus version:", err)
		return err
	}

	bh.CRC32, err = calcNewCRC32(path, &bh)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot compute CRC32 of the new model:", err)
		return err
	}

	f, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot write to the file:", err)
		return err
	}
	defer f.Close()

	if err := binary.Write(f, binary.BigEndian, &bh); err != nil {
		fmt.Fprintln(os.Stderr, "cannot write a new header to the file:", err)
		return err
	}
	return nil
}

func parseVersion(version string) (major uint32, minor uint32, maintenance uint32, err error) {
	vs := strings.Split(version, ".")
	if len(vs) != 3 {
		err = fmt.Errorf("the version doesn't contain three numbers")
		return
	}

	vns := make([]uint32, 3)
	for i, v := range vs {
		n, e := strconv.ParseUint(v, 10, 32)
		if e != nil {
			err = e
			return
		}
		vns[i] = uint32(n)
	}
	major, minor, maintenance = vns[0], vns[1], vns[2]
	return
}

func calcNewCRC32(path string, newHeader *BinaryHeader) (uint32, error) {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.BigEndian, newHeader)
	if err != nil {
		return 0, err
	}

	header := buf.Bytes()
	crc := crc32.Checksum(header[:28], crc32.IEEETable) // before crc32
	crc = crc32.Update(crc, crc32.IEEETable, header[32:])

	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	if _, err := f.Seek(int64(len(header)), os.SEEK_SET); err != nil {
		return 0, err
	}

	data := make([]byte, 4*1024*1024)
	for {
		n, err := f.Read(data)
		crc = crc32.Update(crc, crc32.IEEETable, data[:n])
		if err != nil {
			if err != io.EOF {
				return 0, err
			}
			break
		}
	}
	return crc, nil
}

func ExecRewriteVersion(args []string) {
	if len(args) != 2 {
		RewriteVersionUsage()
	}

	err := Rewrite(args[0], args[1])
	if err != nil {
		os.Exit(1)
	}
}

func RewriteVersionUsage() {
	fmt.Println("Usage: jubamodel rewrite-version file new-version")
	fmt.Println()
	fmt.Println("Rewrite a Jubatus version of the given model file")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("    $ jubamodel rewrite-version /path/to/classifier.model 0.5.8")
	os.Exit(1)
}
