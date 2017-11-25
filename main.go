package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Mikhalevich/argparser"
	"github.com/Mikhalevich/downloader"
)

type Params struct {
	Method        string `json:"method, omitempty"`
	MaxWorkers    int64  `json:"max_workers, omitempty"`
	ChunkSize     int64  `json:"chunk_size, omitempty"`
	UseFileSystem bool   `json:"use_filesystem, omitempty"`
}

func NewParams() *Params {
	return &Params{
		Method:     "GET",
		MaxWorkers: downloader.DefaultMaxWorkers,
		ChunkSize:  downloader.DefaultChunkSize,
	}
}

func loadParams() (*Params, error) {
	basicParams := NewParams()
	params, err, _ := argparser.Parse(basicParams)
	return params.(*Params), err
}

func getUrl() (string, error) {
	if argparser.NArg() <= 0 {
		return "", errors.New("No url for download specified")
	}

	urlString := argparser.Arg(0)
	uri, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	if uri.Scheme == "" {
		urlString = "http://" + urlString
	}

	return urlString, nil
}

func doDownload(url string, params *Params) error {
	startTime := time.Now()

	task := downloader.NewChunkedTask()
	task.Method = params.Method
	task.MaxDownloaders = params.MaxWorkers
	task.ChunkSize = params.ChunkSize
	if params.UseFileSystem {
		task.CS = downloader.NewFileStorer(fmt.Sprintf("%s_%s", url[strings.LastIndex(url, "/")+1:], "chunks"))
	}

	if err := task.Download(url); err != nil {
		return err
	}

	fmt.Printf("Downloaded sucessfully, time elapsed: %s\n", time.Now().Sub(startTime))
	return nil
}

func main() {
	params, err := loadParams()
	if err != nil {
		fmt.Println(err)
		return
	}

	uri, err := getUrl()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = doDownload(uri, params)
	if err != nil {
		fmt.Println(err)
	}
}
