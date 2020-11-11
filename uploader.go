package upload

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Uploader struct {
	FilePath   string
	Size       int64
	Wg         *sync.WaitGroup
	ServerAddr string
	Progress   *Progress
}

func (c *Uploader) Upload() {
	defer c.Wg.Done()

	f, err := os.Open(c.FilePath)
	defer f.Close()
	if err != nil {
		panic(fmt.Sprintf("Open file error %v", err))
	}

	finfo, err := os.Stat(c.FilePath)
	if err != nil {
		panic(fmt.Sprintf("Stat file error %v", err))
	}
	c.Size = finfo.Size()

	conn, err := net.Dial("tcp", c.ServerAddr)
	// defer conn.Close()
	if err != nil {
		panic(fmt.Sprintf("Dial %s error %v", c.ServerAddr, err))
	}

	fLine := fmt.Sprintf("UPLOAD_FILE %s\n", filepath.Base(c.FilePath))
	conn.Write([]byte(fLine))
	uw := &uploaderWapper{reader: f, size: c.Size, fileName: f.Name(), Progress: c.Progress}
	buf := make([]byte, 10)
	io.CopyBuffer(conn, uw, buf)
}

type uploaderWapper struct {
	reader   io.Reader
	fileName string
	size     int64
	uploaded int64
	Progress *Progress
}

func (w *uploaderWapper) Read(b []byte) (readed int, err error) {
	if readed, err = w.reader.Read(b); err == nil {
		time.Sleep(50 * time.Millisecond)
		w.uploaded += int64(readed)
		w.Progress.Event <- Event{fileName: w.fileName, uploaded: w.uploaded, total: w.size}
	}
	return
}
