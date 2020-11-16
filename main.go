package main

import (
    "fmt"
    "time"
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
    Spaceship
    QueenBee
)

type State uint8
const (
    Reproduction State = iota
    Underpopulation
    Overpopulation
    Unchanged
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
    case Spaceship:
        g.b[2][2] = true
        g.b[2][5] = true
        g.b[3][6] = true
        g.b[4][2] = true
        g.b[4][6] = true
        g.b[5][3] = true
        g.b[5][4] = true
        g.b[5][5] = true
        g.b[5][6] = true
    case QueenBee:
        g.b[4][15] = true
        g.b[5][14] = true
        g.b[5][15] = true
        g.b[6][13] = true
        g.b[6][14] = true
        g.b[6][19] = true
        g.b[6][20] = true
        g.b[7][3] = true
        g.b[7][4] = true
        g.b[7][12] = true
        g.b[7][13] = true
        g.b[7][14] = true
        g.b[7][19] = true
        g.b[7][20] = true
        g.b[7][23] = true
        g.b[7][24] = true
        g.b[8][3] = true
        g.b[8][4] = true
        g.b[8][13] = true
        g.b[8][14] = true
        g.b[8][19] = true
        g.b[8][20] = true
        g.b[8][23] = true
        g.b[8][24] = true
        g.b[9][14] = true
        g.b[9][15] = true
        g.b[10][15] = true

    default:
        panic(fmt.Sprintf("Undefined pattern: %s", p))
    }
}

func (g *Game) getState(h, w int) State {
    neighbors := func(h, w int) [][2]int {
        ngbrs := [][2]int{}

        for i:= -1; i <= 1; i++ {
            for j := -1; j <= 1; j++ {
                if i == 0 && j == 0 {
                    continue
                }

                h_i := h + i
                w_j := w + j

                if h_i < 0 || w_j < 0 || h_i >= BOARD_HEIGHT || w_j >= BOARD_WIDTH {
                    continue
                }

                ngbrs = append(ngbrs, [2]int{ h_i, w_j })
            }
        }

        return ngbrs
    }(h, w)

    liveCellCount := func(ngbrs [][2]int) int {
        cnt := 0
        for _, ngbr := range ngbrs {
            if g.b[ngbr[0]][ngbr[1]] {
                cnt++
            }
        }
        return cnt
    }(neighbors)

    if g.b[h][w] {
        if liveCellCount == 2 || liveCellCount == 3 {
            return Unchanged
        } else if liveCellCount >= 4 {
            return Overpopulation
        } else if liveCellCount <= 1 {
            return Underpopulation
        }
    } else {
        if liveCellCount == 3 {
            return Reproduction
        }
    }

    return Unchanged
}

func (g *Game) transition() {
    nextBoard := g.b
    for i := 0; i < BOARD_HEIGHT; i++ {
        for j := 0; j < BOARD_WIDTH; j++ {
            s := g.getState(i, j)
            switch s {
            case Reproduction:
                nextBoard[i][j] = true
            case Underpopulation, Overpopulation:
                nextBoard[i][j] = false
            case Unchanged:
            }
        }
    }
    g.b = nextBoard
}

func main() {
    g, err := newGame()
    if err != nil {
        panic(err)
    }

    g.set(QueenBee)

    running := true
    for running {
        g.draw()
        g.update()
        g.transition()
        time.Sleep(time.Millisecond * 100)
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
