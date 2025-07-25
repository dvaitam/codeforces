package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type edge struct{ to int }

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildTree(n int, parents []int) [][]int {
	g := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		p := parents[i]
		g[p] = append(g[p], i)
		g[i] = append(g[i], p)
	}
	return g
}

func subtreeNodes(g [][]int, root, parent int, mark []bool) []int {
	mark[root] = true
	nodes := []int{root}
	for _, v := range g[root] {
		if v != parent {
			nodes = append(nodes, subtreeNodes(g, v, root, mark)...)
		}
	}
	return nodes
}

func componentSize(g [][]int, banned int, nodesMap map[int]bool, start int, visited map[int]bool) int {
	stack := []int{start}
	visited[start] = true
	size := 0
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		size++
		for _, v := range g[u] {
			if v == banned || !nodesMap[v] || visited[v] {
				continue
			}
			visited[v] = true
			stack = append(stack, v)
		}
	}
	return size
}

func centroidOfSubtree(g [][]int, root int) int {
	mark := make([]bool, len(g))
	nodes := subtreeNodes(g, root, 0, mark)
	nodesMap := make(map[int]bool)
	for _, v := range nodes {
		nodesMap[v] = true
	}
	total := len(nodes)
	for _, c := range nodes {
		visited := make(map[int]bool)
		maxComp := 0
		for _, v := range nodes {
			if v == c || visited[v] {
				continue
			}
			size := componentSize(g, c, nodesMap, v, visited)
			if size > maxComp {
				maxComp = size
			}
		}
		if maxComp*2 <= total {
			return c
		}
	}
	return root
}

func expected(n int, parents []int, queries []int) []int {
	g := buildTree(n, parents)
	res := make([]int, len(queries))
	for i, v := range queries {
		res[i] = centroidOfSubtree(g, v)
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go <binary>")
		return
	}
	bin := os.Args[1]
	rand.Seed(42)
	for t := 0; t < 100; t++ {
		n := rand.Intn(10) + 1
		q := rand.Intn(10) + 1
		parents := make([]int, n+1)
		for i := 2; i <= n; i++ {
			parents[i] = rand.Intn(i-1) + 1
		}
		queries := make([]int, q)
		for i := 0; i < q; i++ {
			queries[i] = rand.Intn(n) + 1
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
		for i := 2; i <= n; i++ {
			if i > 2 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(parents[i]))
		}
		if n > 1 {
			sb.WriteByte('\n')
		} else {
			sb.WriteString("\n")
		}
		for i := 0; i < q; i++ {
			sb.WriteString(fmt.Sprintf("%d\n", queries[i]))
		}
		want := expected(n, parents, queries)
		out, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\noutput:%s", t+1, err, out)
			return
		}
		outLines := strings.Fields(out)
		if len(outLines) != q {
			fmt.Printf("invalid output on test %d\ninput:%soutput:%s", t+1, sb.String(), out)
			return
		}
		for i := 0; i < q; i++ {
			got, err := strconv.Atoi(outLines[i])
			if err != nil || got != want[i] {
				fmt.Printf("wrong answer on test %d\ninput:%sexpected:%v\noutput:%s", t+1, sb.String(), want, out)
				return
			}
		}
	}
	fmt.Println("All tests passed")
}
