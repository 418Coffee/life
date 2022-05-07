package life

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Field is the two-dimensional board of cells.
type Field struct {
	s             [][]bool
	width, height uint
	wrap          bool
}

// NewField allocates a new empty board of the given height and width.
func NewField(width, height uint, wrap bool) *Field {
	s := make([][]bool, height)
	for i := range s {
		s[i] = make([]bool, width)
	}
	return &Field{s: s, width: width, height: height, wrap: wrap}
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
			if !f.wrap {
				if ix+i < 0 || iy+j < 0 {
					continue
				}
			}
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
	wrap          bool
	comment       string
}

// uintn is basically Intn but casted to uintn
func uintn(n uint) uint {
	return uint(rand.Intn(int(n)))
}

// NewGame returns a new Life game state with a random initial state.
func NewGame(width, height uint, wrap bool) *Game {
	current := NewField(width, height, wrap)
	for i := 0; i < int(width*height/4); i++ {
		current.Set(uintn(width), uintn(height), true)
	}
	return &Game{
		current: current,
		next:    NewField(width, height, wrap),
		width:   width,
		height:  height,
	}
}

var (
	widthHeightRegex = regexp.MustCompile(`\d+`)
	lifeRuleRegex    = regexp.MustCompile(`(?i)b3/s23`)
	itemRegex        = regexp.MustCompile(`(?U)\d+[b|o]+|[b|o]`)
	rule             = []byte{'r', 'u', 'l', 'e'}
)

// LoadGame loads a Life game state from a run-length encoded file.
// The file must have the .rle extension.
// An error is returned if an error occurred when reading the file or when parsing the contents.
func LoadGame(filename string, wrap bool) (*Game, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	} else if info.IsDir() {
		return nil, fmt.Errorf("%s is a directory", filename)
	}
	if filepath.Ext(filename) != ".rle" {
		return nil, fmt.Errorf("only RLE files are supported currently")
	}
	comment := new(strings.Builder)
	scanner := bufio.NewScanner(f)
	game := new(Game)
	game.wrap = true
	for scanner.Scan() {
		if line := scanner.Bytes(); len(line) > 0 {
			if len(line) > 70 {
				// Lines in the RLE file must not exceed 70 characters, although it is a good idea for RLE readers to be able to cope with longer lines.
				fmt.Printf("warning: line exceeds 70 characters: %s\n", line)
				continue
			}
			if line[0] == '#' {
				// Comment line
				// Skip the 3 preceding bytes and append a new line for printing purposes.
				comment.Write(append(line[3:], '\n'))
			} else if line[0] == 'x' {
				// Alternative rules are not supported.
				if bytes.Contains(line, rule) && !lifeRuleRegex.Match(line) {
					return nil, fmt.Errorf("rules are not supported")
				}
				widthHeight := widthHeightRegex.FindAll(line, 2)
				if len(widthHeight) != 2 {
					return nil, fmt.Errorf("got %d parameters expected 2", len(widthHeight))
				}
				width, err := strconv.ParseUint(string(widthHeight[0]), 0, 64)
				if err != nil {
					return nil, err
				}
				height, err := strconv.ParseUint(string(widthHeight[1]), 0, 64)
				if err != nil {
					return nil, err
				}
				game.width, game.height = uint(width), uint(height)
				game.current = NewField(game.width, game.height, wrap)
			} else {
				// If we haven't encountered a header line this file is invalid.
				if game.current == nil {
					return nil, fmt.Errorf("invalid RLE format")
				}
				// An exclamation mark marks the end of the configuration.
				for !bytes.ContainsRune(line, '!') {
					scanner.Scan()
					line = append(line, scanner.Bytes()...)
				}
				game.current.s = make([][]bool, game.height)
				// Dead cells at the end of the last line of the pattern do not need to be encoded.
				definedCells := bytes.Split(line, []byte{'$'})
				for i, item := range definedCells {
					if game.current.s[i], err = generateLine(item, game.width); err != nil {
						return nil, err
					}
				}
				// Define any dead cells if necessary.
				for i := uint(len(definedCells)); i < game.height; i++ {
					game.current.s[i] = make([]bool, game.width)
				}
			}
		}

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	game.comment = comment.String()
	game.next = NewField(game.width, game.height, wrap)
	return game, nil
}

// <tag>	description
//   b		dead cell
//   o		alive cell
func tag2bool(tag byte) bool {
	return tag == 'o'
}

// generateLine generates a slice representing one pattern line from an RLE file.
// It follows all standards proposed by: https://conwaylife.com/wiki/Run_Length_Encoded#Description_of_format.
// This function may be considered as dangerous because it makes an allocation with a size specified by the user.
// An attacker could theoretically starve memory with this function.
func generateLine(rawItem []byte, width uint) ([]bool, error) {
	line := make([]bool, 0, width)
	for _, item := range itemRegex.FindAll(rawItem, -1) {
		if len(item) == 1 {
			// run_count was omitted because it's equal to 1.
			line = append(line, tag2bool(item[0]))
		} else {
			// run_count not emitted.
			rawRunCount, tag := item[:len(item)-1], item[len(item)-1]
			runCount, err := strconv.ParseUint(string(rawRunCount), 0, 64)
			if err != nil {
				return nil, err
			}
			values := make([]bool, runCount)
			// Because our slice is initialized with all 0 values (false), we only need set values when it's necessary.
			if tag2bool(tag) {
				for i := range values {
					values[i] = true
				}
			}
			line = append(line, values...)
		}
	}
	// Dead cells at the end of a pattern line do not need to be encoded.
	return append(line, make([]bool, width-uint(len(line)))...), nil
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

// Comment returns the comment(s) of the loaded RLE file.
// A string with length 0 is returned if there are no comments, or the game was created using NewGame.
func (g *Game) Comment() string {
	return g.comment
}
