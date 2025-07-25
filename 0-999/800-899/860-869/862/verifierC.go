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
	n int
	x int
}

func generateTests() []Test {
	rand.Seed(44)
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(6) + 1
		x := rand.Intn(20)
		tests = append(tests, Test{n: n, x: x})
	}
	tests = append(tests, Test{n: 2, x: 0})
	return tests
}

func hasSolution(n, x int) bool {
	return !(n == 2 && x == 0)
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

func parseInts(s string) ([]int, error) {
	fields := strings.Fields(s)
	res := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

func xorAll(arr []int) int {
	x := 0
	for _, v := range arr {
		x ^= v
	}
	return x
}

func allDistinct(arr []int) bool {
	m := make(map[int]bool)
	for _, v := range arr {
		if m[v] {
			return false
		}
		m[v] = true
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: verifierC <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	passed := 0
	for i, t := range tests {
		input := fmt.Sprintf("%d %d\n", t.n, t.x)
		wantExists := hasSolution(t.n, t.x)
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d: exec error %v\n", i+1, err)
			continue
		}
		lines := strings.Split(strings.TrimSpace(output), "\n")
		if !wantExists {
			if strings.TrimSpace(lines[0]) != "NO" {
				fmt.Printf("Test %d: expected NO got %s\n", i+1, lines[0])
				continue
			}
			passed++
			continue
		}
		if strings.TrimSpace(lines[0]) != "YES" {
			fmt.Printf("Test %d: expected YES got %s\n", i+1, lines[0])
			continue
		}
		numsLine := ""
		if len(lines) > 1 {
			numsLine = strings.Join(lines[1:], " ")
		}
		nums, err := parseInts(numsLine)
		if err != nil || len(nums) != t.n || !allDistinct(nums) || xorAll(nums) != t.x {
			fmt.Printf("Test %d: invalid set %v\n", i+1, nums)
			continue
		}
		ok := true
		for _, v := range nums {
			if v < 0 || v > 1000000 {
				ok = false
				break
			}
		}
		if !ok {
			fmt.Printf("Test %d: numbers out of range\n", i+1)
			continue
		}
		passed++
	}
	fmt.Printf("Passed %d/%d tests\n", passed, len(tests))
}
