// Copyright 2013 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program will download files from textfiles.com.
//
// Usage:
//
//   % download -u <url> -d <dest-dir> -w <workers>
//
// Example of url: http://www.textfiles.com/programming/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"sync"
)

var link = regexp.MustCompile(`<a href="([\w-]+\.txt)">[\w-]+\.txt</a>`)
var url, dstdir string
var workers uint

func init() {
	flag.StringVar(&url, "u", "", "URL to download files from")
	flag.StringVar(&dstdir, "d", "", "Destination directory (where to save files)")
	flag.UintVar(&workers, "w", 2, "Number of workers")
}

func extract(url string, files chan<- string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	links := link.FindAllSubmatch(bytes.ToLower(b), -1)
	for _, l := range links {
		files <- string(l[1]) // HL
	}
	close(files) // HL
	return nil
}

func download(files <-chan string, wg *sync.WaitGroup) {
	defer wg.Done() // HL
	var caminho string
	for file := range files {
		caminho = url + file
		resp, err := http.Get(caminho)
		if err != nil {
			fmt.Println("Arquivo", caminho, "com erro de abertura")
			continue
		}
		open, err := os.Create(path.Join(dstdir, file))
		if err != nil {
			fmt.Println("Falhou ao abrir arquivo", err)
			continue
		}
		io.Copy(open, resp.Body)
		open.Close()
		resp.Body.Close()
	}
}

func main() {
	flag.Parse() // HL
	os.MkdirAll(dstdir, 0755)
	var wg sync.WaitGroup // HL
	if workers < 1 {
		workers = 2
	}
	files := make(chan string, workers) // HL
	for i := uint(0); i < workers; i++ {
		wg.Add(1)
		go download(files, &wg) // HL
	}
	err := extract(url, files)
	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()
}
