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

type seg struct{ r, a, b int }

type point struct{ x, y int }

func solve(x0, y0, x1, y1 int, segs []seg) string {
	allowed := make(map[point]bool)
	for _, s := range segs {
		for c := s.a; c <= s.b; c++ {
			allowed[point{s.r, c}] = true
		}
	}
	start := point{x0, y0}
	goal := point{x1, y1}
	dirs := [][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}
	dist := map[point]int{start: 0}
	q := []point{start}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if cur == goal {
			return fmt.Sprintf("%d\n", dist[cur])
		}
		for _, d := range dirs {
			nx, ny := cur.x+d[0], cur.y+d[1]
			np := point{nx, ny}
			if !allowed[np] {
				continue
			}
			if _, ok := dist[np]; !ok {
				dist[np] = dist[cur] + 1
				q = append(q, np)
			}
		}
	}
	return "-1\n"
}

func genCase(rng *rand.Rand) (string, string) {
	rows := rng.Intn(4) + 2
	cols := rng.Intn(4) + 2
	n := rows
	segs := make([]seg, 0, n)
	allowed := make([]point, 0, rows*cols)
	for i := 1; i <= rows; i++ {
		a := rng.Intn(cols) + 1
		b := a + rng.Intn(cols-a+1)
		segs = append(segs, seg{i, a, b})
		for c := a; c <= b; c++ {
			allowed = append(allowed, point{i, c})
		}
	}
	if len(allowed) < 2 {
		return genCase(rng)
	}
	p1 := allowed[rng.Intn(len(allowed))]
	p2 := allowed[rng.Intn(len(allowed))]
	for p2 == p1 {
		p2 = allowed[rng.Intn(len(allowed))]
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", p1.x, p1.y, p2.x, p2.y)
	fmt.Fprintf(&sb, "%d\n", len(segs))
	for _, s := range segs {
		fmt.Fprintf(&sb, "%d %d %d\n", s.r, s.a, s.b)
	}
	out := solve(p1.x, p1.y, p2.x, p2.y, segs)
	return sb.String(), out
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	if strings.TrimSpace(buf.String()) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected \n%s\ngot \n%s", exp, buf.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
