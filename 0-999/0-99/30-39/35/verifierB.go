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

type op struct {
	add  bool
	x, y int
	id   string
}

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func simulate(n, m int, ops []op) string {
	grid := make([][]bool, n+1)
	for i := range grid {
		grid[i] = make([]bool, m+1)
	}
	pos := make(map[string][2]int)
	var out strings.Builder
	for _, op := range ops {
		if op.add {
			placed := false
			for sx := op.x; sx <= n && !placed; sx++ {
				start := 1
				if sx == op.x {
					start = op.y
				}
				for sy := start; sy <= m; sy++ {
					if !grid[sx][sy] {
						grid[sx][sy] = true
						pos[op.id] = [2]int{sx, sy}
						placed = true
						break
					}
				}
			}
		} else {
			if p, ok := pos[op.id]; ok {
				fmt.Fprintf(&out, "%d %d\n", p[0], p[1])
				grid[p[0]][p[1]] = false
				delete(pos, op.id)
			} else {
				out.WriteString("-1 -1\n")
			}
		}
	}
	return strings.TrimSpace(out.String())
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	k := rng.Intn(8) + 1
	var ops []op
	nextID := 1
	for i := 0; i < k; i++ {
		if rng.Intn(2) == 0 || nextID == 1 {
			x := rng.Intn(n) + 1
			y := rng.Intn(m) + 1
			id := fmt.Sprintf("id%d", nextID)
			nextID++
			ops = append(ops, op{add: true, x: x, y: y, id: id})
		} else {
			idn := rng.Intn(nextID-1) + 1
			id := fmt.Sprintf("id%d", idn)
			ops = append(ops, op{add: false, id: id})
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for _, op := range ops {
		if op.add {
			fmt.Fprintf(&sb, "+1 %d %d %s\n", op.x, op.y, op.id)
		} else {
			fmt.Fprintf(&sb, "-1 %s\n", op.id)
		}
	}
	input := sb.String()
	expected := simulate(n, m, ops)
	return input, expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases [][2]string
	// simple fixed case
	cases = append(cases, [2]string{"1 1 2\n+1 1 1 a\n-1 a\n", "1 1"})
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		cases = append(cases, [2]string{in, exp})
	}

	for i, tc := range cases {
		out, err := run(bin, tc[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc[0])
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc[1] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, tc[1], out, tc[0])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
