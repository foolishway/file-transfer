package upload

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/k0kubun/go-ansi"
)

type Event struct {
	fileName string
	uploaded int64
	total    int64
}

type Progress struct {
	Event    chan Event
	progress map[string]string
	rendered bool
}

func (p *Progress) Start(done chan struct{}) {
	for e := range p.Event {
		if _, ok := p.progress[e.fileName]; ok {
			p.progress[e.fileName] = getProgress(e.uploaded, e.total)
			if p.rendered {
				p.clear()
			}
			p.render()
		}
	}
	done <- struct{}{}
}

func (p *Progress) render() {
	if !p.rendered {
		p.rendered = true
	}
	keys := make([]string, 0)
	for fileName := range p.progress {
		keys = append(keys, fileName)
	}
	sort.Strings(keys)

	var progress bytes.Buffer
	for _, fileName := range keys {
		if pro, ok := p.progress[fileName]; ok {
			progress.WriteString(fmt.Sprintf("%s:\n%s\n", fileName, pro))
		}
	}

	fmt.Fprintf(os.Stdout, progress.String())
}

func NewProgress(files []string) *Progress {
	pMap := make(map[string]string)
	for _, file := range files {
		pMap[file] = strings.Repeat("□", 100)
	}
	p := &Progress{Event: make(chan Event, len(files)), progress: pMap}
	return p
}

func getProgress(uploaded, total int64) string {
	var progress bytes.Buffer
	done := strings.Repeat("■", int(float64(uploaded)/float64(total)*100))
	todo := strings.Repeat("□", 100-int(float64(uploaded)/float64(total)*100))
	progress.WriteString(fmt.Sprintf("%s%s", done, todo))
	return progress.String()
}

func (p *Progress) clear() {
	for i := 0; i < len(p.progress)*2; i++ {
		ansi.CursorUp(1)
		ansi.EraseInLine(2)
	}
}
