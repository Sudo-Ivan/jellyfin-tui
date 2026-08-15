package tui

import (
	"syscall"
	"time"
	"unsafe"
)

func (t *Term) Listen(stop <-chan struct{}) <-chan Event {
	ch := make(chan Event, eventQueue)
	go t.readKeys(ch, stop)
	return ch
}

func (t *Term) readKeys(ch chan<- Event, stop <-chan struct{}) {
	hin, err := getStdHandle(stdInputHandle)
	if err != nil {
		return
	}
	recs := make([]inputRecord, winRecCount)
	for {
		select {
		case <-stop:
			return
		default:
		}
		var n uint32
		r1, _, e := procReadConsoleInput.Call(
			uintptr(hin),
			uintptr(unsafe.Pointer(&recs[0])), // #nosec G103
			uintptr(len(recs)),
			uintptr(unsafe.Pointer(&n)), // #nosec G103
		)
		if r1 == 0 {
			if e == syscall.Errno(0) {
				continue
			}
			return
		}
		if !t.emitRecords(ch, stop, recs[:n]) {
			return
		}
	}
}

func (t *Term) emitRecords(ch chan<- Event, stop <-chan struct{}, recs []inputRecord) bool {
	for i := range recs {
		ev, ok := t.winRecord(recs[i])
		if !ok {
			continue
		}
		ev.At = time.Now()
		select {
		case ch <- ev:
		case <-stop:
			return false
		}
	}
	return true
}

func (t *Term) winRecord(r inputRecord) (Event, bool) {
	switch r.EventType {
	case windowSizeEvent:
		w, h, err := t.readSize()
		if err != nil {
			return Event{}, false
		}
		return Event{Kind: KindResize, W: w, H: h}, true
	case keyEventType:
		var ke keyEventRecord
		copy((*[inputEventSize]byte)(unsafe.Pointer(&ke))[:], r.Event[:]) // #nosec G103
		if ke.KeyDown == 0 {
			return Event{}, false
		}
		return winKey(ke)
	case mouseEventType:
		var me mouseEventRecord
		copy((*[inputEventSize]byte)(unsafe.Pointer(&me))[:], r.Event[:]) // #nosec G103
		return winMouse(me)
	default:
		return Event{}, false
	}
}
