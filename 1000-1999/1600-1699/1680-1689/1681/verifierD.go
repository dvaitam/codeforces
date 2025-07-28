package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type TestCase struct {
	input  string
	output string
}

type node struct {
	val   int64
	steps int
}

func digitsLen(x int64) int {
	if x == 0 {
		return 1
	}
	l := 0
	for x > 0 {
		x /= 10
		l++
	}
	return l
}

func uniqueDigits(x int64) []int {
	seen := [10]bool{}
	res := make([]int, 0, 10)
	if x == 0 {
		return []int{0}
	}
	for x > 0 {
		d := int(x % 10)
		if !seen[d] {
			seen[d] = true
			res = append(res, d)
		}
		x /= 10
	}
	return res
}

func bfs(n int, x int64) int {
	limit := int64(1)
	for i := 0; i < n; i++ {
		limit *= 10
	}
	queue := []node{{x, 0}}
	visited := map[int64]bool{x: true}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if digitsLen(cur.val) == n {
			return cur.steps
		}
		digits := uniqueDigits(cur.val)
		for _, d := range digits {
			if d <= 1 {
				continue
			}
			next := cur.val * int64(d)
			if next >= limit {
				continue
			}
			if !visited[next] {
				visited[next] = true
				queue = append(queue, node{next, cur.steps + 1})
			}
		}
	}
	return -1
}

func solveCaseD(n int, x int64) string {
	return fmt.Sprintf("%d\n", bfs(n, x))
}

func generateTests() []TestCase {
	rand.Seed(4)
	tests := make([]TestCase, 0, 20)
	for t := 0; t < 20; t++ {
		n := rand.Intn(4) + 1
		var x int64
		if n == 1 {
			x = int64(rand.Intn(9) + 1)
		} else {
			start := int64(rand.Intn(9) + 1)
			x = start
		}
		input := fmt.Sprintf("%d %d\n", n, x)
		out := solveCaseD(n, x)
		tests = append(tests, TestCase{input, out})
	}
	return tests
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	passed := 0
	for i, tc := range tests {
		got, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			continue
		}
		g := strings.TrimSpace(got)
		e := strings.TrimSpace(tc.output)
		if g != e {
			fmt.Printf("Test %d failed. Expected %s got %s\n", i+1, e, g)
		} else {
			passed++
		}
	}
	fmt.Printf("%d/%d tests passed\n", passed, len(tests))
}
