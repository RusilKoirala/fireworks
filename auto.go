package main

import (
	"math/rand"
	"time"
)

func autoLauncher(out chan<- launchMsg, stop <-chan struct{}) {
	for {
		delay := time.Duration(250+rand.Intn(950)) * time.Millisecond
		select {
		case <-stop:
			return
		case <-time.After(delay):
		}

		n := 1
		if rand.Float64() < 0.18 { // finale-style salvo
			n = 3 + rand.Intn(4)
		}
		for i := 0; i < n; i++ {
			msg := launchMsg{k: kind(rand.Intn(int(numKinds))), xf: 0.1 + rand.Float64()*0.8, yAim: -1}
			select {
			case out <- msg:
			case <-stop:
				return
			}
			if n > 1 {
				time.Sleep(time.Duration(60+rand.Intn(120)) * time.Millisecond)
			}
		}
	}
}

func (m *model) setAuto(on bool) {
	if on == m.auto {
		return
	}
	m.auto = on
	if on {
		m.stopAuto = make(chan struct{})
		go autoLauncher(m.launches, m.stopAuto)
	} else {
		close(m.stopAuto)
	}
}
