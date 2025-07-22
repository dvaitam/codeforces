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

type Gripper struct {
	x, y    int64
	m, p, r int64
}

type Node struct {
	p, r int64
}

func dist2(x1, y1, x2, y2 int64) int64 {
	dx := x1 - x2
	if dx < 0 {
		dx = -dx
	}
	dy := y1 - y2
	if dy < 0 {
		dy = -dy
	}
	return dx*dx + dy*dy
}

func solve(x0, y0, p0, r0 int64, gs []Gripper) int {
	n := len(gs)
	visited := make([]bool, n)
	queue := []Node{{p0, r0}}
	count := 0
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for i, g := range gs {
			if visited[i] {
				continue
			}
			if g.m <= cur.p && dist2(x0, y0, g.x, g.y) <= cur.r*cur.r {
				visited[i] = true
				queue = append(queue, Node{g.p, g.r})
				count++
			}
		}
	}
	return count
}

func generateCase(rng *rand.Rand) (string, string) {
	x := int64(rng.Intn(41) - 20)
	y := int64(rng.Intn(41) - 20)
	p := int64(rng.Intn(20) + 1)
	r := int64(rng.Intn(20) + 1)
	n := rng.Intn(20)
	gr := make([]Gripper, n)
	for i := 0; i < n; i++ {
		gx := int64(rng.Intn(41) - 20)
		gy := int64(rng.Intn(41) - 20)
		if gx == x && gy == y {
			gx++
		}
		m := int64(rng.Intn(20) + 1)
		p2 := int64(rng.Intn(20) + 1)
		r2 := int64(rng.Intn(20) + 1)
		gr[i] = Gripper{gx, gy, m, p2, r2}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d %d\n", x, y, p, r, n)
	for _, g := range gr {
		fmt.Fprintf(&sb, "%d %d %d %d %d\n", g.x, g.y, g.m, g.p, g.r)
	}
	ans := solve(x, y, p, r, gr)
	return sb.String(), fmt.Sprintf("%d", ans)
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
