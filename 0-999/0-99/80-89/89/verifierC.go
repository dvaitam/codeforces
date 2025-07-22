package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// chip struct
type chip struct {
	r, c int
	dir  byte
}

type gridCase struct {
	n, m int
	grid []string
}

func genCaseC(rng *rand.Rand) gridCase {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	grid := make([]string, n)
	dirs := []byte{'L', 'R', 'U', 'D', '.'}
	for i := 0; i < n; i++ {
		b := make([]byte, m)
		for j := 0; j < m; j++ {
			b[j] = dirs[rng.Intn(len(dirs))]
		}
		grid[i] = string(b)
	}
	// ensure at least one chip
	has := false
	for _, row := range grid {
		if strings.ContainsAny(row, "LRUD") {
			has = true
			break
		}
	}
	if !has {
		i := rng.Intn(n)
		j := rng.Intn(m)
		d := dirs[rng.Intn(4)]
		row := []byte(grid[i])
		row[j] = d
		grid[i] = string(row)
	}
	return gridCase{n, m, grid}
}

func collectChips(gc gridCase) ([]chip, [][]int) {
	n, m := gc.n, gc.m
	chips := make([]chip, 0)
	id := make([][]int, n)
	for i := 0; i < n; i++ {
		id[i] = make([]int, m)
		for j := 0; j < m; j++ {
			id[i][j] = -1
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			d := gc.grid[i][j]
			if d == 'L' || d == 'R' || d == 'U' || d == 'D' {
				id[i][j] = len(chips)
				chips = append(chips, chip{i, j, d})
			}
		}
	}
	return chips, id
}

func solveC(gc gridCase) (int, int) {
	chips, _ := collectChips(gc)
	k := len(chips)
	if k == 0 {
		return 0, 0
	}
	l := make([]int, k)
	r := make([]int, k)
	u := make([]int, k)
	d := make([]int, k)
	for i := 0; i < k; i++ {
		l[i], r[i], u[i], d[i] = -1, -1, -1, -1
	}
	rowMap := make(map[int][]int)
	for i, ch := range chips {
		rowMap[ch.r] = append(rowMap[ch.r], i)
	}
	for _, idxs := range rowMap {
		sort.Slice(idxs, func(i, j int) bool { return chips[idxs[i]].c < chips[idxs[j]].c })
		for t := 0; t < len(idxs); t++ {
			if t > 0 {
				l[idxs[t]] = idxs[t-1]
			}
			if t+1 < len(idxs) {
				r[idxs[t]] = idxs[t+1]
			}
		}
	}
	colMap := make(map[int][]int)
	for i, ch := range chips {
		colMap[ch.c] = append(colMap[ch.c], i)
	}
	for _, idxs := range colMap {
		sort.Slice(idxs, func(i, j int) bool { return chips[idxs[i]].r < chips[idxs[j]].r })
		for t := 0; t < len(idxs); t++ {
			if t > 0 {
				u[idxs[t]] = idxs[t-1]
			}
			if t+1 < len(idxs) {
				d[idxs[t]] = idxs[t+1]
			}
		}
	}
	maxPoints := 0
	ways := 0
	for start := 0; start < k; start++ {
		lc := append([]int(nil), l...)
		rc := append([]int(nil), r...)
		uc := append([]int(nil), u...)
		dc := append([]int(nil), d...)
		points := 0
		curr := start
		for {
			var nxt int
			switch chips[curr].dir {
			case 'L':
				nxt = lc[curr]
			case 'R':
				nxt = rc[curr]
			case 'U':
				nxt = uc[curr]
			case 'D':
				nxt = dc[curr]
			}
			if lc[curr] != -1 {
				rc[lc[curr]] = rc[curr]
			}
			if rc[curr] != -1 {
				lc[rc[curr]] = lc[curr]
			}
			if uc[curr] != -1 {
				dc[uc[curr]] = dc[curr]
			}
			if dc[curr] != -1 {
				uc[dc[curr]] = uc[curr]
			}
			points++
			if nxt == -1 {
				break
			}
			curr = nxt
		}
		if points > maxPoints {
			maxPoints = points
			ways = 1
		} else if points == maxPoints {
			ways++
		}
	}
	return maxPoints, ways
}

func runCaseC(bin string, gc gridCase) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", gc.n, gc.m))
	for i := 0; i < gc.n; i++ {
		sb.WriteString(gc.grid[i])
		if i+1 < gc.n {
			sb.WriteByte('\n')
		}
	}
	sb.WriteByte('\n')
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got1, got2 int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got1, &got2); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp1, exp2 := solveC(gc)
	if got1 != exp1 || got2 != exp2 {
		return fmt.Errorf("expected %d %d got %d %d", exp1, exp2, got1, got2)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		gc := genCaseC(rng)
		if err := runCaseC(bin, gc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
