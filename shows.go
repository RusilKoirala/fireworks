package main

import "sort"

type showScript struct {
	name string
	duration float64
	events []launchCue
}

type showPlayer struct {
	active bool 
	script showScript
	elapsed float64
	next int 
}


var showNames = []string{
	"Symphony",
	"Heartbeats",
	"Storm",
}

func buildShow(index int) showScript {
	script := showScript{
		name: showNames[index],
		duration: 32,
	}

	add := func (at float64,shape kind, x float64)  {
		script.events = append(script.events, launchCue{
			At: at,
			Kind: shape,
			X: x,
			Target: -1,
		})
	}

	switch index {
	case 0:
		for i:= 0; i<8;i++ {
			at:= float64(i) * 1.4
			x:= 0.18 + float64(i%4)*0.21

			add(at,ring,x)
			add(at+0.35,doubleRing, 1-x)

		}

		for i:= 0;i<9;i++ {
			at := 12 + float64(i)* 1.1

			add(at, spiral, 0.28)
			add(at+0.25.starburst, 0.72)
		}
	case 1:

		for i := 0; i< 10; i++ {
			at := float64(i) * 2.1

			add(at,heart,0.35)
			add(at+0.35, heart,0.65)
			add(at, willow, 0.5)
		}

		for i := 0;i<10;i++ {
			at := 22 + float64(i) * 0.65
			add(at,starburst, 0.2+float64(i)* 0.2)
		}
	case 2:
		for i:= 0;i<45;i++ {
			at := float64(i) * 0.55
			x:= 0.12 + float64((i*7)%17)/16*0.76

			add(at,kind(i%int(numKinds)),x)

			if i%3 == 0 {
				add(at + 0.12, palm, 1-x)
			}
		}
	}

	for wave := 0; wave < 5 ; wave++ {
		for position := 0; position < 5; position++ {
			at := 26 + float64(wave)* 0.75 + float64(position) * 0.08
			x:= 0.14 + float64(position) * 0.18

			shape := kind((wave+position) % int(numKinds))

			add(at,shape,x)
		}

	}

	sort.SliceStable(script.events, func(i,j int)bool {
		return script.events[i].At < script.events[j].At
	})

	return  script
}


func (m *model) toggleShow() {
	if m.studio.show.active {
		m.studio.show.active = false;
		m.tell("show stopped")
		return
	}

	m.setAuto(false)
	m.studio.replay.active = false 
	m.clearSky()

	script := buildShow(m.studio.selected)

	m.studio.show = showPlayer{
		active:  true,
		script: script,
	}

	m.tell("playing "+ script.name+ "; g stops the show")
}


func (m *model) advanceShow() {
	p := &m.studio.show

	if !p.active {
		return
	}

	for p.next < len(p.script.events) && p.script.events[p.next].At <= p.elapsed {
		m.launchCue(p.script.events[p.next])
		p.next++
	}

	if p.elapsed >= p.script.duration {
		p.active = false 
		m.tell(p.script.name + " finished")
	}

	p.elapsed += dt
}