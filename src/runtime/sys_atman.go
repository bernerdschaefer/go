package runtime

import "unsafe"

const (
	_PAGESIZE = 0x1000
)

var (
	_atman_stack          [8 * _PAGESIZE]byte
	_atman_start_info     *xenStartInfo
	_atman_hypercall_page [_PAGESIZE]byte
)

func forceReachability() {
	_atman_hypercall_page[0] = 'a'
}

//go:nosplit
func getRandomData(r []byte) {
	forceReachability() // TODO: remove this
	extendRandom(r, 0)
}

// lock

const (
	active_spin     = 4
	active_spin_cnt = 30
)

func lock(l *mutex)   {}
func unlock(l *mutex) {}

func noteclear(n *note)                  {}
func notewakeup(n *note)                 {}
func notesleep(n *note)                  {}
func notetsleep(n *note, ns int64) bool  { return false }
func notetsleepg(n *note, ns int64) bool { return false }

// env

func gogetenv(key string) string { return "" }

var _cgo_setenv unsafe.Pointer   // pointer to C function
var _cgo_unsetenv unsafe.Pointer // pointer to C function

// mem

func sysAlloc(n uintptr, sysStat *uint64) unsafe.Pointer                    { return nil }
func sysFree(v unsafe.Pointer, n uintptr, sysStat *uint64)                  {}
func sysMap(v unsafe.Pointer, n uintptr, reserved bool, sysStat *uint64)    {}
func sysReserve(v unsafe.Pointer, n uintptr, reserved *bool) unsafe.Pointer { return nil }
func sysUnused(v unsafe.Pointer, n uintptr)                                 {}
func sysUsed(v unsafe.Pointer, n uintptr)                                   {}
func sysFault(v unsafe.Pointer, n uintptr)                                  {}

// os

func sigpanic() {}
func crash()    {}
func goenvs()   {}

func newosproc(mp *m, stk unsafe.Pointer) {}

func resetcpuprofiler(hz int32) {}

func minit()         {}
func unminit()       {}
func mpreinit(mp *m) {}
func msigsave(mp *m) {}

//go:nosplit
func osyield() {}

func osinit() {}

// signals

const _NSIG = 0

func initsig()                 {}
func sigdisable(uint32)        {}
func sigenable(uint32)         {}
func sigignore(uint32)         {}
func raisebadsignal(sig int32) {}

// net

func netpoll(block bool) *g { return nil }
func netpollinited() bool   { return false }

type xenStartInfo struct {
	Magic          [32]byte
	NrPages        uint32
	SharedInfoAddr uint32 // machine address of share info struct
	SIFFlags       uint32
	StoreMfn       uint32 // machine page number of shared page
	StoreEventchn  uint32
	Console        struct {
		Mfn      uint32 // machine page number of console page
		Eventchn uint32 // event channel
	}
	PageTableBase     uint32 // virtual address of page directory
	NrPageTableFrames uint32
	MfnList           uint32 // virtual address of page-frame list
	ModStart          uint32 // virtual address of pre-loaded module
	ModLen            uint32 // size (bytes) of pre-loaded module
	CmdLine           [1024]byte

	// The pfn range here covers both page table and p->m table frames
	FirstP2mPfn uint32 // 1st pfn forming initial P->M table
	NrP2mFrames uint32 // # of pgns forming initial P->M table
}
