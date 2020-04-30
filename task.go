package vortex

import "time"

var tw *TimeWheel

/**
 * init tw
 */
func TimeWheelInit() *TimeWheel {
	tw = New(time.Second, 180)
	tw.Start()
	return tw
}

