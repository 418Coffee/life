package life

import (
	"math/rand"
	"strings"
)

// Field is the two-dimensional board of cells.
type Field struct {
	s             [][]bool
	width, height uint
}

// NewField allocates a new empty board of the given height and width.
func NewField(width, height uint) *Field {
	s := make([][]bool, height)
	for i := range s {
		s[i] = make([]bool, width)
	}
	return &Field{s: s, width: width, height: height}
}

// Set sets the value v to the cell with position x,y on the field.
func (f *Field) Set(x, y uint, v bool) {
	f.s[y][x] = v
}

// Alive reports whether cell at position x,y is alive or dead.
// If the x or y coordinates are outside the field boundaries they are wrapped toroidally.
// x=-1 -> width-1
func (f *Field) Alive(x, y int) bool {
	w, h := int(f.width), int(f.height)
	x += w
	x %= w
	y += h
	y %= h
	return f.s[y][x]
}

// Future returns the state of the cell at position x,y at the next tick.
// A cell is alive if:
//  - Any live cell with fewer than two live neighbours dies, as if by underpopulation.
//  - Any live cell with two or three live neighbours lives on to the next generation.
//  - Any live cell with more than three live neighbours dies, as if by overpopulation.
//  - Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
func (f *Field) Future(x, y uint) bool {
	var aliveNeighbours uint8
	ix, iy := int(x), int(y)
	// Start at position x-1,y-1 (top left corner) and work our way through the neighbouring cells.
	//	[
	//		[0, 0, 0]
	//		[0, 1, 0]
	//		[0, 0, 0]
	//	]
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (j != 0 || i != 0) && f.Alive(ix+i, iy+j) {
				aliveNeighbours++
			}
		}
	}
	return aliveNeighbours == 3 || aliveNeighbours == 2 && f.Alive(ix, iy)
}

// String is a string representation of the current state.
func (f *Field) String() string {
	w := new(strings.Builder)
	for y := 0; y < int(f.height); y++ {
		for x := 0; x < int(f.width); x++ {
			if f.Alive(x, y) {
				w.WriteRune('*')
			} else {
				w.WriteRune(' ')
			}
		}
		w.WriteRune('\n')
	}
	return w.String()
}

// Game stores the state of a round of Conway's Game of Life.
type Game struct {
	current, next *Field
	width, height uint
}

// uintn is basically Intn but casted to uintn
func uintn(n uint) uint {
	return uint(rand.Intn(int(n)))
}

// NewGame returns a new Life game state with a random initial state.
func NewGame(width, height uint) *Game {
	current := NewField(width, height)
	for i := 0; i < int(width*height/4); i++ {
		current.Set(uintn(width), uintn(height), true)
	}
	return &Game{
		current: current,
		next:    NewField(width, height),
		width:   width,
		height:  height,
	}
}

// Tick is a single discrete moment when births and deaths are processed.
func (g *Game) Tick() {
	for y := uint(0); y < g.height; y++ {
		for x := uint(0); x < g.width; x++ {
			g.next.Set(x, y, g.current.Future(x, y))
		}
	}
	g.current, g.next = g.next, g.current
}

// String is a string representation of the current game state.
func (g *Game) String() string {
	return g.current.String()
}
