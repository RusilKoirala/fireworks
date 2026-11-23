package main

import "fmt"


const (
	maxTapeSeconds = 20 * 60
	maxTapeEvents = 3000
)


type launchCue struct {
	At float64 `json:"at"`
	Kind kind `json:"kind"`
	X float64 `json:"x"`
	Target float64 `json:"target"`
}


type recording struct {
	Version int `json:"version"`
	Duration float64 `json:"duration"`
	Events []launchCue `json:"events"`
}


func (t recording) validate() error {
	if t.Version != 1 {
		return fmt.Errorf("unsupported recording version")
	}

	if !finite(t.Duration) || t.Duration <= 0 || t.Duration > maxTapeEvents {
		return  fmt.Errorf("invalid recording duration")
	}
	if len(t.Events) == 0 || len(t.Events) > maxTapeEvents {
		return fmt.Errorf("recording needs 1 to %d launches", maxTapeEvents)
	}

	previous:= 0.0

	for i,cue := range t.Events {
		if !finite(cue.At) || cue.At < previous || cue.At > t.Duration {
			return  fmt.Errorf("invalid time at launch %d", i+ 1)
		}

		if cue.Kind < classic || cue.Kind >= numKinds {
			return fmt.Errorf("invalid shape at launch %d", i +1)
		}

		if !finite(cue.X) || cue.X < 0 || cue.X > 1 {
			return fmt.Errorf("invalid position at launch %d ", i+1)
		}

		if !finite(cue.Target) || (cue.Target != -1 && (cue.Target< 0 || cue.Target>1)) {
			return  fmt.Errorf("invalid target at launch %d", i+1)
		}

		previous = cue.At
	}

	return  nil
}

func (m *model) toggleRecording() {
	if m.studio.recorder.active {
		m.finishRecording()
		return
	}

	if m.studio.replay.active {
		m.tell("Stop replay with t before recording")
		return
	}

	m.studio.tape = recording{
		Version: 1,
	}

	m.studio.recorder =recorderState{
		active: true,
		started: m.clock,
	}

	m.tell("recording launches; press r to finish")
}


func (m *model) finishRecording() {
	if !m.studio.recorder.active {
		return
	}

	duration := max(dt,m.clock-m.studio.recorder.started)
	m.studio.tape.Duration = min(duration, float64(maxTapeSeconds))
	m.studio.recorder.active = false

	m.tell(fmt.Sprintf(
		"recorded %d launches; t replay; w save",
		len(m.studio.tape.Events),
	))


}

func (m *model) captureLaunch(l launchMsg) {
	m.studio.launched++ 

	if !m.studio.recorder.active {
		return
	}

	at := m.clock - m.studio.recorder.started 

	if at >= maxTapeEvents || len(m.studio.tape.Events) >= maxTapeEvents {
		m.finishRecording()
		m.tell("recording stopped: size or duration limit reached");
		return 
	}

	target := -1.0

	if l.yAim >= 0 {
		target = clamp(l.yAim/float64(m.skyH()),0,1)
	}

	m.studio.tape.Events = append(m.studio.tape.Events, launchCue{
		At: at,
		Kind: l.k,
		X: clamp(l.xf,0,1),
		Target: target,
	})



}

func (m *model) launchCue(cue launchCue) {
	target := -1.0

	if cue.Target >= 0{ 
		target = cue.Target*float64(m.skyH())
	}

	m.launches(launchMsg{
		k: cue.Kind,
		xf: cue.X,
		yAim: target,
	})
}


func (m *model) toggleReplay() {
	if m.studio.replay.active {
		m.studio.replay.active =false
		m.tell("replay stopped")
		return
	}

	m.finishRecording()

	if err := m.studio.tape.validate(); err != nil {
		m.tell("cannot replay"+ err.Error());
		return
	}

	m.setAuto(false)
	m.studio.show.active = false 
	m.clearSky()
	m.studio.replay = replayState{active : true}

	m.tell("replay started; press t to stop")
}

func (m *model) advanceReplay() {
	p := &m.studio.replay

	if !p.active {
		return
	}

	tape := &m.studio.tape

	for p.next < len(tape.Events) && tape.Events[p.next].At <= p.elapsed {
		m.launchCue(tape.Events[p.next])
		p.next++
	}

	if p.elapsed >= tape.Duration && p.next == len(tape.Events) {
		m.active = false 
		m.tell("replay finished")
	}

	p.elapsed += dt
}