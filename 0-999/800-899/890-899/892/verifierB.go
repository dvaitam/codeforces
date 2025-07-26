package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	n int
	l []int
}

func generateTest(i int) testCase {
	r := rand.New(rand.NewSource(int64(i + 1000)))
	n := r.Intn(50) + 1
	l := make([]int, n)
	for j := 0; j < n; j++ {
		l[j] = r.Intn(1000)
	}
	return testCase{n: n, l: l}
}

func buildTests() []testCase {
	tests := make([]testCase, 100)
	for i := range tests {
		tests[i] = generateTest(i)
	}
	return tests
}

func expected(t testCase) int {
	n := t.n
	left := n + 1
	alive := 0
	for i := n; i >= 1; i-- {
		if i < left {
			alive++
		}
		if v := i - t.l[i-1]; v < left {
			left = v
		}
	}
	return alive
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), fmt.Errorf("timeout")
	}
	if err != nil {
		return out.String(), err
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierB <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := buildTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d\n", t.n)
		for j, v := range t.l {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprint(v)
		}
		input += "\n"
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nOutput:%s\n", i+1, err, out)
			os.Exit(1)
		}
		exp := fmt.Sprint(expected(t))
		got := strings.TrimSpace(out)
		if got != exp {
			fmt.Printf("Test %d failed: expected %s got %s\n", i+1, exp, got)
			fmt.Printf("Input:\n%s\n", input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
