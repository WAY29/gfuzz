package payloads

import (
	"io/ioutil"
	"time"
)

func init() {
	AddPayload("dirwalk", &PayloadDirwalk{})
	AddPayloadInfo("dirwalk", "Returns filename's recursively from a local directory.")
}

type PayloadDirwalk struct {
	channel chan interface{}
}

//
func (p *PayloadDirwalk) Channel() chan interface{} {
	return p.channel
}

func (p *PayloadDirwalk) New(s ...interface{}) error {
	p.channel = make(chan interface{})
	ch := p.channel
	go func() {
		defer func() {
			close(ch)
		}()
		startPath := s[0].(string)
		lenOfStartPath := len(startPath)
		filePaths := make(chan string, 1)
		filePaths <- startPath
		for {
			select {
			case filePath, ok := <-filePaths:
				if !ok {
					return
				}
				files, err := ioutil.ReadDir(filePath)
				if err != nil {
					return
				}
				for _, v := range files {
					fullPath := filePath + "/" + v.Name()
					if v.IsDir() {
						go func() { filePaths <- fullPath }()
					} else {
						ch <- fullPath[lenOfStartPath+1:]
					}
				}
			case <-time.After(1):
				return
			}
		}

	}()
	return nil
}
