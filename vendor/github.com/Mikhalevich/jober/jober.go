package jober

import (
	"sync"
)

type WorkerFunc func() (interface{}, error)

type Jober interface {
	Add(f WorkerFunc)
	Wait()
	Get() ([]interface{}, []error)
}

type processer interface {
	processData()
	processError()
}

type job struct {
	waitGroup       sync.WaitGroup
	once            sync.Once
	data            []interface{}
	dataChan        chan interface{}
	dataFinishFlag  chan bool
	dataErrors      []error
	errorChan       chan error
	errorFinishFlag chan bool
	cancelationChan chan bool
}

func newJob() *job {
	return &job{
		data:            make([]interface{}, 0),
		dataChan:        make(chan interface{}),
		dataFinishFlag:  make(chan bool),
		dataErrors:      make([]error, 0),
		errorChan:       make(chan error),
		errorFinishFlag: make(chan bool),
		cancelationChan: make(chan bool),
	}
}

func (j *job) cancel() {
	close(j.cancelationChan)
}

func (j *job) processData() {
	for d := range j.dataChan {
		j.data = append(j.data, d)
	}
	j.dataFinishFlag <- true
}

func (j *job) processError() {
	for err := range j.errorChan {
		j.dataErrors = append(j.dataErrors, err)
	}
	j.errorFinishFlag <- true
}

func (j *job) startProcess(p processer) {
	j.once.Do(func() {
		go p.processData()
		go p.processError()
	})
}

func (j *job) Wait() {
	j.waitGroup.Wait()
	close(j.dataChan)
	close(j.errorChan)
	<-j.dataFinishFlag
	<-j.errorFinishFlag
}

func (j *job) Get() ([]interface{}, []error) {
	return j.data, j.dataErrors
}

func (j *job) Add(f WorkerFunc) {
	j.waitGroup.Add(1)
	go func() {
		defer j.waitGroup.Done()
		d, err := f()
		if err != nil {
			select {
			case j.errorChan <- err:
			case <-j.cancelationChan:
			}
			return
		}
		select {
		case j.dataChan <- d:
		case <-j.cancelationChan:
		}
	}()
}

func (j *job) addCallback(f WorkerFunc, callback func()) {
	j.waitGroup.Add(1)
	go func() {
		defer callback()
		defer j.waitGroup.Done()
		d, err := f()
		if err != nil {
			select {
			case j.errorChan <- err:
			case <-j.cancelationChan:
			}
			return
		}
		select {
		case j.dataChan <- d:
		case <-j.cancelationChan:
		}
	}()
}
