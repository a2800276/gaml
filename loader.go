package gaml

import (
	"fmt"
	"os"
	"strings"
)

// Loader, Entry are not particularly nice and custom tailored to the http handlers
// needs. Probably should be made private.

// The whole convolute is a victim of premature optimization ...

type Loader interface {
	Load(id interface{}) (root *node, err error)
}

type fileSystemLoader struct {
	baseDir string
}

func NewFileSystemLoader(dir string) (loader Loader, err error) {
	var f *os.File
	if f, err = os.Open(dir); err != nil {
		return
	}

	defer f.Close()

	var fi os.FileInfo
	if fi, err = f.Stat(); err != nil {
		return
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("%s: not a directory", fi.Name())
	}

	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	return &fileSystemLoader{dir}, nil
}

func (l *fileSystemLoader) Load(id_string interface{}) (root *node, err error) {
	// check
	id, ok := id_string.(string)
	if !ok {
		err = fmt.Errorf("id: %s is not a string", id)
		return
	}

	var file *os.File
	// check fs
	var path = l.baseDir + id
	if file, err = os.Open(path); err != nil {
		return
	}

	defer file.Close()

	parser := NewParser(file)
	// the parser inherits this loader to handle loading includes that
	// it may encounter.
	parser.IncludeLoader = l
	return parser.Parse()
}
