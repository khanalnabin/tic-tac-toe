package game

import (
	"math"
	"time"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var fontPath string = "/usr/share/fonts/ttf-mononoki/mononoki-Regular.ttf"

func (game *Game) Initialize() (err error) {
	if err != nil {
		return
	}
	game.Height = 600
	game.Width = 600

	grid := GameGrid{}
	grid.Array = [3][3]int8{{Empty, Empty, Empty}, {Empty, Empty, Empty}, {Empty, Empty, Empty}}

	game.Grid = &grid
	grid.Width = 400
	grid.Height = 400

	grid.PosX = 100
	grid.PosY = 100

	grid.CellCount = 3
	grid.CellHeight = grid.Height / grid.CellCount
	grid.CellWidth = grid.Width / grid.CellCount

	game.Running = true

	game.Mouse = &Mouse{}

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return
	}

	game.Window, err = sdl.CreateWindow("Tic Tac Toe", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, game.Width, game.Height, sdl.WINDOW_SHOWN)

	if err != nil {
		return
	}

	game.Window.SetBordered(true)
	game.Renderer, err = sdl.CreateRenderer(game.Window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return
	}

	return err
}
func (game *Game) HandleEvents() (err error) {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			game.Running = false
		case *sdl.MouseButtonEvent:
			if e.State == sdl.PRESSED {
				if e.Button == sdl.BUTTON_LEFT {
					game.Mouse.Clicked = true
					game.Mouse.X = e.X
					game.Mouse.Y = e.Y
				}
			}
		case *sdl.KeyboardEvent:
			if e.State == sdl.PRESSED {
				switch e.Keysym.Sym {
				case sdl.K_ESCAPE:
					if game.Grid.State != Running {
						game.Grid.State = Running
						game.Grid.Turn = PlayerX
						game.Grid.Array = [3][3]int8{{Empty, Empty, Empty}, {Empty, Empty, Empty}, {Empty, Empty, Empty}}
						game.Mouse.Clicked = false
					}
				case sdl.K_UP:
					if !game.Selected {
						game.Grid.Mode = Single
					}
				case sdl.K_DOWN:
					if !game.Selected {
						game.Grid.Mode = Multi
					}
				case sdl.K_RETURN:
					if !game.Selected {
						game.Selected = true

					}
				}
			}
		}
	}
	return nil
}

func (game *Game) Update() (err error) {
	grid := game.Grid
	if game.Selected && grid.State == Running {
		if grid.Mode == Multi || grid.Turn == PlayerX {
			mX := game.Mouse.X
			mY := game.Mouse.Y
			clicked := game.Mouse.Clicked
			if clicked && mX >= grid.PosX && mX <= (grid.PosX+grid.Width) && mY >= grid.PosY && mY <= (grid.Height+grid.PosY) {
				game.Mouse.Clicked = false
				mX = mX - grid.PosX
				mY = mY - grid.PosY
				column := mX / grid.CellWidth
				row := mY / grid.CellHeight
				if grid.Array[row][column] == Empty {
					if grid.Turn == PlayerX {
						grid.Array[row][column] = X
						grid.Turn = PlayerO
					} else {
						grid.Array[row][column] = O
						grid.Turn = PlayerX
					}
					grid.CheckLogic()
				}
			}
		} else {
			game.playComputer()
		}
	}
	return nil
}
func (game *Game) playComputer() {
	time.Sleep(500 * time.Millisecond)
	bestVal := math.MinInt
	var row, col int
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if game.Grid.Array[i][j] == Empty {
				game.Grid.Array[i][j] = O
				moveVal := game.Grid.miniMax(9, false)
				if moveVal > bestVal {
					bestVal = moveVal
					row = i
					col = j
				}
				game.Grid.Array[i][j] = Empty
			}
		}
	}
	game.Grid.Array[row][col] = O
	game.Grid.Turn = PlayerX
	game.Grid.CheckLogic()
}

func (grid *GameGrid) isMovesLeft() bool {
	for i := int32(0); i < grid.CellCount; i++ {
		for j := int32(0); j < grid.CellCount; j++ {
			if grid.Array[i][j] == Empty {
				return true
			}
		}
	}
	return false
}

func (grid *GameGrid) evaluate() int {
	for i := int32(0); i < grid.CellCount; i++ {
		if grid.Array[0][i] == grid.Array[1][i] && grid.Array[1][i] == grid.Array[2][i] {
			if grid.Array[0][i] != Empty {
				if grid.Array[0][i] == X {
					return -10
				} else {
					return +10
				}
			}
		}
		if grid.Array[i][0] == grid.Array[i][1] && grid.Array[i][1] == grid.Array[i][2] {
			if grid.Array[i][0] != Empty {
				if grid.Array[i][0] == X {
					return -10
				} else {
					return +10
				}
			}
		}
	}
	if grid.Array[0][0] == grid.Array[1][1] && grid.Array[1][1] == grid.Array[2][2] {
		if grid.Array[1][1] != Empty {
			if grid.Array[1][1] == X {
				return -10
			} else {
				return +10
			}
		}
	}
	if grid.Array[0][2] == grid.Array[1][1] && grid.Array[1][1] == grid.Array[2][0] {
		if grid.Array[1][1] != Empty {
			if grid.Array[1][1] == X {
				return -10
			} else {
				return +10
			}
		}
	}
	return 0
}

func (grid *GameGrid) miniMax(depth int, isMax bool) int {
	score := grid.evaluate()
	if score == 10 || score == -10 {
		return score
	}
	if !grid.isMovesLeft() {
		return 0
	}
	if isMax {
		best := math.MinInt
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if grid.Array[i][j] == Empty {
					grid.Array[i][j] = O
					new := grid.miniMax(depth-1, !isMax)
					if new > best {
						best = new
					}
					grid.Array[i][j] = Empty
				}
			}
		}
		return best
	} else {
		best := math.MaxInt
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if grid.Array[i][j] == Empty {
					grid.Array[i][j] = X
					new := grid.miniMax(depth-1, !isMax)
					if new < best {
						best = new
					}
					grid.Array[i][j] = Empty
				}
			}
		}
		return best
	}
}

func (grid *GameGrid) CheckLogic() {
	for i := int32(0); i < grid.CellCount; i++ {
		if grid.Array[0][i] == grid.Array[1][i] && grid.Array[1][i] == grid.Array[2][i] {
			if grid.Array[0][i] != Empty {
				if grid.Array[0][i] == X {
					grid.State = XWon
				} else {
					grid.State = OWon
				}
				grid.EndIndex = [2]int32{0*grid.CellCount + i, 2*grid.CellCount + i}
				return
			}
		}
		if grid.Array[i][0] == grid.Array[i][1] && grid.Array[i][1] == grid.Array[i][2] {
			if grid.Array[i][0] != Empty {
				if grid.Array[i][0] == X {
					grid.State = XWon
				} else {
					grid.State = OWon
				}
				grid.EndIndex = [2]int32{i*grid.CellCount + 0, i*grid.CellCount + 2}
				return
			}
		}
	}
	if grid.Array[0][0] == grid.Array[1][1] && grid.Array[1][1] == grid.Array[2][2] {
		if grid.Array[1][1] != Empty {
			if grid.Array[1][1] == X {
				grid.State = XWon
			} else {
				grid.State = OWon
			}
			grid.EndIndex = [2]int32{0*grid.CellCount + 0, 2*grid.CellCount + 2}
			return
		}
	}
	if grid.Array[0][2] == grid.Array[1][1] && grid.Array[1][1] == grid.Array[2][0] {
		if grid.Array[1][1] != Empty {
			if grid.Array[1][1] == X {
				grid.State = XWon
			} else {
				grid.State = OWon
			}
			grid.EndIndex = [2]int32{2*grid.CellCount + 0, 0*grid.CellCount + 2}
			return
		}
	}
	if !game.Grid.isMovesLeft() {
		grid.State = Draw
		return
	}
}
func (game *Game) Render() (err error) {
	game.Renderer.SetDrawColor(0, 0, 0, 0)
	game.Renderer.Clear()
	if err = game.renderTitleText(); err != nil {
		return
	}
	if game.Selected {
		game.renderGrid()
		game.renderExtraText()
	} else {
		game.renderSelection()
	}
	game.Renderer.Present()
	return
}

func (game *Game) renderSelection() (err error) {
	err = ttf.Init()
	if err != nil {
		return
	}
	defer ttf.Quit()
	font, err := ttf.OpenFont(fontPath, 50)
	if err != nil {
		return
	}
	defer font.Close()
	red := sdl.Color{R: 255, G: 0, B: 0, A: 0}
	yellow := sdl.Color{R: 255, G: 255, B: 0, A: 0}
	var multiColor, singleColor sdl.Color
	if game.Grid.Mode == Multi {
		multiColor = red
		singleColor = yellow
	} else {
		multiColor = yellow
		singleColor = red
	}
	textMulti, err := font.RenderUTF8Blended("Multiplayer", multiColor)
	if err != nil {
		return
	}
	defer textMulti.Free()
	textureMulti, err := game.Renderer.CreateTextureFromSurface(textMulti)
	if err != nil {
		return
	}
	defer textureMulti.Destroy()
	srcMulti := sdl.Rect{X: 0, Y: 0, W: textMulti.W, H: textMulti.H}
	destMulti := sdl.Rect{X: (game.Width - textMulti.W) / 2, Y: (game.Grid.PosY-textMulti.H)/2 + game.Height*2/3, W: textMulti.W, H: textMulti.H}
	game.Renderer.Copy(textureMulti, &srcMulti, &destMulti)

	textSingle, err := font.RenderUTF8Blended("Computer", singleColor)
	if err != nil {
		return
	}
	defer textSingle.Free()
	textureSingle, err := game.Renderer.CreateTextureFromSurface(textSingle)
	if err != nil {
		return
	}
	defer textureSingle.Destroy()
	srcSingle := sdl.Rect{X: 0, Y: 0, W: textSingle.W, H: textSingle.H}
	destSingle := sdl.Rect{X: (game.Width - textSingle.W) / 2, Y: 200 + 1/3*game.Height, W: textSingle.W, H: textSingle.H}
	game.Renderer.Copy(textureSingle, &srcSingle, &destSingle)
	return
}

func (game *Game) renderTitleText() (err error) {
	err = ttf.Init()
	if err != nil {
		return
	}
	defer ttf.Quit()
	font, err := ttf.OpenFont(fontPath, 80)
	if err != nil {
		return
	}
	defer font.Close()
	text, err := font.RenderUTF8Blended("Tic Tac Toe", sdl.Color{R: 0, G: 255, B: 255, A: 0})
	if err != nil {
		return
	}
	defer text.Free()
	texture, err := game.Renderer.CreateTextureFromSurface(text)
	if err != nil {
		return
	}
	defer texture.Destroy()
	src := sdl.Rect{X: 0, Y: 0, W: text.W, H: text.H}
	dest := sdl.Rect{X: (game.Width - text.W) / 2, Y: (game.Grid.PosY - text.H) / 2, W: text.W, H: text.H}
	game.Renderer.Copy(texture, &src, &dest)
	return
}

func (game *Game) renderExtraText() (err error) {
	err = ttf.Init()
	if err != nil {
		return
	}
	defer ttf.Quit()
	font, err := ttf.OpenFont(fontPath, 50)
	if err != nil {
		return
	}
	defer font.Close()
	var str string

	if game.Grid.State == Running {
		if game.Grid.Turn == PlayerX {
			str = "X's Turn"
		} else {
			str = "O's Turn"
		}
	} else if game.Grid.State == Draw {
		str = "Draw"
	} else {
		if game.Grid.State == XWon {
			str = "X Wins"
		} else {
			str = "O Wins"
		}
	}
	text, err := font.RenderUTF8Blended(str, sdl.Color{R: 255, G: 255, B: 0, A: 0})
	if err != nil {
		return
	}
	defer text.Free()
	texture, err := game.Renderer.CreateTextureFromSurface(text)
	if err != nil {
		return
	}
	defer texture.Destroy()
	src := sdl.Rect{X: 0, Y: 0, W: text.W, H: text.H}
	dest := sdl.Rect{X: (game.Width - text.W) / 2, Y: game.Grid.PosY + game.Grid.Height + (game.Grid.PosY-text.H)/2, W: text.W, H: text.H}
	game.Renderer.Copy(texture, &src, &dest)

	return
}

func (game *Game) renderTitle() (err error) {
	image, err := img.Load("assets/ttt.png")
	if err != nil {
		return
	}
	defer image.Free()

	texture, err := game.Renderer.CreateTextureFromSurface(image)
	if err != nil {
		return
	}
	defer texture.Destroy()
	src := sdl.Rect{X: 0, Y: 0, W: image.W, H: image.H}
	dest := sdl.Rect{X: (game.Width - image.W) / 2, Y: (game.Grid.PosY - image.H) / 2, W: image.W, H: image.H}
	game.Renderer.Copy(texture, &src, &dest)

	return
}

func (game *Game) Clean() {
	game.Renderer.Destroy()
	game.Window.Destroy()
	sdl.Quit()
}

func (game *Game) renderGrid() {
	grid := game.Grid
	for count := int32(0); count < 4; count++ {
		var lineWidth int32
		if count == 0 || count == 3 {
			lineWidth = 4
		} else {
			lineWidth = 2
		}
		gfx.ThickLineColor(game.Renderer, grid.PosX+0, grid.PosY+count*grid.CellHeight, grid.PosX+grid.Width, grid.PosY+count*grid.CellHeight, lineWidth, sdl.Color{R: 255, G: 0, B: 0, A: 255})
		gfx.ThickLineColor(game.Renderer, grid.PosX+count*grid.CellWidth, grid.PosY+0, grid.PosX+count*grid.CellWidth, grid.PosY+grid.Height, lineWidth, sdl.Color{R: 255, G: 0, B: 0, A: 255})
	}
	for i := int32(0); i < grid.CellCount; i++ {
		for j := int32(0); j < grid.CellCount; j++ {
			if grid.Array[i][j] == X {
				gfx.ThickLineColor(game.Renderer, grid.PosX+j*grid.CellWidth+grid.CellWidth/4, grid.PosY+i*grid.CellHeight+grid.CellHeight/4, grid.PosX+(j+1)*grid.CellWidth-grid.CellWidth/4, grid.PosY+(i+1)*grid.CellHeight-grid.CellHeight/4, 4, sdl.Color{R: 255, G: 0, B: 255, A: 255})
				gfx.ThickLineColor(game.Renderer, grid.PosX+(j+1)*grid.CellWidth-grid.CellWidth/4, grid.PosY+i*grid.CellHeight+grid.CellHeight/4, grid.PosX+(j)*grid.CellWidth+grid.CellWidth/4, grid.PosY+(i+1)*grid.CellHeight-grid.CellHeight/4, 4, sdl.Color{R: 255, G: 0, B: 255, A: 255})
			} else if grid.Array[i][j] == O {
				gfx.AACircleColor(game.Renderer, grid.PosX+j*grid.CellWidth+grid.CellWidth/2, grid.PosY+i*grid.CellHeight+grid.CellHeight/2, grid.CellHeight/4, sdl.Color{R: 0, G: 255, B: 255, A: 255})
				gfx.AACircleColor(game.Renderer, grid.PosX+j*grid.CellWidth+grid.CellWidth/2, grid.PosY+i*grid.CellHeight+grid.CellHeight/2, grid.CellHeight/4-1, sdl.Color{R: 0, G: 255, B: 255, A: 255})
				gfx.AACircleColor(game.Renderer, grid.PosX+j*grid.CellWidth+grid.CellWidth/2, grid.PosY+i*grid.CellHeight+grid.CellHeight/2, grid.CellHeight/4-2, sdl.Color{R: 0, G: 255, B: 255, A: 255})
			}
		}
	}
	if grid.State == XWon || grid.State == OWon {
		startRow := grid.EndIndex[0] / grid.CellCount
		startColumn := grid.EndIndex[0] % grid.CellCount
		endRow := grid.EndIndex[1] / grid.CellCount
		endColumn := grid.EndIndex[1] % grid.CellCount
		if startColumn == endColumn {
			gfx.ThickLineColor(game.Renderer, grid.PosX+startColumn*grid.CellWidth+grid.CellWidth/2, grid.PosY+0, grid.PosX+startColumn*grid.CellWidth+grid.CellWidth/2, grid.PosY+grid.Height, 4, sdl.Color{R: 255, G: 255, B: 255, A: 100})
		} else if startRow == endRow {
			gfx.ThickLineColor(game.Renderer, grid.PosX+0, grid.PosY+startRow*grid.CellHeight+grid.CellHeight/2, grid.PosX+grid.Width, grid.PosY+startRow*grid.CellHeight+grid.CellHeight/2, 4, sdl.Color{R: 255, G: 255, B: 255, A: 100})
		} else {
			if startRow == startColumn {
				gfx.ThickLineColor(game.Renderer, grid.PosX+0, grid.PosY+0, grid.PosX+grid.Width, grid.PosY+grid.Height, 4, sdl.Color{R: 255, G: 255, B: 255, A: 100})
			} else {
				gfx.ThickLineColor(game.Renderer, grid.PosX+grid.Width, grid.PosY+0, grid.PosX+0, grid.PosY+grid.Height, 4, sdl.Color{R: 255, B: 255, G: 255, A: 100})
			}
		}
	}
}
