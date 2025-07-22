package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

type point struct{ x, y int }

func neighbors(p point, n, m int) []point {
	dirs := []point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	res := make([]point, 0, 4)
	for _, d := range dirs {
		nx, ny := p.x+d.x, p.y+d.y
		if nx >= 0 && nx < n && ny >= 0 && ny < m {
			res = append(res, point{nx, ny})
		}
	}
	return res
}

func bfs(n, m int, grid [][]byte, snake []point, apple point) int {
	length := len(snake)
	start := make([]point, length)
	copy(start, snake)
	type state struct {
		s []point
		d int
	}
	key := func(s []point) string {
		sb := strings.Builder{}
		for _, p := range s {
			sb.WriteByte(byte(p.x))
			sb.WriteByte(byte(p.y))
		}
		return sb.String()
	}
	visited := map[string]bool{key(start): true}
	q := []state{{start, 0}}
	for head := 0; head < len(q); head++ {
		cur := q[head]
		if cur.s[0] == apple {
			return cur.d
		}
		for _, nb := range neighbors(cur.s[0], n, m) {
			if grid[nb.x][nb.y] == '#' {
				continue
			}
			collide := false
			for i := 0; i < length-1; i++ {
				if cur.s[i] == nb {
					collide = true
					break
				}
			}
			if collide {
				continue
			}
			ns := make([]point, length)
			copy(ns[1:], cur.s[:length-1])
			ns[0] = nb
			k := key(ns)
			if !visited[k] {
				visited[k] = true
				q = append(q, state{ns, cur.d + 1})
			}
		}
	}
	return -1
}

func genCase(r *rand.Rand) (string, int) {
	for {
		n := r.Intn(3) + 3 //3..5
		m := r.Intn(3) + 3
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			grid[i] = make([]byte, m)
			for j := 0; j < m; j++ {
				grid[i][j] = '.'
			}
		}
		// snake length
		length := r.Intn(3) + 3 //3..5
		// random head
		hx := r.Intn(n)
		hy := r.Intn(m)
		snake := []point{{hx, hy}}
		used := map[point]bool{{hx, hy}: true}
		cur := point{hx, hy}
		ok := true
		for i := 1; i < length; i++ {
			opts := neighbors(cur, n, m)
			ps := make([]point, 0, len(opts))
			for _, p := range opts {
				if !used[p] {
					ps = append(ps, p)
				}
			}
			if len(ps) == 0 {
				ok = false
				break
			}
			nxt := ps[r.Intn(len(ps))]
			snake = append(snake, nxt)
			used[nxt] = true
			cur = nxt
		}
		if !ok {
			continue
		}
		for i, p := range snake {
			grid[p.x][p.y] = byte('1' + i)
		}
		// apple
		var ax, ay int
		for {
			ax = r.Intn(n)
			ay = r.Intn(m)
			if grid[ax][ay] == '.' {
				break
			}
		}
		grid[ax][ay] = '@'
		// some walls
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if grid[i][j] == '.' && r.Intn(5) == 0 {
					grid[i][j] = '#'
				}
			}
		}
		lines := make([]string, n)
		for i := 0; i < n; i++ {
			lines[i] = string(grid[i])
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for i := 0; i < n; i++ {
			sb.WriteString(lines[i])
			sb.WriteByte('\n')
		}
		input := sb.String()
		expect := bfs(n, m, grid, snake, point{ax, ay})
		return input, expect
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(4))
	for t := 1; t <= 100; t++ {
		input, expected := genCase(r)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\nInput:\n%s", t, err, input)
			return
		}
		if strings.TrimSpace(out) != fmt.Sprintf("%d", expected) {
			fmt.Printf("Test %d FAILED\nInput:\n%sExpected:%d Got:%s\n", t, input, expected, out)
			return
		}
	}
	fmt.Println("All tests passed")
}
