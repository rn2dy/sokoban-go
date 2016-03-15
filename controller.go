package main

import (
	"github.com/nsf/termbox-go"
	"time"
)

const animationSpeed = 10 * time.Millisecond

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	g := NewGame()
	render(g)

	for {
		ev := <-eventQueue
		if ev.Type == termbox.EventKey {
			switch {
			case ev.Key == termbox.KeyArrowUp || ev.Ch == 'k':
				g.move(UP)
			case ev.Key == termbox.KeyArrowDown || ev.Ch == 'j':
				g.move(DOWN)
			case ev.Key == termbox.KeyArrowLeft || ev.Ch == 'h':
				g.move(LEFT)
			case ev.Key == termbox.KeyArrowRight || ev.Ch == 'l':
				g.move(RIGHT)
			case ev.Ch == 'n':
				g.nextLevel()
			case ev.Ch == 'p':
				g.prevLevel()
			case ev.Ch == 'r':
				g.reset()
			case ev.Ch == 'd':
				g.toggleDebug()
			case ev.Key == termbox.KeyEsc:
				return
			}
		}
		render(g)
		time.Sleep(animationSpeed)
	}
}
