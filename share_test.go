package main

import (
	"bytes"
	crand "crypto/rand"
	mrand "math/rand"
	"os"
	"testing"
)

const (
	FileChunkNum int = 100
	TestNumber       = 1000
	ShareDir         = "/tmp/lightsync/"
	TestFile         = "test.001"
)

var Sh *Share

var Success bool = true

func InitShare(t *testing.T) {
	t.Log("Creating tempdir for share...")
	err := os.Mkdir(ShareDir, os.ModeDir|os.ModePerm)

	defer func() { Success = (err == nil) }()

	if err != nil && !os.IsExist(err) {
		t.Error("Unable to create a temp dir for share testing! ", err)
		return
	}

	Sh, err = NewShare("test", ShareDir)

	if err != nil {
		t.Error("Unable to create test share in ", ShareDir)
		return
	}
}

func TestCreateShare(t *testing.T) {
	if Sh == nil {
		InitShare(t)
	}
	var err error

	defer func() { Success = (err == nil) }()

	if Sh.Name != "test" || Sh.Path != ShareDir {
		t.Error("Share name mismatch! Got: ", Sh.Name, ". Expected: test")
		return
	}
}

func TestCreateFile(t *testing.T) {
	if !Success || Sh == nil {
		t.Skip("Previous test failed can't run this one!")
	}

	err := Sh.CreateFile(TestFile)

	defer func() { Success = (err == nil) }()

	if err != nil {
		t.Error("Error while creating test file: ", err)
		return
	}
}

func TestReadChunk(t *testing.T) {
	var err error

	if !Success || Sh == nil {
		t.Skip("Previous test failed can't run this one!")
	}

	defer func() { Success = (err == nil) }()

	var chunks = make([][]byte, FileChunkNum)
	var chunk = make([]byte, FileChunkSize)

	f, err := os.OpenFile(ShareDir+TestFile, os.O_RDWR, 0666)

	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < FileChunkNum; i++ {
		chunks[i] = make([]byte, FileChunkSize)
		crand.Read(chunks[i])
		_, err := f.Write(chunks[i])

		if err != nil {
			t.Error(err)
			return
		}
	}

	for i := 0; i < FileChunkNum; i++ {
		_, err = f.ReadAt(chunk, int64(i)*int64(FileChunkSize))
		if err != nil {
			t.Error(err)
			return
		}

		if bytes.Compare(chunk, chunks[i]) != 0 {
			t.Error("Read a different chunk!")
			return
		}
	}
}

func TestWriteChunk(t *testing.T) {
	if !Success || Sh == nil {
		t.Skip("Previous test failed can't test this!")
	}

	var err error

	defer func() { Success = (err == nil) }()

	for i := 0; i < TestNumber; i++ {
		var partnum int64 = int64(mrand.Intn(FileChunkNum))

		data := make([]byte, FileChunkSize)

		err = Sh.WriteChunk(TestFile, partnum, data)

		if err != nil {
			t.Error("Could not write chunk: ", err)
			return
		}

		rddata, err := Sh.ReadChunk(TestFile, partnum)

		if err != nil {
			t.Error("Could not read chunk: ", err)
			return
		}

		if bytes.Compare(data, rddata) != 0 {
			t.Error("Read chunk is different from written chunk!")
			return
		}
	}
}
