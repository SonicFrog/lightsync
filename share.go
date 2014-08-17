package main

import (
	"bytes"
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"fsnotify"
	"io"
	"io/ioutil"
	"os"
	"path"
)

//Can't import it normally because it's only used for driver
import _ "github.com/gwenn/gosqlite"

type Share struct {
	Name    string
	Clients map[string]*Client
	Path    string
	Watcher *fsnotify.Watcher

	Database *sql.DB
}

const (
	FileChunkSize int64 = 1024 ^ 2
)

func NewShare(name, path string) (s *Share, err error) {

	wat, err := fsnotify.NewWatcher()

	if err != nil {
		return
	}

	db, err := sql.Open("sqlite3", "temp.db")

	if err != nil {
		LogObj.Println("Could not open the database for share ", name, ": ", err)
		return
	}

	s = &Share{name, make(map[string]*Client), path, wat, db}

	s.Watch(path)

	return
}

func (s *Share) AddClient(client *Client) {
	_, contains := s.Clients[client.Name]

	if !contains {
		s.Clients[client.Name] = client
	} else {
		LogObj.Println("Tried to add already connected client ", client.Name)
	}
}

func (s *Share) RemoveClient(client *Client) {
	_, contains := s.Clients[client.Name]

	if !contains {
		LogObj.Println("Unregistered client for share tried to disconnect ", client.Name)
	} else {
		delete(s.Clients, client.Name)
	}
}

func (s *Share) NotifyClients(msg Message) {
	for _, c := range s.Clients {
		c.WriteMessage(msg)
	}
}

func (s *Share) CreateFile(file string) (err error) {
	file = path.Join(s.Path, file)

	_, err = os.Stat(file)

	if err != nil && !os.IsNotExist(err) {
		return
	}

	if os.IsNotExist(err) {
		_, err = os.Create(file)
		if err != nil {
			LogObj.Println("Error while creating ", file, ": ", err)
			return
		}

	}

	return nil
}

func (s *Share) CreateDir(dir string) error {
	dir = path.Join(s.Path, dir)

	stat, err := os.Stat(path.Clean(dir))

	if err == nil {
		if !stat.IsDir() {
			return errors.New("A file with the name " + dir + " already exists!\n")
		}
		return nil
	}

	err = os.Mkdir(dir, 0755)

	return err
}

func (s *Share) WriteChunk(file string, partnum int64, part []byte) (err error) {
	var fd *os.File

	file = path.Join(s.Path, file)

	fd, err = os.OpenFile(file, os.O_RDWR, os.ModeExclusive)

	defer fd.Close()

	if err != nil {
		LogObj.Println("Fatal error while opening ", file, " for update: ", err)
		return err
	}

	if int64(len(part)) != FileChunkSize {
		LogObj.Println("Invalid chunk size ", len(part), " in ", file, ". Last chunk?")
	}

	_, err = fd.Seek(partnum*FileChunkSize, 0)

	fd.Write(part)

	if err != nil {
		fmt.Printf("Error updating chunk %d in %s:", partnum, file)
		LogObj.Println(err)
		return err
	}

	return nil
}

func (s *Share) ReadChunk(file string, partnum int64) (chunk []byte, err error) {
	file = path.Join(s.Path, file)

	fd, err := os.OpenFile(file, os.O_RDWR, os.ModeExclusive)

	if err != nil {
		return
	}

	defer fd.Close()

	chunk = make([]byte, FileChunkSize)

	n, err := fd.ReadAt(chunk, FileChunkSize*partnum)

	if err != nil {
		return
	}

	if int64(n) != FileChunkSize {
		var b bytes.Buffer
		b.Write(chunk)
		b.Truncate(n)
		chunk = b.Bytes()
	}

	return
}

func (s *Share) Remove(object string) {
	object = path.Join(s.Path, object)

	os.Remove(object)
	//TODO: notify DB of removal
}

func (s *Share) FromOffline(dir string) error {
	if len(dir) == 0 {
		dirinfo, err := os.Stat(s.Path)
		if dirinfo.IsDir() && err == nil {
			s.FromOffline(s.Path)
		} else {
			return errors.New("Invalid share path: " + dir + "!")
		}
	} else {
		dirinfo, err := os.Stat(dir)

		if !dirinfo.IsDir() || err != nil {
			return nil
		}

		files, err := ioutil.ReadDir(dir)

		if err != nil {
			LogObj.Println("Error reading directory ", dir, ": ", err)
			return err //Resuming recursing in other subdirs
		}

		for _, f := range files {
			mod, err := s.CheckFileShallow(f.Name())
			if err != nil {
				LogObj.Print("Could not read ", f.Name(), ": ", err)
				continue
			}

			if mod {
				s.CheckFileDeep(f.Name())
			}
		}
	}
	return nil
}

func (s *Share) CheckFileShallow(path string) (modified bool, err error) {
	stat, err := os.Stat(path)

	if err != nil {
		return
	}

	mtime, err := s.StoredModTime(path)

	modified = (mtime == stat.ModTime().UTC().Unix())
	//Time stored as UTC to avoid problems with timezones

	return
}

func (s *Share) CheckFileDeep(path string) (modified bool, err error) {
	file, err := os.Open(path)

	if err != nil {
		return
	}

	defer file.Close()

	hasher := sha1.New()

	_, err = io.Copy(hasher, file)

	if err != nil {
		return
	}

	lastHash, err := s.StoredHash(path)
	currentHash := hasher.Sum(nil)

	if len(lastHash) != len(currentHash) {
		//Consider file modified if stored hash is invalid ?
		return true, nil
	}

	modified = (bytes.Compare(currentHash, lastHash) == 0)

	return
}

func (s *Share) StoredModTime(path string) (mtime int64, err error) {
	//TODO: Retrieve stored info from database
	return
}

func (s *Share) StoredHash(path string) (hash []byte, err error) {
	return
}

func (s *Share) Events() chan fsnotify.Event {
	return s.Watcher.Events
}

func (s *Share) Errors() chan error {
	return s.Watcher.Errors
}

func (s *Share) Close() {
	s.Watcher.Close()
}

func (s *Share) Watch(dir string) error {
	finfo, err := ioutil.ReadDir(dir)

	if err != nil {
		return err
	}

	for _, f := range finfo {
		if f.IsDir() {
			err := s.Watch(dir + f.Name())
			if err != nil {
				return err
			}
		} else {
			s.Watcher.Add(dir + f.Name())
		}
	}
	return nil
}
