package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseA struct{ x, y int }

func genTestsA() []testCaseA {
	rng := rand.New(rand.NewSource(42))
	tests := []testCaseA{
		{0, 0}, {3, 4}, {5, 12}, {1, 0}, {0, 1},
	}
	for len(tests) < 100 {
		tests = append(tests, testCaseA{rng.Intn(51), rng.Intn(51)})
	}
	return tests
}

func solveA(tc testCaseA) string {
	if tc.x == 0 && tc.y == 0 {
		return "0"
	}
	d2 := tc.x*tc.x + tc.y*tc.y
	r := int(math.Round(math.Sqrt(float64(d2))))
	if r*r == d2 {
		return "1"
	}
	return "2"
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
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	tests := genTestsA()
	for i, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.x, tc.y)
		exp := solveA(tc)
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
