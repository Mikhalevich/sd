package jober

import "sync"

type FirstError struct {
	job
	once sync.Once
}

func NewFirstError() *FirstError {
	return &FirstError{
		job: *newJob(),
	}
}

func (fr *FirstError) processError() {
	err, ok := <-fr.errorChan
	if ok {
		fr.dataErrors = append(fr.dataErrors, err)
		fr.cancel()
	}
	fr.errorFinishFlag <- true
}

func (fr *FirstError) processFunc() {
	fr.waitGroup.Wait()
	close(fr.dataChan)
	close(fr.errorChan)
	<-fr.dataFinishFlag
}

func (fr *FirstError) startProcess(p processer) {
	fr.once.Do(func() {
		fr.job.startProcess(fr)
		go fr.processFunc()
	})
}

func (fr *FirstError) Add(f WorkerFunc) {
	fr.startProcess(fr)
	fr.job.Add(f)
}

func (fr *FirstError) addCallback(f WorkerFunc, callback func()) {
	fr.startProcess(fr)
	fr.job.addCallback(f, callback)
}

func (fr *FirstError) Wait() {
	<-fr.errorFinishFlag
}
