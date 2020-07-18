package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

const Format = "2006-01-02 15:04:05.000"

func main() {
	fmt.Println("Hello World")

	t := NewDynamicTicker(10*time.Second, func(tickTime time.Time) {
		log.Println("业务代码，执行回调 ticked = ", tickTime.Format(Format))
	})

	for {
		time.Sleep(30 * time.Second)

		n := time.Duration(rand.Int31n(20)+10) * time.Second
		log.Println("业务代码，调整间隔", n)
		t.ChangeInterval(n)
	}
}

// DynamicTicker 定义动态间隔的滴答器结构.
type DynamicTicker struct {
	Interval       time.Duration
	IntervalChange chan time.Duration
	Fn             func(tickTime time.Time)
}

// NewDynamicTicker 创建一个新的动态滴答器.
func NewDynamicTicker(interval time.Duration, fn func(time.Time)) *DynamicTicker {
	d := &DynamicTicker{
		Interval:       interval,
		IntervalChange: make(chan time.Duration, 1),
		Fn:             fn,
	}

	go d.start()

	return d
}

// ChangeInterval 调整滴答器的滴答时间间隔.
func (d *DynamicTicker) ChangeInterval(newInterval time.Duration) {
	d.IntervalChange <- newInterval
}

// start 开始周期性运行任务.
func (d *DynamicTicker) start() {
	ticker := time.NewTicker(d.Interval)
	defer ticker.Stop()

	log.Println("滴答要开始干活了，初始间隔为", d.Interval)

	for {
		select {
		case t, ok := <-ticker.C:
			if ok {
				log.Println("滴答，时间到", t.Format(Format))
				go d.Fn(t)
			} else {
				log.Println("管道歇菜了")
			}
		case ic := <-d.IntervalChange:
			if ic != d.Interval {
				log.Println("收到，滴答间隔调整为", ic)
				ticker.Stop()

				if ic > 0 {
					d.Interval = ic
					ticker = time.NewTicker(ic)
				}
			} else {
				log.Println("收到，滴答间隔还是", d.Interval)
			}
		}
	}
}

/*
https://play.golang.org/p/XEjWBKhBKly

Hello World
2009/11/10 23:00:00 滴答要开始干活了，初始间隔为 10s
2009/11/10 23:00:10 滴答，时间到 2009-11-10 23:00:10.000
2009/11/10 23:00:10 业务代码，执行回调 ticked =  2009-11-10 23:00:10.000
2009/11/10 23:00:20 滴答，时间到 2009-11-10 23:00:20.000
2009/11/10 23:00:20 业务代码，执行回调 ticked =  2009-11-10 23:00:20.000
2009/11/10 23:00:30 业务代码，调整间隔 11s
2009/11/10 23:00:30 滴答，时间到 2009-11-10 23:00:30.000
2009/11/10 23:00:30 收到，滴答间隔调整为 11s
2009/11/10 23:00:30 业务代码，执行回调 ticked =  2009-11-10 23:00:30.000
2009/11/10 23:00:41 滴答，时间到 2009-11-10 23:00:41.000
2009/11/10 23:00:41 业务代码，执行回调 ticked =  2009-11-10 23:00:41.000
2009/11/10 23:00:52 滴答，时间到 2009-11-10 23:00:52.000
2009/11/10 23:00:52 业务代码，执行回调 ticked =  2009-11-10 23:00:52.000
2009/11/10 23:01:00 业务代码，调整间隔 17s
2009/11/10 23:01:00 收到，滴答间隔调整为 17s
2009/11/10 23:01:17 滴答，时间到 2009-11-10 23:01:17.000
2009/11/10 23:01:17 业务代码，执行回调 ticked =  2009-11-10 23:01:17.000
2009/11/10 23:01:30 业务代码，调整间隔 17s
2009/11/10 23:01:30 收到，滴答间隔还是 17s
2009/11/10 23:01:34 滴答，时间到 2009-11-10 23:01:34.000
2009/11/10 23:01:34 业务代码，执行回调 ticked =  2009-11-10 23:01:34.000
2009/11/10 23:01:51 滴答，时间到 2009-11-10 23:01:51.000
2009/11/10 23:01:51 业务代码，执行回调 ticked =  2009-11-10 23:01:51.000
2009/11/10 23:02:00 业务代码，调整间隔 29s
2009/11/10 23:02:00 收到，滴答间隔调整为 29s
2009/11/10 23:02:29 滴答，时间到 2009-11-10 23:02:29.000
2009/11/10 23:02:29 业务代码，执行回调 ticked =  2009-11-10 23:02:29.000
2009/11/10 23:02:30 业务代码，调整间隔 11s
2009/11/10 23:02:30 收到，滴答间隔调整为 11s
2009/11/10 23:02:41 滴答，时间到 2009-11-10 23:02:41.000
2009/11/10 23:02:41 业务代码，执行回调 ticked =  2009-11-10 23:02:41.000
2009/11/10 23:02:52 滴答，时间到 2009-11-10 23:02:52.000
2009/11/10 23:02:52 业务代码，执行回调 ticked =  2009-11-10 23:02:52.000
2009/11/10 23:03:00 业务代码，调整间隔 28s
2009/11/10 23:03:00 收到，滴答间隔调整为 28s
2009/11/10 23:03:28 滴答，时间到 2009-11-10 23:03:28.000
2009/11/10 23:03:28 业务代码，执行回调 ticked =  2009-11-10 23:03:28.000
*/