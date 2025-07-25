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

type Query struct {
	t int
	a int
	b int
	y string
}

type Test struct {
	n       int
	q       int
	s       []string
	queries []Query
}

func generateTests() []Test {
	rand.Seed(47)
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(3) + 2
		q := rand.Intn(3) + 1
		s := make([]string, n)
		for j := 0; j < n; j++ {
			s[j] = randString(rand.Intn(3) + 1)
		}
		queries := make([]Query, q)
		for k := 0; k < q; k++ {
			if rand.Intn(2) == 0 {
				a := rand.Intn(n) + 1
				b := rand.Intn(n-a+1) + a
				queries[k] = Query{t: 1, a: a, b: b}
			} else {
				x := rand.Intn(n) + 1
				y := randString(rand.Intn(3) + 1)
				queries[k] = Query{t: 2, a: x, y: y}
			}
		}
		tests = append(tests, Test{n: n, q: q, s: s, queries: queries})
	}
	tests = append(tests, Test{n: 2, q: 1, s: []string{"a", "b"}, queries: []Query{{t: 1, a: 1, b: 2}}})
	return tests
}

const letters = "abc"

func randString(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func lcp(a, b string) int {
	if len(a) > len(b) {
		a, b = b, a
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return len(a)
}

func maxArea(heights []int) int {
	n := len(heights)
	stack := make([]int, 0, n)
	maxA := 0
	for i := 0; i <= n; i++ {
		var h int
		if i < n {
			h = heights[i]
		} else {
			h = -1
		}
		for len(stack) > 0 && h < heights[stack[len(stack)-1]] {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			left := -1
			if len(stack) > 0 {
				left = stack[len(stack)-1]
			}
			width := i - left
			area := width * heights[top]
			if area > maxA {
				maxA = area
			}
		}
		stack = append(stack, i)
	}
	return maxA
}

func solveQuery(s []string, q Query) int {
	if q.t == 2 {
		s[q.a-1] = q.y
		return -1
	}
	a, b := q.a-1, q.b-1
	ans := 0
	for i := a; i <= b; i++ {
		if len(s[i]) > ans {
			ans = len(s[i])
		}
	}
	if a < b {
		heights := make([]int, b-a)
		for i := a; i < b; i++ {
			heights[i-a] = lcp(s[i], s[i+1])
		}
		area := maxArea(heights)
		if area > ans {
			ans = area
		}
	}
	return ans
}

func solve(t Test) []int {
	res := make([]int, 0, len(t.queries))
	s := append([]string(nil), t.s...)
	for _, qu := range t.queries {
		val := solveQuery(s, qu)
		if qu.t == 1 {
			res = append(res, val)
		}
	}
	return res
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
		fmt.Println("Usage: verifierF <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	passed := 0
	for i, t := range tests {
		input := fmt.Sprintf("%d %d\n", t.n, t.q)
		for _, str := range t.s {
			input += str + " "
		}
		input += "\n"
		for _, qu := range t.queries {
			if qu.t == 1 {
				input += fmt.Sprintf("1 %d %d\n", qu.a, qu.b)
			} else {
				input += fmt.Sprintf("2 %d %s\n", qu.a, qu.y)
			}
		}
		wants := solve(t)
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d exec err %v\n", i+1, err)
			continue
		}
		outLines := strings.Fields(strings.TrimSpace(output))
		if len(outLines) != len(wants) {
			fmt.Printf("Test %d expected %d outputs got %d\n", i+1, len(wants), len(outLines))
			continue
		}
		ok := true
		for k, w := range wants {
			g, err := strconv.Atoi(outLines[k])
			if err != nil || g != w {
				ok = false
				break
			}
		}
		if !ok {
			fmt.Printf("Test %d wrong output got %s expected %v\n", i+1, output, wants)
			continue
		}
		passed++
	}
	fmt.Printf("Passed %d/%d tests\n", passed, len(tests))
}
