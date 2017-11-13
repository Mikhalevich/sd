package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"time"

	"github.com/Mikhalevich/downloader"
)

var (
	cConfig = "sd_config.json"
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

func loadParams(configFile string) (*Params, error) {
	params := NewParams()

	file, err := os.Open(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return params, nil
		}
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, params)
	if err != nil {
		return nil, err
	}

	return params, nil
}

func getUrl() (string, error) {
	flag.Parse()

	if flag.NArg() <= 0 {
		return "", errors.New("No url for download specified")
	}

	urlString := flag.Arg(0)
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
		task.CS = downloader.NewFileStorer(fmt.Sprintf("%s_%s", url, "chunks"))
	}

	if err := task.Download(url); err != nil {
		return err
	}

	fmt.Printf("Downloaded sucessfully, time elapsed: %s\n", time.Now().Sub(startTime))
	return nil
}

func main() {
	uri, err := getUrl()
	if err != nil {
		fmt.Println(err)
		return
	}

	params, err := loadParams(cConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = doDownload(uri, params)
	if err != nil {
		fmt.Println(err)
	}
}
