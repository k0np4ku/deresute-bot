package main

import (
	"fmt"
	"image"
	"time"

	"github.com/go-vgo/robotgo"
)

var grooveButtons = []string{
	"groove_entrance",
	"groove_live",
	"groove_confirm",
	"groove_start",
	"groove_continue",
	"groove_continue2",
}

var liveOkButtons = []string{
	"live_ok",
	"live_ok2",
}

func doClick(coord *image.Point) {
	if coord == nil {
		return
	}

	if err := searchImage("gError_OK"); err != nil {
		robotgo.MoveClick(err.X, err.Y, "right")
		time.Sleep(1 * time.Second)
	}
	robotgo.MoveClick(coord.X, coord.Y, "right")
	time.Sleep(1 * time.Second)
}

func isPopHidden() bool {
	flags := make([]bool, 5)
	iterateBooleans(flags, func(v bool, i int) {
		flags[i] = searchImage(fmt.Sprintf("pop%d", i+1)) == nil
	})
	for _, flag := range flags {
		if !flag {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println("Deresute-bot v0.1.0")
	for {
		if isPopHidden() {
			if searchImage("live_pause") != nil {
				if err := searchImage("gError_OK"); err != nil {
					robotgo.MoveMouseSmooth(err.X, err.Y)
					robotgo.MouseClick("right")
				}
				time.Sleep(3 * time.Second)
				continue
			}

			if coord := searchImage("gameStart"); coord != nil {
				doClick(coord)
				continue
			}

			if searchImage("home_active") != nil {
				doClick(searchImage("live_entrance"))
				continue
			}

			if coord := searchImage("live_success"); coord != nil {
				var okCoord *image.Point
				for {
					if okCoord == nil {
						doClick(coord)
					} else {
						break
					}

					iterateStrings(liveOkButtons, func(v string, _ int) {
						if coord := searchImage(v); okCoord == nil {
							okCoord = coord
						}
					})
				}
				doClick(okCoord)
				continue
			}

			iterateStrings(grooveButtons, func(v string, _ int) {
				doClick(searchImage(v))
			})
			time.Sleep(1 * time.Second)
		} else {
			if searchImage("live_paused") != nil {
				doClick(searchImage("live_continue"))
			} else if searchImage("announcement") != nil {
				doClick(searchImage("announcement_close"))
			} else if searchImage("difficulty_select") != nil {
				doClick(searchImage("difficulty_target"))
				doClick(searchImage("difficulty_confirm"))
			} else if searchImage("multiplier_select") != nil {
				doClick(searchImage("multiplier_target"))
				doClick(searchImage("multiplier_confirm"))
			} else if searchImage("stamina") != nil || searchImage("stamina2") != nil {
				doClick(searchImage("stamina_refill"))
				doClick(searchImage("stamina_OK"))
			} else if searchImage("staminaAfter") != nil {
				doClick(searchImage("staminaAfter_close"))
			}
		}
	}
}
