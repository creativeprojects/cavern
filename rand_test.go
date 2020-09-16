package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandomInt(t *testing.T) {
	timeout := time.After(3 * time.Second)
	done := make(chan bool)

	go func() {
		// test first and last value are picked
		first := 70
		last := 200
		firstPick, lastPick := false, false
		for !firstPick || !lastPick {
			pick := randomInt(first, last+1)
			if pick == first {
				firstPick = true
				continue
			}
			if pick == last {
				lastPick = true
			}
		}
		assert.True(t, firstPick)
		assert.True(t, lastPick)

		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("test didn't finish in time")
	case <-done:
	}
}
