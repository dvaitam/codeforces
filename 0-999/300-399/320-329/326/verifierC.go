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

func expectedResult(r int, reclaimed [][2]int) string {
	avail := make([][3]bool, r+2)
	for i := 1; i <= r; i++ {
		avail[i][1], avail[i][2] = true, true
	}
	rec := make([][3]bool, r+2)
	for _, rc := range reclaimed {
		rec[rc[0]][rc[1]] = true
	}
	for i := 1; i <= r; i++ {
		for c := 1; c <= 2; c++ {
			if rec[i][c] {
				avail[i][c] = false
				oc := 3 - c
				for d := -1; d <= 1; d++ {
					j := i + d
					if j >= 1 && j <= r {
						avail[j][oc] = false
					}
				}
			}
		}
	}
	maxr := r
	g := make([][][3]int, maxr+1)
	for i := 0; i <= maxr; i++ {
		g[i] = make([][3]int, 3)
	}
	for length := 1; length <= maxr; length++ {
		for tb := 0; tb < 3; tb++ {
			for bb := 0; bb < 3; bb++ {
				used := make([]bool, length*2+5)
				for i := 1; i <= length; i++ {
					for c := 1; c <= 2; c++ {
						if (i == 1 && tb == c) || (i == length && bb == c) {
							continue
						}
						var g1 int
						if i > 1 {
							bb1 := 3 - c
							g1 = g[i-1][tb][bb1]
						}
						var g2 int
						if i < length {
							tb2 := 3 - c
							g2 = g[length-i][tb2][bb]
						}
						x := g1 ^ g2
						if x < len(used) {
							used[x] = true
						}
					}
				}
				mex := 0
				for mex < len(used) && used[mex] {
					mex++
				}
				g[length][tb][bb] = mex
			}
		}
	}
	total := 0
	segStart := 0
	segTB := 0
	inSeg := false
	for i := 1; i <= r; i++ {
		if !avail[i][1] && !avail[i][2] {
			if inSeg {
				length := i - segStart
				bb := 0
				if !avail[i-1][1] {
					bb = 1
				}
				if !avail[i-1][2] {
					bb = 2
				}
				total ^= g[length][segTB][bb]
			}
			inSeg = false
		} else {
			if !inSeg {
				inSeg = true
				segStart = i
				segTB = 0
				if !avail[i][1] {
					segTB = 1
				}
				if !avail[i][2] {
					segTB = 2
				}
			}
		}
	}
	if inSeg {
		length := r - segStart + 1
		bb := 0
		if !avail[r][1] {
			bb = 1
		}
		if !avail[r][2] {
			bb = 2
		}
		total ^= g[length][segTB][bb]
	}
	if total != 0 {
		return "WIN"
	}
	return "LOSE"
}

func runCase(bin string, r int, cells [][2]int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", r, len(cells)))
	for _, c := range cells {
		sb.WriteString(fmt.Sprintf("%d %d\n", c[0], c[1]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := expectedResult(r, cells)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func generateCase(rng *rand.Rand) (int, [][2]int) {
	r := rng.Intn(6) + 1
	n := rng.Intn(r + 1)
	board := make([][3]bool, r+2)
	cells := make([][2]int, 0, n)
	for len(cells) < n {
		row := rng.Intn(r) + 1
		col := rng.Intn(2) + 1
		if board[row][col] {
			continue
		}
		oc := 3 - col
		if board[row][oc] || board[row-1][oc] || board[row+1][oc] {
			continue
		}
		board[row][col] = true
		cells = append(cells, [2]int{row, col})
	}
	return r, cells
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		r, cells := generateCase(rng)
		if err := runCase(bin, r, cells); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
