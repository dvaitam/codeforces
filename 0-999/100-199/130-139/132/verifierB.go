package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Block struct {
	color                  byte
	minx, miny, maxx, maxy int
}

func computeColor(grid []string, steps int64) byte {
	h := len(grid)
	w := len(grid[0])
	blkid := make([][]int, h)
	for i := range blkid {
		blkid[i] = make([]int, w)
		for j := range blkid[i] {
			blkid[i][j] = -1
		}
	}
	var blocks []Block
	dirs := [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if grid[y][x] != '0' && blkid[y][x] < 0 {
				col := grid[y][x]
				q := [][2]int{{x, y}}
				blkid[y][x] = len(blocks)
				minx, maxx, miny, maxy := x, x, y, y
				for qi := 0; qi < len(q); qi++ {
					cx, cy := q[qi][0], q[qi][1]
					if cx < minx {
						minx = cx
					}
					if cx > maxx {
						maxx = cx
					}
					if cy < miny {
						miny = cy
					}
					if cy > maxy {
						maxy = cy
					}
					for _, d := range dirs {
						nx, ny := cx+d[0], cy+d[1]
						if nx >= 0 && nx < w && ny >= 0 && ny < h && blkid[ny][nx] < 0 && grid[ny][nx] == col {
							blkid[ny][nx] = len(blocks)
							q = append(q, [2]int{nx, ny})
						}
					}
				}
				blocks = append(blocks, Block{color: col, minx: minx, miny: miny, maxx: maxx, maxy: maxy})
			}
		}
	}
	B := len(blocks)
	dx := [4]int{1, 0, -1, 0}
	dy := [4]int{0, 1, 0, -1}
	trans := make([]int, B*8)
	for id := 0; id < B; id++ {
		blk := blocks[id]
		for dp := 0; dp < 4; dp++ {
			for cp := 0; cp < 2; cp++ {
				state := id*8 + dp*2 + cp
				var tx, ty int
				switch dp {
				case 0:
					tx = blk.maxx
					if cp == 0 {
						ty = blk.miny
					} else {
						ty = blk.maxy
					}
				case 2:
					tx = blk.minx
					if cp == 0 {
						ty = blk.maxy
					} else {
						ty = blk.miny
					}
				case 1:
					ty = blk.maxy
					if cp == 0 {
						tx = blk.maxx
					} else {
						tx = blk.minx
					}
				case 3:
					ty = blk.miny
					if cp == 0 {
						tx = blk.minx
					} else {
						tx = blk.maxx
					}
				}
				nx, ny := tx+dx[dp], ty+dy[dp]
				if nx >= 0 && nx < w && ny >= 0 && ny < h && grid[ny][nx] != '0' {
					nid := blkid[ny][nx]
					trans[state] = nid*8 + dp*2 + cp
				} else {
					if cp == 0 {
						trans[state] = id*8 + dp*2 + 1
					} else {
						ndp := (dp + 1) & 3
						trans[state] = id*8 + ndp*2 + 0
					}
				}
			}
		}
	}
	startBlk := blkid[0][0]
	state := startBlk*8 + 0*2 + 0
	visited := make([]int, B*8)
	for i := range visited {
		visited[i] = -1
	}
	order := make([]int, 0, 10000)
	var step int64
	for step = 0; step < steps && visited[state] < 0; step++ {
		visited[state] = int(step)
		order = append(order, state)
		state = trans[state]
	}
	if step < steps {
		cycleStart := visited[state]
		cycleLen := int(step) - cycleStart
		rem := (steps - step) % int64(cycleLen)
		state = order[cycleStart+int(rem)]
	}
	bid := state / 8
	return blocks[bid].color
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildGrid(rng *rand.Rand) ([]string, int, string) {
	h := rng.Intn(4) + 1
	w := rng.Intn(4) + 1
	grid := make([][]byte, h)
	for i := range grid {
		grid[i] = make([]byte, w)
		for j := range grid[i] {
			grid[i][j] = '0'
		}
	}
	// place one rectangle covering (0,0)
	rx := rng.Intn(w)
	ry := rng.Intn(h)
	grid[ry][rx] = byte('1' + rng.Intn(9))
	// random other rectangles
	rects := rng.Intn(3)
	for r := 0; r < rects; r++ {
		x1 := rng.Intn(w)
		y1 := rng.Intn(h)
		x2 := rng.Intn(w)
		y2 := rng.Intn(h)
		if x1 > x2 {
			x1, x2 = x2, x1
		}
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		col := byte('1' + rng.Intn(9))
		for y := y1; y <= y2; y++ {
			for x := x1; x <= x2; x++ {
				grid[y][x] = col
			}
		}
	}
	rows := make([]string, h)
	for i := range grid {
		rows[i] = string(grid[i])
	}
	n := rng.Int63n(200) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", h, n)
	for i := range rows {
		sb.WriteString(rows[i])
		sb.WriteByte('\n')
	}
	return rows, int(n), sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		rows, n, input := buildGrid(rng)
		exp := computeColor(rows, int64(n))
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		if len(got) != 1 || got[0] != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %c got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
