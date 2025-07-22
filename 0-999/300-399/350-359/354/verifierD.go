package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func expected(n int, pts [][2]int) int64 {
	k := len(pts)
	rowCols := make(map[int][]int, k)
	exist := make(map[int]map[int]bool, k)
	for _, p := range pts {
		r, c := p[0], p[1]
		if exist[r] == nil {
			exist[r] = make(map[int]bool)
		}
		exist[r][c] = true
		rowCols[r] = append(rowCols[r], c)
	}
	rows := make([]int, 0, len(rowCols))
	for r, cols := range rowCols {
		sort.Ints(cols)
		rowCols[r] = cols
		rows = append(rows, r)
	}
	sort.Ints(rows)
	dpRow := make(map[int]map[int]int, len(rows))
	for i := len(rows) - 1; i >= 0; i-- {
		r := rows[i]
		dpRow[r] = make(map[int]int, len(rowCols[r]))
		next := dpRow[r+1]
		for _, c := range rowCols[r] {
			d1, d2 := 0, 0
			if next != nil {
				d1 = next[c]
				d2 = next[c+1]
			}
			m := d1
			if d2 < m {
				m = d2
			}
			dpRow[r][c] = m + 1
		}
	}
	type cand struct{ dp, r, c int }
	cands := make([]cand, 0, k)
	for _, p := range pts {
		r, c := p[0], p[1]
		d := dpRow[r][c]
		if d > 0 {
			cands = append(cands, cand{d, r, c})
		}
	}
	sort.Slice(cands, func(i, j int) bool { return cands[i].dp > cands[j].dp })
	covered := make(map[int]map[int]bool, len(rowCols))
	for r := range rowCols {
		covered[r] = make(map[int]bool, len(rowCols[r]))
	}
	rem := k
	var cost int64 = 0
	for _, cd := range cands {
		if rem <= 0 {
			break
		}
		r, c, d := cd.r, cd.c, cd.dp
		if covered[r][c] {
			continue
		}
		h := d - 1
		tsize := d * (d + 1) / 2
		cost += int64(tsize + 2)
		for p := 0; p <= h; p++ {
			rr := r + p
			cols := rowCols[rr]
			if cols == nil {
				continue
			}
			l := sort.Search(len(cols), func(i int) bool { return cols[i] >= c })
			for idx := l; idx < len(cols) && cols[idx] <= c+p; idx++ {
				cc := cols[idx]
				if !covered[rr][cc] {
					covered[rr][cc] = true
					rem--
				}
			}
		}
	}
	if rem > 0 {
		cost += int64(rem * 3)
	}
	return cost
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		k := rng.Intn(n) + 1
		pts := make([][2]int, k)
		for j := 0; j < k; j++ {
			r := rng.Intn(n) + 1
			c := rng.Intn(r) + 1
			pts[j] = [2]int{r, c}
		}
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for _, p := range pts {
			input.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
		}
		expect := strconv.FormatInt(expected(n, pts), 10)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
