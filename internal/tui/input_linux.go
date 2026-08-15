package tui

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (t *Term) Listen(stop <-chan struct{}) <-chan Event {
	ch := make(chan Event, eventQueue)
	go t.readKeys(ch, stop)
	go t.watchResize(ch, stop)
	return ch
}

func (t *Term) watchResize(ch chan<- Event, stop <-chan struct{}) {
	winch := make(chan os.Signal, 1)
	signal.Notify(winch, syscall.SIGWINCH)
	defer signal.Stop(winch)
	for {
		select {
		case <-stop:
			return
		case <-winch:
			w, h, err := t.readSize()
			if err != nil {
				continue
			}
			select {
			case ch <- Event{Kind: KindResize, W: w, H: h, At: time.Now()}:
			case <-stop:
				return
			}
		}
	}
}

func (t *Term) readKeys(ch chan<- Event, stop <-chan struct{}) {
	buf := make([]byte, keyReadBuf)
	var leftover []byte
	for {
		select {
		case <-stop:
			return
		default:
		}
		n, err := t.in.Read(buf)
		if n == 0 && err != nil {
			return
		}
		data := append(leftover, buf[:n]...)
		leftover = t.drainKeys(ch, stop, data, buf)
		if leftover == nil {
			return
		}
	}
}

func (t *Term) drainKeys(ch chan<- Event, stop <-chan struct{}, data, buf []byte) []byte {
	for len(data) > 0 {
		ev, rest, ok := decodeKey(data)
		if !ok {
			next, keep, ok := t.finishPartial(ch, stop, data, buf)
			if !ok {
				return nil
			}
			if keep != nil {
				return keep
			}
			data = next
			continue
		}
		ev.At = time.Now()
		select {
		case ch <- ev:
		case <-stop:
			return nil
		}
		data = rest
	}
	return []byte{}
}

func (t *Term) finishPartial(ch chan<- Event, stop <-chan struct{}, data, buf []byte) (next []byte, leftover []byte, ok bool) {
	if data[0] != asciiESC || len(data) >= escMaxPending {
		return data, append([]byte(nil), data...), true
	}
	_ = t.in.SetReadDeadline(time.Now().Add(escWait))
	n, err := t.in.Read(buf)
	_ = t.in.SetReadDeadline(time.Time{})
	if n > 0 {
		return append(data, buf[:n]...), nil, true
	}
	if err != nil && len(data) > 0 {
		ev := Event{Kind: KindKey, Key: KeyEsc, At: time.Now()}
		select {
		case ch <- ev:
		case <-stop:
			return data, nil, false
		}
		return data[1:], nil, true
	}
	return data, append([]byte(nil), data...), true
}
