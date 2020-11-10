package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Listen error %v", err)
	}

	// var bf *bufio.Reader
	addr := l.Addr().String()
	log.Printf("Server listen at %s", addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("Accept error %v", err)
		}

		go downloadHandler(conn)
	}
}

func downloadHandler(conn net.Conn) {
	defer conn.Close()

	newName := strconv.FormatInt(time.Now().Unix(), 10)

	line := ReadFirstLine(conn)

	s := strings.SplitN(string(line), " ", 2)
	flag, fileName := s[0], s[1]
	if flag != "UPLOAD_FILE" {
		return
	}

	newName = newName + "_" + fileName
	uploadPath := filepath.Join("upload_files", newName)

	dir := filepath.Dir(uploadPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Printf("Create %q error %v", "upload_files", err)
		return
	}

	f, err := os.Create(uploadPath)
	defer f.Close()
	if err != nil {
		log.Printf("Create %s error %v", uploadPath, err)
		return
	}

	_, err = io.Copy(f, conn)
	if err != nil {
		log.Printf("Copy error %v", err)
		return
	}

	err = f.Sync()
	if err != nil {
		log.Printf("Sync error %v", err)
		return
	}
}

func ReadFirstLine(conn net.Conn) string {
	var bf bytes.Buffer
	for {
		b := make([]byte, 1)
		_, err := conn.Read(b)
		if err != nil {
			panic("Read first line error.")
		}
		if b[0] == '\n' || b[0] == '\r' {
			break
		}

		bf.Write(b)
	}

	return bf.String()
}
