package main

import (
	"flag"
	"io/ioutil"
	"log"
	"sync"

	upload "github.com/foolishway/go-multiupload"
)

func main() {
	serverAddr := flag.String("s", "", "Server addr")
	flag.Parse()

	if *serverAddr == "" {
		flag.Usage()
		panic("Server address is required.")
	}

	var files []string
	if flag.NArg() == 0 {
		// panic("Files to upload is requred.")
		finfo, _ := ioutil.ReadDir("./testData")
		for _, fi := range finfo {
			if files == nil {
				files = make([]string, 0)
			}
			files = append(files, "./testData/"+fi.Name())
		}
	} else {
		files = flag.Args()
	}

	done := make(chan struct{}, 0)
	//progress
	progress := upload.NewProgress(files)
	go progress.Start(done)

	var wg *sync.WaitGroup = &sync.WaitGroup{}
	for _, filePath := range files {
		wg.Add(1)
		go func(filePath string) {
			c := upload.Uploader{FilePath: filePath, Wg: wg, ServerAddr: *serverAddr, Progress: progress}
			c.Upload()
		}(filePath)
	}
	wg.Wait()
	close(progress.Event)
	<-done
	log.Println("All files upload completed.")
}
