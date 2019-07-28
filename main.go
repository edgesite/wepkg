package wepkg

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

type Pack struct {
	FileCount int
	Files     []*File
}

type File struct {
	Filename string
	Offset   int
	Size     int
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func readStr(reader io.Reader) (int, string) {
	n, size := readUint32LE(reader)

	buf := make([]byte, size)
	_, err := io.ReadFull(reader, buf)
	must(err)

	return n + size, string(buf)
}

func readUint32LE(r io.Reader) (int, int) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(r, buf)
	must(err)
	return 4, int(binary.LittleEndian.Uint32(buf))
}

func Unpack(filename string, verbose bool) int {
	fp, err := os.Open(filename)
	must(err)

	pos, header := readStr(fp)
	if strings.Compare(header, "PKGV0001") != 0 {
		fmt.Println("Unrecognized format:", header)
		return 1
	}

	n, count := readUint32LE(fp)
	pos += n

	pkg := &Pack{FileCount: count}
	for i := 0; i < pkg.FileCount; i++ {
		n, filename := readStr(fp)
		pos += n

		n, offset := readUint32LE(fp)
		pos += n

		n, size := readUint32LE(fp)
		pos += n

		file := &File{
			Filename: filename,
			Offset:   offset,
			Size:     size,
		}
		if verbose {
			fmt.Printf("%s(%d)\n", file.Filename, file.Size)
		}
		pkg.Files = append(pkg.Files, file)
	}

	unpackDir := "unpacked"
	os.Mkdir(unpackDir, os.ModeDir)
	for i := 0; i < pkg.FileCount; i++ {
		file := pkg.Files[i]
		subfilename := path.Join("unpacked", file.Filename)
		fmt.Println(subfilename)
		if strings.Index("../", subfilename) > -1 || strings.Index("..\\", subfilename) > -1 {
			fmt.Println("unsafe filename, refuse to extract")
			return -1
		}
		must(os.MkdirAll(path.Dir(subfilename), os.ModeDir))
		fp.Seek(int64(file.Offset+pos), 0)
		must(createAndCopy(subfilename, fp, int64(file.Size)))
	}
	fmt.Println("Create project.json")

	pjFp, err := os.Create(path.Join(unpackDir, "project.json"))
	must(err)
	_, err = pjFp.WriteString(`{
	"file" : "scene.json",
	"general" : 
	{
		"properties" : 
		{
			"schemecolor" : 
			{
				"order" : 0,
				"text" : "ui_browse_properties_scheme_color",
				"type" : "color",
				"value" : "0 0 0"
			}
		}
	},
	"title" : "Unpacked Project"
}`)
	must(err)

	return 0
}

func createAndCopy(filename string, r io.Reader, size int64) error {
	dst, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.CopyN(dst, r, size)
	return err
}
