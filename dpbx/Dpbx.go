package dpbx

type Dpbx struct {
	dropbox    int
	empty      bool
	full       bool // redundant (always !empty) but reads better
	access     sync.Mutex
	untilEmpty *sync.Cond
	untilFull  *sync.Cond
}

func New() (*Dpbx) {
	d := &Dpbx{dropbox:0, empty:true, full:false}
	cve := sync.NewCond(&d.access)
	cvf := sync.NewCond(&d.access)
	d.untilEmpty = cve
	d.untilFull = cvf
	return d
}

func (d *dbpx) Send(value int) {
	defer d.access.Unlock()
	d.access.Lock()
	for d.full {
		d.untilEmpty.Wait()
	}
	d.dropbox = value
	d.empty = false
	d.full  = true
	d.untilFull.Signal()
}

func (d *Dpbx) Recv(value *int) {
	defer d.access.Unlock()
	d.access.Lock()
	for d.empty {
		d.untilFull.Wait()
	}
	(*value) = d.dropbox
	d.empty = true
	d.full  = false
	d.untilEmpty.Signal()
}
