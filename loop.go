/**
 * Created by nazarigonzalez on 10/10/16.
 */

package loop

import (
  "time"
)

type Loop struct {
  fps                  float64
  nanoFps              time.Duration
  IsRunning            bool
  ticker               *time.Ticker
  last, lastTime, time int64
  Tick                 chan float64
}

func NewLoop(fps float64) Loop {
  loop := Loop{
    Tick: make(chan float64),
  }
  loop.SetFPS(fps)
  return loop
}

func (loop *Loop) SetFPS(fps float64) {
  loop.fps = fps
  loop.nanoFps = time.Duration((1/fps)*1e9)*time.Nanosecond
}

func (loop *Loop) Start() {
  if loop.IsRunning {
    return
  }

  loop.IsRunning = true
  loop.ticker = time.NewTicker(loop.nanoFps)

  go func() {
    var now, delta int64

    for _ = range loop.ticker.C {
      now = time.Now().UnixNano()
      loop.time += (now-loop.last)/1000
      delta = loop.time -loop.lastTime
      loop.lastTime = loop.time
      loop.last = now
      loop.Tick <- float64(delta)/1e6
    }
  }()
}

func (loop *Loop) Stop() {
  if !loop.IsRunning {
    return
  }

  loop.ticker.Stop()
  loop.IsRunning = false
}