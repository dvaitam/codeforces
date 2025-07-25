package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

// solveA implements the logic from 557A.go for a single test case.
func solveA(n, min1, max1, min2, max2, min3, max3 int) (int, int, int) {
	if n < 0 {
		return 0, 0, 0
	}
	x1 := max1
	if t := n - min2 - min3; t < x1 {
		x1 = t
	}
	if x1 < min1 {
		x1 = min1
	}
	remaining := n - x1
	x2 := max2
	if t := remaining - min3; t < x2 {
		x2 = t
	}
	if x2 < min2 {
		x2 = min2
	}
	x3 := n - x1 - x2
	return x1, x2, x3
}

func genCase() (string, string) {
	// generate random but small test case fulfilling constraints
	for {
		min1 := rand.Intn(10) + 1
		max1 := min1 + rand.Intn(10)
		min2 := rand.Intn(10) + 1
		max2 := min2 + rand.Intn(10)
		min3 := rand.Intn(10) + 1
		max3 := min3 + rand.Intn(10)
		minTotal := min1 + min2 + min3
		maxTotal := max1 + max2 + max3
		if minTotal > maxTotal {
			continue
		}
		n := rand.Intn(maxTotal-minTotal+1) + minTotal
		if n < 3 {
			continue
		}
		x1, x2, x3 := solveA(n, min1, max1, min2, max2, min3, max3)
		input := fmt.Sprintf("%d\n%d %d\n%d %d\n%d %d\n", n, min1, max1, min2, max2, min3, max3)
		output := fmt.Sprintf("%d %d %d", x1, x2, x3)
		return input, output
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(1)
	for i := 0; i < 100; i++ {
		input, expected := genCase()
		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}
