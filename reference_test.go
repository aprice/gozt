package gozt

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestReference(t *testing.T) {
	files, err := ioutil.ReadDir("testfiles/reference/")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("found %d reference files", len(files))
	for _, file := range files {
		t.Logf("found reference file: %s", file.Name())
		t.Skip("full reference implementation pending")
		name := strings.Split(file.Name(), ".")[0]
		t.Run(name, func(tt *testing.T) {
			refFile, err := os.Open("testfiles/reference/" + name + ".html")
			if err != nil {
				tt.Error(err)
				return
			}

			tplFile, err := os.Open("testfiles/templates/" + name + ".html")
			if err != nil {
				tt.Error(err)
				return
			}

			template, err := ReadTemplate(tplFile)
			if err != nil {
				tt.Error(err)
				return
			}
			modelFile, err := os.Open("testfiles/templates/" + name + ".json")
			if err != nil {
				tt.Error(err)
				return
			}
			var model interface{}
			json.NewDecoder(modelFile).Decode(&model)

			parser := New(template, model)
			err = parser.Parse()
			if err != nil {
				tt.Error(err)
				return
			}

			b := 0
			chunkSize := 512
			for {
				buffx := make([]byte, chunkSize)
				len, err1 := refFile.Read(buffx)

				buffa := make([]byte, chunkSize)
				_, err2 := parser.Read(buffa)

				if err1 != nil || err2 != nil {
					if err1 == io.EOF && err2 == io.EOF {

					} else if err1 == io.EOF {
						tt.Error("Result longer than expected.")
					} else if err2 == io.EOF {
						tt.Error("Result shorter than expected.")
					} else {
						tt.Error(err1, err2)
					}
					return
				} else if !bytes.Equal(buffx, buffa) {
					tt.Errorf("Byte mismatch at index %d:\n\tExpected: %s\n\tActual: %s", b, buffx, buffa)
					break
				} else {
					b += len
				}
			}
		})
	}
}
