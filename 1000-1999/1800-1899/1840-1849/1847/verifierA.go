package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func solveA(n, k int, a []int) int {
	if n == 1 {
		return 0
	}
	diffs := make([]int, n-1)
	total := 0
	for i := 0; i < n-1; i++ {
		d := a[i+1] - a[i]
		if d < 0 {
			d = -d
		}
		diffs[i] = d
		total += d
	}
	sort.Slice(diffs, func(i, j int) bool { return diffs[i] > diffs[j] })
	remove := 0
	for i := 0; i < k-1 && i < len(diffs); i++ {
		remove += diffs[i]
	}
	return total - remove
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(1)
	const t = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(100) + 1
		k := rand.Intn(n) + 1
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rand.Intn(500) + 1
		}
		fmt.Fprintf(&input, "%d %d\n", n, k)
		for j := 0; j < n; j++ {
			if j+1 == n {
				fmt.Fprintf(&input, "%d\n", a[j])
			} else {
				fmt.Fprintf(&input, "%d ", a[j])
			}
		}
		expected[i] = strconv.Itoa(solveA(n, k, a))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error executing binary: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Fields(string(bytes.TrimSpace(out)))
	if len(lines) != t {
		fmt.Fprintf(os.Stderr, "expected %d lines of output, got %d\n", t, len(lines))
		fmt.Fprint(os.Stderr, string(out))
		os.Exit(1)
	}
	for i := 0; i < t; i++ {
		if lines[i] != expected[i] {
			fmt.Fprintf(os.Stderr, "mismatch on case %d: expected %s got %s\n", i+1, expected[i], lines[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
