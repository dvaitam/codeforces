package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(7) + 4 // 4..10
	m := rng.Intn(7) + 4
	x1 := rng.Intn(n-2) + 2
	y1 := rng.Intn(m-2) + 2
	var x2, y2 int
	for {
		x2 = rng.Intn(n-2) + 2
		y2 = rng.Intn(m-2) + 2
		if x2 != x1 && y2 != y1 {
			break
		}
	}
	return fmt.Sprintf("%d %d\n%d %d\n%d %d\n", n, m, x1, y1, x2, y2)
}

func verify(input, output string) error {
	in := bufio.NewReader(strings.NewReader(input))
	var n, m int
	var x1, y1, x2, y2 int
	fmt.Fscan(in, &n, &m)
	fmt.Fscan(in, &x1, &y1)
	fmt.Fscan(in, &x2, &y2)

	out := bufio.NewReader(strings.NewReader(output))
	var k int
	if _, err := fmt.Fscan(out, &k); err != nil {
		return fmt.Errorf("parse k: %v", err)
	}
	coords := make([][2]int, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(out, &coords[i][0], &coords[i][1]); err != nil {
			return fmt.Errorf("parse coord %d: %v", i+1, err)
		}
	}
	if _, err := fmt.Fscan(out, new(int)); err == nil {
		return fmt.Errorf("extra output data")
	}
	if k != n*m {
		return fmt.Errorf("length %d expected %d", k, n*m)
	}
	if coords[0][0] != x1 || coords[0][1] != y1 {
		return fmt.Errorf("path doesn't start at (%d,%d)", x1, y1)
	}
	if coords[k-1][0] != x2 || coords[k-1][1] != y2 {
		return fmt.Errorf("path doesn't end at (%d,%d)", x2, y2)
	}
	seen := make(map[[2]int]bool)
	prev := coords[0]
	seen[prev] = true
	for i := 1; i < k; i++ {
		cur := coords[i]
		if cur[0] < 1 || cur[0] > n || cur[1] < 1 || cur[1] > m {
			return fmt.Errorf("cell out of bounds: %v", cur)
		}
		dx := cur[0] - prev[0]
		if dx < 0 {
			dx = -dx
		}
		dy := cur[1] - prev[1]
		if dy < 0 {
			dy = -dy
		}
		if dx+dy != 1 {
			return fmt.Errorf("cells %v and %v not adjacent", prev, cur)
		}
		if seen[cur] {
			return fmt.Errorf("cell %v repeated", cur)
		}
		seen[cur] = true
		prev = cur
	}
	return nil
}

func runCase(exe, input string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return verify(input, out.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		if err := runCase(exe, in); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
