package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveA(p []int) int {
	n := len(p)
	depth := make([]int, n)
	var dfs func(int) int
	dfs = func(u int) int {
		if depth[u] != 0 {
			return depth[u]
		}
		if p[u] == -1 {
			depth[u] = 1
		} else {
			depth[u] = dfs(p[u]) + 1
		}
		return depth[u]
	}
	maxD := 0
	for i := 0; i < n; i++ {
		d := dfs(i)
		if d > maxD {
			maxD = d
		}
	}
	return maxD
}

func genCase() (string, int) {
	n := rand.Intn(20) + 1
	parents := make([]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i == 0 || rand.Intn(2) == 0 {
			parents[i] = -1
		} else {
			parents[i] = rand.Intn(i)
		}
		if parents[i] == -1 {
			sb.WriteString("-1")
		} else {
			sb.WriteString(fmt.Sprintf("%d", parents[i]+1))
		}
		if i+1 < n {
			sb.WriteByte('\n')
		}
	}
	expect := solveA(parents)
	return sb.String(), expect
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(42)
	for t := 1; t <= 100; t++ {
		in, expect := genCase()
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", t, err)
			fmt.Println(in)
			return
		}
		if strings.TrimSpace(got) != fmt.Sprint(expect) {
			fmt.Printf("Test %d failed\nInput:\n%s\nExpected: %d\nGot: %s\n", t, in, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
