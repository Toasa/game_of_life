package main

import (
    "fmt"
    "time"
    "strings"
    "github.com/veandco/go-sdl2/sdl"
)

const (
    WINDOW_NAME = "game of life"

    BOARD_HEIGHT = 70
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
    GliderGun
    Spaceship
    QueenBee
)

const FigGlider string =
`......
......
....#.
.....#
...###`

const FigGliderGun string =
`.......................................
.......................................
.......................................
...........................#...........
.........................#.#...........
...............##......##............##
..............#...#....##............##
...##........#.....#...##..............
...##........#...#.##....#.#...........
.............#.....#.......#...........
..............#...#....................
...............##......................`

const FigSpaceship string =
`.......
.......
..#..#.
......#
..#...#
...####`

const FigQueenBee string =
`.........................
.........................
.........................
.........................
...............#.........
..............##.........
.............##....##....
...##.......###....##..##
...##........##....##..##
..............##.........
...............#.........`

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
    var fig string

    switch p {
    case Glider:
        fig = FigGlider
    case GliderGun:
        fig = FigGliderGun
    case Spaceship:
        fig = FigSpaceship
    case QueenBee:
        fig = FigQueenBee
    default:
        panic(fmt.Sprintf("Undefined pattern: %s", p))
    }

    figLines := strings.Split(fig, "\n")
    for i, line := range(figLines) {
        for j, c := range(line) {
            if c == '#' {
                g.b[i][j] = true
            } else {
                g.b[i][j] = false
            }
        }
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

    g.set(GliderGun)

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
