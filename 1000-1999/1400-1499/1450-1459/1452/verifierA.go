package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func solve(x, y int) int {
	if x < y {
		x, y = y, x
	}
	if x == y {
		return 2 * x
	}
	return 2*x - 1
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("run error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	const t = 100
	xs := make([]int, t)
	ys := make([]int, t)
	exp := make([]int, t)
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", t)
	for i := 0; i < t; i++ {
		xs[i] = rand.Intn(10001)
		ys[i] = rand.Intn(10001)
		exp[i] = solve(xs[i], ys[i])
		fmt.Fprintf(&b, "%d %d\n", xs[i], ys[i])
	}
	out, err := runBinary(binary, b.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fields := strings.Fields(strings.TrimSpace(out))
	if len(fields) != t {
		fmt.Printf("expected %d lines, got %d\noutput:\n%s\n", t, len(fields), out)
		os.Exit(1)
	}
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil || v != exp[i] {
			fmt.Printf("test %d failed: input=(%d,%d) expected=%d got=%s\n", i+1, xs[i], ys[i], exp[i], f)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
