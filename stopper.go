package main

import (
	ps "github.com/mitchellh/go-ps"
	"strings"
	"time"
)

func (s *server) stopper() {

	for range time.Tick(time.Second * 5){
		processList, err := ps.Processes()
		if err != nil {
			continue
		}

		select {
		case <-s.ctx.Done():
			return
		default:
			for x := range processList {
				var process ps.Process
				process = processList[x]
				if strings.Contains(strings.ToLower(process.Executable()), strings.ToLower("AthenaAIO")) {
					s.cancel()
				}

			}
		}

	}
}
