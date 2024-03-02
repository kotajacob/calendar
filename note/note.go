package note

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Path returns the filepath for a given note.
//
// Environment variables, such as $HOME may be used in the dir and will be
// expanded appropriately.
func Path(t time.Time, dir string) string {
	return filepath.Join(os.ExpandEnv(dir), t.Format("2006-01-02")) + ".md"
}

// Exists stats a note file for a given time.
// If the files Exists, but is empty it is counted as not existing.
func Exists(t time.Time, dir string) bool {
	path := Path(t, dir)
	stat, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false
		} else {
			log.Println(err)
		}
	}
	if stat.IsDir() {
		return false
	}
	if stat.Size() == 0 {
		return false
	}
	return true
}

// Load reads a note file for a given time.
//
// If the file is missing it is simply treated as an empty file. All other
// errors will return the error string itself (which is meant to be displayed
// to the user).
func Load(t time.Time, dir string) string {
	path := Path(t, dir)
	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		data = []byte(err.Error())
	}
	return string(data)
}
