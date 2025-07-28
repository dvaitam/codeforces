package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Point struct{ r, c int }

func expectedF(n, m int, grid []string) int {
	sr, sc := 0, 0
	empty := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			ch := grid[i][j]
			if ch != '#' {
				empty++
			}
			if ch == 'V' {
				sr = i
				sc = j
			}
		}
	}
	total := empty
	INF := n*m + 5
	distS := make([]int, n*m)
	for i := range distS {
		distS[i] = INF
	}
	idx := func(r, c int) int { return r*m + c }
	q := []int{}
	start := idx(sr, sc)
	distS[start] = 0
	q = append(q, start)
	dir := []int{-1, 0, 1, 0, -1}
	for head := 0; head < len(q); head++ {
		v := q[head]
		r := v / m
		c := v % m
		for k := 0; k < 4; k++ {
			nr := r + dir[k]
			nc := c + dir[k+1]
			if nr < 0 || nr >= n || nc < 0 || nc >= m || grid[nr][nc] == '#' {
				continue
			}
			ni := idx(nr, nc)
			if distS[ni] == INF {
				distS[ni] = distS[v] + 1
				q = append(q, ni)
			}
		}
	}
	exits := make([]Point, 0)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if i != 0 && i != n-1 && j != 0 && j != m-1 {
				continue
			}
			if grid[i][j] == '#' {
				continue
			}
			if distS[idx(i, j)] == INF {
				continue
			}
			exits = append(exits, Point{i, j})
		}
	}
	k := len(exits)
	if k == 0 {
		return total - 1
	}
	if k == 1 {
		d := distS[idx(exits[0].r, exits[0].c)]
		return total - (d + 1)
	}
	dist1 := make([]int, n*m)
	dist2 := make([]int, n*m)
	id1 := make([]int, n*m)
	id2 := make([]int, n*m)
	for i := range dist1 {
		dist1[i] = INF
		dist2[i] = INF
		id1[i] = -1
		id2[i] = -1
	}
	type state struct {
		pos int
		id  int
		d   int
	}
	queue := make([]state, 0)
	for idxE, e := range exits {
		p := idx(e.r, e.c)
		dist1[p] = 0
		id1[p] = idxE
		queue = append(queue, state{p, idxE, 0})
	}
	for head := 0; head < len(queue); head++ {
		st := queue[head]
		p := st.pos
		id := st.id
		d := st.d
		if dist1[p] == d && id1[p] == id {
			// propagate best
		} else if dist2[p] == d && id2[p] == id {
			// propagate second best
		} else {
			continue
		}
		r := p / m
		c := p % m
		for k := 0; k < 4; k++ {
			nr := r + dir[k]
			nc := c + dir[k+1]
			if nr < 0 || nr >= n || nc < 0 || nc >= m || grid[nr][nc] == '#' {
				continue
			}
			np := idx(nr, nc)
			nd := d + 1
			if nd < dist1[np] {
				if id1[np] != id {
					if dist1[np] < dist2[np] {
						dist2[np] = dist1[np]
						id2[np] = id1[np]
					}
				}
				dist1[np] = nd
				id1[np] = id
				queue = append(queue, state{np, id, nd})
			} else if id1[np] != id && nd < dist2[np] {
				dist2[np] = nd
				id2[np] = id
				queue = append(queue, state{np, id, nd})
			}
		}
	}
	best := INF
	for p := 0; p < n*m; p++ {
		if distS[p] == INF {
			continue
		}
		if id1[p] == -1 || id2[p] == -1 {
			continue
		}
		if id1[p] == id2[p] {
			continue
		}
		cost := distS[p] + dist1[p] + dist2[p] + 1
		if cost < best {
			best = cost
		}
	}
	return total - best
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(8) + 3
	m := rng.Intn(8) + 3
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if (i == 0 || i == n-1) && (j == 0 || j == m-1) {
				row[j] = '#'
				continue
			}
			if rng.Intn(4) == 0 {
				row[j] = '#'
			} else {
				row[j] = '.'
			}
		}
		grid[i] = string(row)
	}
	// place V
	for {
		r := rng.Intn(n)
		c := rng.Intn(m)
		if grid[r][c] != '#' {
			b := []byte(grid[r])
			b[c] = 'V'
			grid[r] = string(b)
			break
		}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(grid[i])
		sb.WriteByte('\n')
	}
	exp := expectedF(n, m, grid)
	return sb.String(), exp
}

func runCase(bin, input string, exp int) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != 1 {
		return fmt.Errorf("expected single integer output")
	}
	got, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid integer: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
