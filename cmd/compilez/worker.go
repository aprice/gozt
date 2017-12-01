package main

import (
	"bytes"
	"io"
)

type worker struct {
	queue chan job
	buf   bytes.Buffer
}

func (w *worker) run() {
	for job := range w.queue {
		w.process(job)
	}
}

func (w *worker) process(job job) {

}

type job struct {
	in    io.Reader
	out   io.Writer
	model interface{}
}
