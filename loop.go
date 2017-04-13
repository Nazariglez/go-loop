/**
 * Created by nazarigonzalez on 10/10/16.
 */

package loop

import (
	"time"
	"sync"
)

type Loop struct {
	sync.Mutex
	fps                  float64
	nanoFps              time.Duration
	isRunning            bool
	ticker               *time.Ticker
	last, lastTime, time int64
	Tick                 func (delta float64)
}

func NewLoop(fps float64) Loop {
	loop := Loop{
		Tick: func (delta float64) {},
	}
	loop.SetFPS(fps)
	return loop
}



func (loop *Loop) SetFPS(fps float64) {
	loop.Lock()
	restart := false
	if loop.isRunning {
		loop.stop()
		restart = true
	}

	loop.fps = fps
	loop.nanoFps = time.Duration((1/fps)*1e9) * time.Nanosecond

	if restart {
		loop.start()
	}
	loop.Unlock()
}

func (loop *Loop) start() {
	if loop.isRunning {
		return
	}

	loop.last = time.Now().UnixNano()
	loop.isRunning = true
	loop.ticker = time.NewTicker(loop.nanoFps)

	go func() {
		var now, delta int64

		for _ = range loop.ticker.C {
			loop.Lock()
			now = time.Now().UnixNano()
			loop.time += (now - loop.last)
			delta = loop.time - loop.lastTime
			loop.lastTime = loop.time
			loop.last = now
			loop.Unlock()
			loop.Tick(float64(delta) / 1e9)
		}
	}()
}

func (loop *Loop) stop() {
	if !loop.isRunning {
		return
	}

	loop.ticker.Stop()
	loop.isRunning = false
}

func (loop *Loop) Start() {
	loop.Lock()
	loop.start()
	loop.Unlock()
}

func (loop *Loop) Stop() {
	loop.Lock()
	loop.stop()
	loop.Unlock()
}

func (loop *Loop) IsRunning() bool {
	loop.Lock()
	defer loop.Unlock()
	return loop.isRunning
}