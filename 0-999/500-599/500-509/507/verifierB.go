package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type TestB struct {
	r  int64
	x1 int64
	y1 int64
	x2 int64
	y2 int64
}

func generateTests() []TestB {
	rand.Seed(2)
	tests := make([]TestB, 100)
	for i := range tests {
		r := rand.Int63n(1000) + 1
		x1 := rand.Int63n(2001) - 1000
		y1 := rand.Int63n(2001) - 1000
		x2 := rand.Int63n(2001) - 1000
		y2 := rand.Int63n(2001) - 1000
		tests[i] = TestB{r, x1, y1, x2, y2}
	}
	return tests
}

func expected(t TestB) int64 {
	dx := float64(t.x1 - t.x2)
	dy := float64(t.y1 - t.y2)
	dist := math.Hypot(dx, dy)
	maxMove := 2 * float64(t.r)
	return int64(math.Ceil(dist / maxMove))
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d %d %d %d %d\n", t.r, t.x1, t.y1, t.x2, t.y2)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(strings.TrimSpace(got), 10, 64)
		if err != nil {
			fmt.Printf("test %d: invalid output\n", i+1)
			os.Exit(1)
		}
		exp := expected(t)
		if val != exp {
			fmt.Printf("test %d: expected %d got %d\n", i+1, exp, val)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
