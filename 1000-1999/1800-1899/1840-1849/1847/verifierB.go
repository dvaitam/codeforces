package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solveB(n int, a []int) int {
	all := a[0]
	for i := 1; i < n; i++ {
		all &= a[i]
	}
	if all > 0 {
		return 1
	}
	cur := -1
	cnt := 0
	for i := 0; i < n; i++ {
		cur &= a[i]
		if cur == 0 {
			cnt++
			cur = -1
		}
	}
	return cnt
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(2)
	const t = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(50) + 1
		a := make([]int, n)
		for j := range a {
			a[j] = rand.Intn(1000)
		}
		fmt.Fprintln(&input, n)
		for j := 0; j < n; j++ {
			if j+1 == n {
				fmt.Fprintf(&input, "%d\n", a[j])
			} else {
				fmt.Fprintf(&input, "%d ", a[j])
			}
		}
		expected[i] = strconv.Itoa(solveB(n, a))
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
