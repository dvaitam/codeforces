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

const modD = 1000000007

func solveD(parents []int) []int {
	n := len(parents) + 1
	adj := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		p := parents[i-2]
		adj[p] = append(adj[p], i)
		adj[i] = append(adj[i], p)
	}
	dpDown := make([]int64, n+1)
	dpUp := make([]int64, n+1)
	ans := make([]int64, n+1)

	var dfs1 func(u, p int)
	dfs1 = func(u, p int) {
		prod := int64(1)
		for _, v := range adj[u] {
			if v == p {
				continue
			}
			dfs1(v, u)
			prod = prod * (dpDown[v] + 1) % modD
		}
		dpDown[u] = prod
	}
	var dfs2 func(u, p int)
	dfs2 = func(u, p int) {
		deg := len(adj[u])
		f := make([]int64, deg)
		for i, v := range adj[u] {
			if v == p {
				f[i] = dpUp[u] + 1
			} else {
				f[i] = dpDown[v] + 1
			}
		}
		pre := make([]int64, deg+1)
		suf := make([]int64, deg+1)
		pre[0] = 1
		for i := 0; i < deg; i++ {
			pre[i+1] = pre[i] * f[i] % modD
		}
		suf[deg] = 1
		for i := deg - 1; i >= 0; i-- {
			suf[i] = suf[i+1] * f[i] % modD
		}
		ans[u] = pre[deg]
		for i, v := range adj[u] {
			if v == p {
				continue
			}
			dpUp[v] = pre[i] * suf[i+1] % modD
			dfs2(v, u)
		}
	}

	dfs1(1, 0)
	dfs2(1, 0)
	res := make([]int, n)
	for i := 1; i <= n; i++ {
		res[i-1] = int(ans[i] % modD)
	}
	return res
}

func genTestD() (string, string) {
	n := rand.Intn(8) + 2 // at least 2
	parents := make([]int, n-1)
	for i := 2; i <= n; i++ {
		parents[i-2] = rand.Intn(i-1) + 1
	}
	input := fmt.Sprintf("%d\n", n)
	for i := 2; i <= n; i++ {
		if i > 2 {
			input += " "
		}
		input += fmt.Sprintf("%d", parents[i-2])
	}
	input += "\n"
	ans := solveD(parents)
	expected := ""
	for i, v := range ans {
		if i > 0 {
			expected += " "
		}
		expected += fmt.Sprintf("%d", v)
	}
	return input, expected
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 1; t <= 100; t++ {
		input, expected := genTestD()
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		if strings.TrimSpace(output) != expected {
			fmt.Printf("Test %d failed\nInput:\n%sExpected: %s\nGot: %s\n", t, input, expected, output)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
