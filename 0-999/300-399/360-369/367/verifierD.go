package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func checkSubset(mask int, sets [][]int, n, d int) bool {
	var arr []int
	for i, set := range sets {
		if mask&(1<<uint(i)) != 0 {
			arr = append(arr, set...)
		}
	}
	if len(arr) == 0 {
		return false
	}
	sort.Ints(arr)
	if arr[0] > d {
		return false
	}
	for i := 0; i < len(arr)-1; i++ {
		if arr[i+1]-arr[i] > d {
			return false
		}
	}
	if arr[len(arr)-1] < n-d+1 {
		return false
	}
	return true
}

func solveD(n, m, d int, sets [][]int) int {
	best := m + 1
	total := 1 << uint(m)
	for mask := 1; mask < total; mask++ {
		if checkSubset(mask, sets, n, d) {
			c := bitsCount(mask)
			if c < best {
				best = c
			}
		}
	}
	if best == m+1 {
		best = m
	}
	return best
}

func bitsCount(x int) int {
	c := 0
	for x > 0 {
		x &= x - 1
		c++
	}
	return c
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	m := rng.Intn(min(5, n)) + 1
	d := rng.Intn(n) + 1
	// assign each number 1..n to one of m sets randomly
	sets := make([][]int, m)
	for v := 1; v <= n; v++ {
		idx := rng.Intn(m)
		sets[idx] = append(sets[idx], v)
	}
	input := fmt.Sprintf("%d %d %d\n", n, m, d)
	for i := 0; i < m; i++ {
		sort.Ints(sets[i])
		input += fmt.Sprintf("%d", len(sets[i]))
		for _, v := range sets[i] {
			input += fmt.Sprintf(" %d", v)
		}
		input += "\n"
	}
	ans := solveD(n, m, d, sets)
	expected := fmt.Sprintf("%d", ans)
	return input, expected
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
	var stderr bytes.Buffer
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
		input, expect := generateCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
