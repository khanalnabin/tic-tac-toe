package game

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Player = int8
type Cell = int8
type GameState = int8
type Mode = int8

const (
	PlayerX Player = iota
	PlayerO
)

const (
	Running GameState = iota
	XWon
	OWon
	Draw
)

const (
	Empty Cell = iota
	X
	O
)

const (
	Single Mode = iota
	Multi
)

type GameGrid struct {
	Array      [3][3]int8
	Mode       Mode
	Turn       Player
	State      GameState
	CellCount  int32
	Width      int32
	Height     int32
	CellHeight int32
	CellWidth  int32
	PosX       int32
	PosY       int32
	EndIndex   [2]int32
}

type Game struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Grid     *GameGrid
	Mouse    *Mouse
	Height   int32
	Width    int32
	Running  bool
	Selected bool
}

type Mouse struct {
	X       int32
	Y       int32
	Clicked bool
}
