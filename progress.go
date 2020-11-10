package upload

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

type Event struct {
	fileName string
	uploaded int64
	total    int64
}

type Progress struct {
	Event    chan Event
	progress map[string]string
}

func (p *Progress) Start() {
	for e := range p.Event {
		if _, ok := p.progress[e.fileName]; ok {
			p.progress[e.fileName] = getProgress(e.uploaded, e.total)
			p.render()
		}
	}
}

func (p *Progress) render() {
	var progress bytes.Buffer
	progress.WriteByte('\r')
	for fileName, pro := range p.progress {
		progress.WriteString(fmt.Sprintf("%s:\n%s", fileName, pro))
	}
	progress.WriteByte('\n')

	fmt.Fprintf(os.Stdout, progress.String())
}

func NewProgress(files []string) *Progress {
	pMap := make(map[string]string)
	for _, file := range files {
		pMap[file] = ""
	}
	p := &Progress{Event: make(chan Event, len(files)), progress: pMap}
	return p
}

func getProgress(uploaded, total int64) string {
	var progress bytes.Buffer
	done := strings.Repeat("#", int(float64(uploaded)/float64(total)*100))
	todo := strings.Repeat("*", 100-int(float64(uploaded)/float64(total)*100))
	progress.WriteString(fmt.Sprintf("%s%s", done, todo))
	return progress.String()
}
