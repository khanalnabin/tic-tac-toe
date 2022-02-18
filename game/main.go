package game

import (
	"fmt"
	"os"
)

var game *Game = &Game{}

func Run() {
	err := game.Initialize()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to initialize")
		panic(err)
	}
	for game.Running {
		err := game.HandleEvents()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unable to Handle Event")
			panic(err)
		}
		err = game.Update()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unable to Update")
			panic(err)
		}
		err = game.Render()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unable to Render")
			panic(err)
		}
	}
	game.Clean()
}
