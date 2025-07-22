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

type point struct{ x, y int }

func simulate(n, m int, blocks map[point]bool, xs, ys int, dx, dy int) int {
	visitedCells := map[point]bool{{xs, ys}: true}
	ans := 1
	states := map[[4]int]bool{}
	x, y := xs, ys
	for {
		key := [4]int{x, y, dx, dy}
		if states[key] {
			break
		}
		states[key] = true
		tx := n - x
		if dx < 0 {
			tx = x - 1
		}
		ty := m - y
		if dy < 0 {
			ty = y - 1
		}
		steps := tx
		if ty < steps {
			steps = ty
		}
		tObs := steps + 1
		for s := 1; s <= steps; s++ {
			nx := x + dx*s
			ny := y + dy*s
			if blocks[point{nx, ny}] {
				tObs = s
				break
			}
		}
		t := steps
		hitObs := false
		if tObs <= steps {
			t = tObs
			hitObs = true
		}
		for s := 1; s <= t; s++ {
			if hitObs && s == t {
				break
			}
			nx := x + dx*s
			ny := y + dy*s
			p := point{nx, ny}
			if !visitedCells[p] {
				visitedCells[p] = true
				ans++
			}
		}
		x += dx * t
		y += dy * t
		if hitObs {
			dx = -dx
			dy = -dy
		} else {
			if tx < ty {
				dx = -dx
			} else if ty < tx {
				dy = -dy
			} else {
				dx = -dx
				dy = -dy
			}
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	maxCells := n*m - 1
	k := 0
	if maxCells > 0 {
		k = rng.Intn(min(maxCells, 4) + 1)
	}
	blocks := map[point]bool{}
	for len(blocks) < k {
		x := rng.Intn(n) + 1
		y := rng.Intn(m) + 1
		blocks[point{x, y}] = true
	}
	xs := rng.Intn(n) + 1
	ys := rng.Intn(m) + 1
	for blocks[point{xs, ys}] {
		xs = rng.Intn(n) + 1
		ys = rng.Intn(m) + 1
	}
	dirs := []string{"NE", "NW", "SE", "SW"}
	dir := dirs[rng.Intn(4)]
	dx, dy := 0, 0
	switch dir {
	case "NE":
		dx = -1
		dy = 1
	case "NW":
		dx = -1
		dy = -1
	case "SE":
		dx = 1
		dy = 1
	case "SW":
		dx = 1
		dy = -1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for p := range blocks {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
	}
	sb.WriteString(fmt.Sprintf("%d %d %s\n", xs, ys, dir))
	res := simulate(n, m, blocks, xs, ys, dx, dy)
	return sb.String(), fmt.Sprintf("%d\n", res)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func runCase(exe string, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
