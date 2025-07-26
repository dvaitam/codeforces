package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	n int
	a []int
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(1))
	tests := make([]testCase, 100)
	for i := range tests {
		n := rnd.Intn(10) + 2
		a := make([]int, n-1)
		for j := range a {
			a[j] = rnd.Intn(n) + 1
		}
		tests[i] = testCase{n, a}
	}
	return tests
}

// expected using algorithm from 1283F.go
func expected(tc testCase) (int, [][2]int) {
	n := tc.n
	v := append([]int{0}, tc.a...)
	root := v[1]
	findPath := make([]bool, n+1)
	findPath[root] = true
	p := 2
	last := n
	ans := make([][2]int, 0, n-1)
	for {
		for last > 0 && findPath[last] {
			last--
		}
		if last <= 0 {
			break
		}
		findPath[last] = true
		for p < n && !findPath[v[p]] {
			ans = append(ans, [2]int{v[p-1], v[p]})
			findPath[v[p]] = true
			p++
		}
		ans = append(ans, [2]int{last, v[p-1]})
		p++
	}
	return root, ans
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
		return out.String(), fmt.Errorf("exec error: %v", err)
	}
	return out.String(), nil
}

func parseEdges(lines []string) (int, [][2]int, error) {
	var root int
	if _, err := fmt.Sscan(lines[0], &root); err != nil {
		return 0, nil, fmt.Errorf("invalid root")
	}
	edges := make([][2]int, len(lines)-1)
	for i := 1; i < len(lines); i++ {
		if _, err := fmt.Sscan(lines[i], &edges[i-1][0], &edges[i-1][1]); err != nil {
			return 0, nil, fmt.Errorf("bad edge")
		}
	}
	return root, edges, nil
}

func normalize(edges [][2]int) [][2]int {
	out := make([][2]int, len(edges))
	for i, e := range edges {
		if e[0] > e[1] {
			e[0], e[1] = e[1], e[0]
		}
		out[i] = e
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i][0] == out[j][0] {
			return out[i][1] < out[j][1]
		}
		return out[i][0] < out[j][0]
	})
	return out
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for idx, tc := range tests {
		input := fmt.Sprintf("%d\n", tc.n)
		for i, x := range tc.a {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", x)
		}
		input += "\n"
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		rootGot, edgesGot, err := parseEdges(lines)
		if err != nil || len(edgesGot) != tc.n-1 {
			fmt.Printf("test %d: parse fail\n", idx+1)
			os.Exit(1)
		}
		rootExp, edgesExp := expected(tc)
		if rootGot != rootExp {
			fmt.Printf("test %d: expected root %d got %d\n", idx+1, rootExp, rootGot)
			os.Exit(1)
		}
		ne := normalize(edgesExp)
		ng := normalize(edgesGot)
		if len(ne) != len(ng) {
			fmt.Printf("test %d: edge count mismatch\n", idx+1)
			os.Exit(1)
		}
		for i := 0; i < len(ne); i++ {
			if ne[i] != ng[i] {
				fmt.Printf("test %d: edges mismatch\n", idx+1)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
