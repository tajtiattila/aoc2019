package input

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func ClearCache() error {
	return os.RemoveAll(cachedir())
}

func daydata(n int) ([]byte, error) {
	if !IgnoreCache {
		b, err := ioutil.ReadFile(dayinputfn(n))
		if err == nil {
			return b, nil
		}
	}

	if err := downloaddaydata(n); err != nil {
		return nil, err
	}

	return ioutil.ReadFile(dayinputfn(n))
}

func cachedir() string {
	d, err := os.UserCacheDir()
	if err == nil {
		return filepath.Join(d, cachesubdir)
	}
	return ".input"
}

func dayinputfn(day int) string {
	return filepath.Join(cachedir(), fmt.Sprint(day))
}

var dataclient struct {
	once sync.Once

	c   *http.Client
	err error
}

func createFile(path string) (io.WriteCloser, error) {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return nil, err
		}
	}
	return os.Create(path)
}
