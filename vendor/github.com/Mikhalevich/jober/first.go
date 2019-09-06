package jober

import "sync"

type First struct {
	job
	once sync.Once
}

func NewFirst() *First {
	return &First{
		job: *newJob(),
	}
}

func (f *First) processData() {
	d, ok := <-f.dataChan
	if ok {
		f.data = append(f.data, d)
		f.cancel()
	}
	f.dataFinishFlag <- true
}

func (f *First) processFunc() {
	f.waitGroup.Wait()
	close(f.dataChan)
	close(f.errorChan)
	<-f.errorFinishFlag
}

func (f *First) startProcess(p processer) {
	f.once.Do(func() {
		f.job.startProcess(f)
		go f.processFunc()
	})
}

func (f *First) Add(fn WorkerFunc) {
	f.startProcess(f)
	f.job.Add(fn)
}

func (f *First) addCallback(fn WorkerFunc, callback func()) {
	f.startProcess(f)
	f.job.addCallback(fn, callback)
}

func (f *First) Wait() {
	<-f.dataFinishFlag
}
