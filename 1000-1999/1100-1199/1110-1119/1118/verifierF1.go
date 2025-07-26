package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func solve(n int, colors []int, edges [][2]int) int {
	totalRed := 0
	totalBlue := 0
	for _, c := range colors {
		if c == 1 {
			totalRed++
		} else if c == 2 {
			totalBlue++
		}
	}
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	parent := make([]int, n)
	parent[0] = -1
	stack := []int{0}
	order := make([]int, 0, n)
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			parent[v] = u
			stack = append(stack, v)
		}
	}
	redCnt := make([]int, n)
	blueCnt := make([]int, n)
	ans := 0
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		for _, v := range adj[u] {
			if parent[v] == u {
				redCnt[u] += redCnt[v]
				blueCnt[u] += blueCnt[v]
			}
		}
		if colors[u] == 1 {
			redCnt[u]++
		} else if colors[u] == 2 {
			blueCnt[u]++
		}
		if u != 0 {
			if redCnt[u] == 0 && blueCnt[u] == totalBlue {
				ans++
			} else if blueCnt[u] == 0 && redCnt[u] == totalRed {
				ans++
			}
		}
	}
	return ans
}

func randomTree(n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		p := rand.Intn(i)
		edges = append(edges, [2]int{i, p})
	}
	return edges
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF1.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(7)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(8) + 2
		colors := make([]int, n)
		hasR := false
		hasB := false
		for i := 0; i < n; i++ {
			colors[i] = rand.Intn(3)
			if colors[i] == 1 {
				hasR = true
			} else if colors[i] == 2 {
				hasB = true
			}
		}
		if !hasR {
			colors[0] = 1
			hasR = true
		}
		if !hasB {
			colors[1%n] = 2
			hasB = true
		}
		edges := randomTree(n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("%d ", colors[i]))
		}
		sb.WriteString("\n")
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
		}
		expected := solve(n, colors, edges)
		out, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", t, err)
			fmt.Println(out)
			return
		}
		var got int
		fmt.Fscan(strings.NewReader(out), &got)
		if got != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected %d got %d\n", t, sb.String(), expected, got)
			return
		}
	}
	fmt.Println("all tests passed")
}
