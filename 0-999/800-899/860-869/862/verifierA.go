package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Test struct {
	n   int
	x   int
	arr []int
}

func generateTests() []Test {
	rand.Seed(42)
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(10) + 1
		x := rand.Intn(10)
		vals := rand.Perm(20)[:n]
		tests = append(tests, Test{n: n, x: x, arr: vals})
	}
	// edge cases
	tests = append(tests, Test{n: 1, x: 0, arr: []int{1}})
	tests = append(tests, Test{n: 1, x: 5, arr: []int{0}})
	tests = append(tests, Test{n: 3, x: 2, arr: []int{0, 3, 4}})
	return tests
}

func solve(t Test) int {
	present := make(map[int]bool)
	for _, v := range t.arr {
		present[v] = true
	}
	ops := 0
	for i := 0; i < t.x; i++ {
		if !present[i] {
			ops++
		}
	}
	if present[t.x] {
		ops++
	}
	return ops
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
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: verifierA <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	passed := 0
	for i, t := range tests {
		input := fmt.Sprintf("%d %d\n", t.n, t.x)
		for j, v := range t.arr {
			if j > 0 {
				input += " "
			}
			input += strconv.Itoa(v)
		}
		input += "\n"
		want := solve(t)
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d: execution error: %v\n", i+1, err)
			continue
		}
		outStr := strings.TrimSpace(output)
		got, err := strconv.Atoi(outStr)
		if err != nil || got != want {
			fmt.Printf("Test %d: expected %d got %s\n", i+1, want, outStr)
		} else {
			passed++
		}
	}
	fmt.Printf("Passed %d/%d tests\n", passed, len(tests))
}
