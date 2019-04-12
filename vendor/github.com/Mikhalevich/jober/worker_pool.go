package jober

type joberFinishNotifier interface {
	Jober
	addCallback(f WorkerFunc, callback func())
}

type WorkerPool struct {
	job   joberFinishNotifier
	count chan bool
}

func NewWorkerPool(j joberFinishNotifier, c int) *WorkerPool {
	return &WorkerPool{
		job:   j,
		count: make(chan bool, c),
	}
}

func (wp *WorkerPool) Add(f WorkerFunc) {
	wp.count <- true
	wp.job.addCallback(f, func() {
		<-wp.count
	})
}

func (wp *WorkerPool) Wait() {
	wp.job.Wait()
}

func (wp *WorkerPool) Get() ([]interface{}, []error) {
	return wp.job.Get()
}
