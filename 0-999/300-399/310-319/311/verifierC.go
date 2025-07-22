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

func recomputeReach(methods []int, maxPos int) []bool {
	reachable := make([]bool, maxPos+1)
	queue := []int{0}
	reachable[0] = true
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, step := range methods {
			next := cur + step
			if next <= maxPos && !reachable[next] {
				reachable[next] = true
				queue = append(queue, next)
			}
		}
	}
	return reachable
}

type Operation struct {
	typ int
	x   int
	y   int64
}

func solveC(h int, n, m, k0 int, pos []int, c []int64, ops []Operation) string {
	methods := []int{k0}
	exists := make([]bool, n)
	for i := range exists {
		exists[i] = true
	}
	maxPos := h
	for _, p := range pos {
		if p > maxPos {
			maxPos = p
		}
	}
	var sb strings.Builder
	for _, op := range ops {
		if op.typ == 1 {
			methods = append(methods, op.x)
		} else if op.typ == 2 {
			idx := op.x - 1
			c[idx] -= op.y
		}
		reachable := recomputeReach(methods, maxPos)
		if op.typ == 3 {
			bestVal := int64(0)
			bestIdx := -1
			for i := 0; i < n; i++ {
				if exists[i] && pos[i]-1 <= maxPos && reachable[pos[i]-1] {
					if c[i] > bestVal || (c[i] == bestVal && (bestIdx == -1 || i < bestIdx)) {
						bestVal = c[i]
						bestIdx = i
					}
				}
			}
			if bestIdx == -1 {
				sb.WriteString("0\n")
			} else {
				sb.WriteString(fmt.Sprintf("%d\n", bestVal))
				exists[bestIdx] = false
			}
		}
	}
	return strings.TrimSpace(sb.String())
}

func generateCaseC(rng *rand.Rand) (string, string) {
	h := rng.Intn(50) + 10
	n := rng.Intn(5) + 1
	m := rng.Intn(10) + 1
	k0 := rng.Intn(5) + 1
	pos := make([]int, n)
	used := map[int]bool{}
	for i := 0; i < n; i++ {
		for {
			p := rng.Intn(h) + 1
			if !used[p] {
				used[p] = true
				pos[i] = p
				break
			}
		}
	}
	c := make([]int64, n)
	for i := range c {
		c[i] = int64(rng.Intn(50) + 1)
	}
	ops := make([]Operation, 0, m)
	queryCount := 0
	for i := 0; i < m; i++ {
		typ := rng.Intn(3) + 1
		if typ == 1 {
			x := rng.Intn(5) + 1
			ops = append(ops, Operation{typ: 1, x: x})
		} else if typ == 2 {
			x := rng.Intn(n) + 1
			y := int64(rng.Intn(10) + 1)
			ops = append(ops, Operation{typ: 2, x: x, y: y})
		} else {
			ops = append(ops, Operation{typ: 3})
			queryCount++
		}
	}
	if queryCount == 0 {
		ops = append(ops, Operation{typ: 3})
		m++
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", h, n, m, k0)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", pos[i], c[i])
	}
	for _, op := range ops {
		if op.typ == 1 {
			fmt.Fprintf(&sb, "1 %d\n", op.x)
		} else if op.typ == 2 {
			fmt.Fprintf(&sb, "2 %d %d\n", op.x, op.y)
		} else {
			fmt.Fprintf(&sb, "3\n")
		}
	}
	expected := solveC(h, n, m, k0, pos, c, ops)
	return sb.String(), expected
}

func runCase(bin, input, expected string) error {
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
	outStr := strings.TrimSpace(out.String())
	if outStr != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, outStr)
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
		in, exp := generateCaseC(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
