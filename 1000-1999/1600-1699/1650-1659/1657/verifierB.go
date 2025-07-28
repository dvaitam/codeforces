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

// test case for problem B

type testCaseB struct {
	n int
	B int64
	x int64
	y int64
}

func genTestsB() []testCaseB {
	rng := rand.New(rand.NewSource(43))
	tests := []testCaseB{
		{1, 1, 1, 1}, {3, 5, 2, 1}, {5, 10, 3, 4},
	}
	for len(tests) < 100 {
		n := rng.Intn(20) + 1
		B := rng.Int63n(100) + 1
		x := rng.Int63n(50) + 1
		y := rng.Int63n(50) + 1
		tests = append(tests, testCaseB{n, int64(B), x, y})
	}
	return tests
}

func solveB(tc testCaseB) string {
	cur := int64(0)
	sum := int64(0)
	for i := 0; i < tc.n; i++ {
		if cur+tc.x <= tc.B {
			cur += tc.x
		} else {
			cur -= tc.y
		}
		sum += cur
	}
	return fmt.Sprint(sum)
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	tests := genTestsB()
	for i, tc := range tests {
		input := fmt.Sprintf("%d %d %d %d\n", tc.n, tc.B, tc.x, tc.y)
		exp := solveB(tc)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != exp {
			fmt.Printf("test %d failed: expected %q got %q\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
