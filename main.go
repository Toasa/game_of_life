package main

import (
    "fmt"
    "github.com/veandco/go-sdl2/sdl"
)

const (
    WINDOW_NAME = "game of life"

    BOARD_HEIGHT = 30
    BOARD_WIDTH = 70

    CELL_LEN = 10
    GRID_WIDTH_LEN = 2
    HEIGHT_LEN = BOARD_HEIGHT * CELL_LEN + (BOARD_HEIGHT+1)*GRID_WIDTH_LEN
    WIDTH_LEN = BOARD_WIDTH * CELL_LEN + (BOARD_WIDTH+1)*GRID_WIDTH_LEN
)

const (
    COLOR_BG = 0x00DCDCDC
    COLOR_GRID = 0x00FFFFFF
    COLOR_CELL = 0x00000000
)

type Pattern uint8
const (
    Glider Pattern = iota
)

type Board [BOARD_HEIGHT][BOARD_WIDTH]bool

type Game struct {
    w *sdl.Window
    s *sdl.Surface
    b Board
}

func newGame() (*Game, error) {
    if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
        return nil, err
    }

    w, err := sdl.CreateWindow(
        WINDOW_NAME,
        sdl.WINDOWPOS_UNDEFINED,
        sdl.WINDOWPOS_UNDEFINED,
        int32(WIDTH_LEN),
        int32(HEIGHT_LEN),
        sdl.WINDOW_SHOWN,
    )
    if err != nil {
        return nil, err
    }

    s, err := w.GetSurface()
    if err != nil {
        return nil, err
    }

    return &Game{ w: w, s: s }, nil
}

func (g *Game) update() {
    g.w.UpdateSurface()
}

func (g *Game) draw() {
    // Draw back ground.
    bg := sdl.Rect { 0, 0, WIDTH_LEN, HEIGHT_LEN }
    g.s.FillRect(&bg, COLOR_BG)

    // Draw vertical grids.
    for i := 0; i < BOARD_HEIGHT+1; i++ {
        grid := sdl.Rect {
            0,
            int32((CELL_LEN+GRID_WIDTH_LEN)*i),
            WIDTH_LEN,
            GRID_WIDTH_LEN,
        }
        g.s.FillRect(&grid, COLOR_GRID)
    }

    // Draw horizontal grids.
    for i := 0; i < BOARD_WIDTH+1; i++ {
        grid := sdl.Rect {
            int32((CELL_LEN+GRID_WIDTH_LEN)*i),
            0,
            GRID_WIDTH_LEN,
            HEIGHT_LEN,
        }
        g.s.FillRect(&grid, COLOR_GRID)
    }

    // Draw cell.
    for i := 0; i < BOARD_HEIGHT; i++ {
        for j := 0; j < BOARD_WIDTH; j++ {
            if g.b[i][j] {
                cell := sdl.Rect {
                    int32((CELL_LEN+GRID_WIDTH_LEN)*j + GRID_WIDTH_LEN),
                    int32((CELL_LEN+GRID_WIDTH_LEN)*i + GRID_WIDTH_LEN),
                    CELL_LEN,
                    CELL_LEN,
                }
                g.s.FillRect(&cell, COLOR_CELL)
            }
        }
    }
}

func (g *Game) set(p Pattern) {
    switch p {
    case Glider:
        g.b[2][4] = true
        g.b[3][5] = true
        g.b[4][3] = true
        g.b[4][4] = true
        g.b[4][5] = true
    default:
        panic(fmt.Sprintf("Undefined pattern: %s", p))
    }
}

func main() {
    g, err := newGame()
    if err != nil {
        panic(err)
    }

    g.set(Glider)

    g.draw()
    g.update()

    running := true
    for running {
        for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
            switch e.(type) {
            case *sdl.QuitEvent:
                fmt.Println("Quit")
                running = false
                break
            }
        }
    }

}
