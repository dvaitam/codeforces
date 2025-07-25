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
	s string
}

func generateTests() []Test {
	rand.Seed(45)
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(10) + 2
		for {
			b := make([]byte, n)
			has0, has1 := false, false
			for j := 0; j < n; j++ {
				if rand.Intn(2) == 0 {
					b[j] = '0'
					has0 = true
				} else {
					b[j] = '1'
					has1 = true
				}
			}
			if has0 && has1 {
				tests = append(tests, Test{s: string(b)})
				break
			}
		}
	}
	tests = append(tests, Test{s: "01"})
	tests = append(tests, Test{s: "10"})
	return tests
}

func solve(t Test) (int, int) {
	pos0, pos1 := -1, -1
	for i, ch := range t.s {
		if ch == '0' && pos0 == -1 {
			pos0 = i + 1
		}
		if ch == '1' && pos1 == -1 {
			pos1 = i + 1
		}
	}
	return pos0, pos1
}

func runBinary(bin, input string) (string, error) {
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
		fmt.Println("Usage: verifierD <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	passed := 0
	for i, t := range tests {
		input := fmt.Sprintf("%d\n%s\n", len(t.s), t.s)
		want0, want1 := solve(t)
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d exec err %v\n", i+1, err)
			continue
		}
		outStr := strings.TrimSpace(output)
		nums := strings.Fields(outStr)
		if len(nums) != 2 {
			fmt.Printf("Test %d bad output %s\n", i+1, outStr)
			continue
		}
		g0, err0 := strconv.Atoi(nums[0])
		g1, err1 := strconv.Atoi(nums[1])
		if err0 != nil || err1 != nil || g0 != want0 || g1 != want1 {
			fmt.Printf("Test %d expected %d %d got %s\n", i+1, want0, want1, outStr)
			continue
		}
		passed++
	}
	fmt.Printf("Passed %d/%d tests\n", passed, len(tests))
}
