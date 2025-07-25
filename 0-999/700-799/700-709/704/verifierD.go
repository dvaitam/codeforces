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

func isValid(colors []byte, n, m int, rCost, bCost int, xs, ys []int, cons [][3]int) (bool, int) {
	cost := 0
	for i := 0; i < n; i++ {
		if colors[i] == 'r' {
			cost += rCost
		} else {
			cost += bCost
		}
	}
	for _, c := range cons {
		t, l, d := c[0], c[1], c[2]
		red, blue := 0, 0
		for i := 0; i < n; i++ {
			if (t == 1 && xs[i] == l) || (t == 2 && ys[i] == l) {
				if colors[i] == 'r' {
					red++
				} else {
					blue++
				}
			}
		}
		if abs(red-blue) > d {
			return false, 0
		}
	}
	return true, cost
}

func solveD(n, m, rCost, bCost int, xs, ys []int, cons [][3]int) string {
	best := -1
	bestMask := 0
	for mask := 0; mask < (1 << n); mask++ {
		colors := make([]byte, n)
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				colors[i] = 'r'
			} else {
				colors[i] = 'b'
			}
		}
		ok, cost := isValid(colors, n, m, rCost, bCost, xs, ys, cons)
		if !ok {
			continue
		}
		if best == -1 || cost < best {
			best = cost
			bestMask = mask
		}
	}
	if best == -1 {
		return "-1"
	}
	colors := make([]byte, n)
	for i := 0; i < n; i++ {
		if bestMask&(1<<i) != 0 {
			colors[i] = 'r'
		} else {
			colors[i] = 'b'
		}
	}
	return fmt.Sprintf("%d\n%s", best, string(colors))
}

func generateCase(rng *rand.Rand) (string, string, int, int, int, int, []int, []int, [][3]int) {
	n := rng.Intn(4) + 1
	m := rng.Intn(3)
	r := rng.Intn(5) + 1
	b := rng.Intn(5) + 1
	xs := make([]int, n)
	ys := make([]int, n)
	for i := 0; i < n; i++ {
		xs[i] = rng.Intn(5) + 1
		ys[i] = rng.Intn(5) + 1
	}
	cons := make([][3]int, m)
	for i := 0; i < m; i++ {
		t := rng.Intn(2) + 1
		l := rng.Intn(5) + 1
		d := rng.Intn(n + 1)
		cons[i] = [3]int{t, l, d}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	sb.WriteString(fmt.Sprintf("%d %d\n", r, b))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", xs[i], ys[i]))
	}
	for i := 0; i < m; i++ {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", cons[i][0], cons[i][1], cons[i][2]))
	}
	expected := solveD(n, m, r, b, xs, ys, cons)
	return sb.String(), expected, n, m, r, b, xs, ys, cons
}

func runCase(bin, input, expected string, n, m, rCost, bCost int, xs, ys []int, cons [][3]int) error {
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
	expLines := strings.Split(strings.TrimSpace(expected), "\n")
	outLines := strings.Split(outStr, "\n")
	if expLines[0] == "-1" {
		if outLines[0] != "-1" {
			return fmt.Errorf("expected -1 got %s", outLines[0])
		}
		return nil
	}
	if outLines[0] != expLines[0] {
		return fmt.Errorf("expected cost %s got %s", expLines[0], outLines[0])
	}
	if len(outLines) < 2 {
		return fmt.Errorf("missing color string")
	}
	colors := outLines[1]
	if len(colors) != n {
		return fmt.Errorf("invalid color length")
	}
	// verify assignment validity
	xsCopy := make([]int, len(xs))
	ysCopy := make([]int, len(ys))
	copy(xsCopy, xs)
	copy(ysCopy, ys)
	ok, cost := isValid([]byte(colors), n, m, rCost, bCost, xsCopy, ysCopy, cons)
	if !ok {
		return fmt.Errorf("assignment violates constraints")
	}
	if fmt.Sprintf("%d", cost) != expLines[0] {
		return fmt.Errorf("cost mismatch")
	}
	return nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp, n, m, r, b, xs, ys, cons := generateCase(rng)
		if err := runCase(bin, in, exp, n, m, r, b, xs, ys, cons); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
