package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Mikhalevich/argparser"
	"github.com/Mikhalevich/downloader"
	"github.com/Mikhalevich/pbw"
)

type Params struct {
	Method        string `json:"method,omitempty"`
	MaxWorkers    int64  `json:"max_workers,omitempty"`
	ChunkSize     int64  `json:"chunk_size,omitempty"`
	UseFileSystem bool   `json:"use_filesystem,omitempty"`
}

func NewParams() *Params {
	return &Params{
		Method:        "GET",
		MaxWorkers:    downloader.DefaultMaxWorkers,
		ChunkSize:     downloader.DefaultChunkSize,
		UseFileSystem: true,
	}
}

func getURL(p *argparser.Parser) (string, error) {
	arguments := p.Arguments()
	if len(arguments) <= 0 {
		return "", errors.New("No url for download specified")
	}

	urlString := arguments[0]
	uri, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	if uri.Scheme == "" {
		urlString = "http://" + urlString
	}

	return urlString, nil
}

func loadParams() (*Params, string, error) {
	p := argparser.NewParser()
	basicParams := NewParams()
	params, err, _ := p.Parse(basicParams)
	if err != nil {
		return nil, "", err
	}

	uri, err := getURL(p)
	if err != nil {
		return nil, "", err
	}

	return params.(*Params), uri, err
}

func doDownload(url string, params *Params) error {
	startTime := time.Now()

	task := downloader.NewChunkedTask()
	task.Method = params.Method
	task.MaxDownloaders = params.MaxWorkers
	task.ChunkSize = params.ChunkSize
	if params.UseFileSystem {
		chunksPath := fmt.Sprintf("%s_%s", url[strings.LastIndex(url, "/")+1:], "chunks")
		task.CS = downloader.NewFileStorer(chunksPath)
		defer os.RemoveAll(chunksPath)
	}
	task.Notifier = make(chan int64, task.MaxDownloaders*3)

	pbw.Show(task.Notifier)

	fileName, err := task.Download(url)
	if err != nil {
		return err
	}

	fmt.Printf("Downloaded sucessfully into %s, time elapsed: %s\n", fileName, time.Now().Sub(startTime))
	return nil
}

func main() {
	params, uri, err := loadParams()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = doDownload(uri, params)
	if err != nil {
		fmt.Println(err)
	}
}
