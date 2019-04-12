package jober

type All struct {
	job
}

func NewAll() *All {
	return &All{
		job: *newJob(),
	}
}

func (a *All) Add(f WorkerFunc) {
	a.startProcess(a)
	a.job.Add(f)
}

func (a *All) addCallback(f WorkerFunc, callback func()) {
	a.startProcess(a)
	a.job.addCallback(f, callback)
}
