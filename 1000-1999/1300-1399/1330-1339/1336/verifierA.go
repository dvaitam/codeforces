package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n, k int, edges [][2]int) int {
	g := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	depth := make([]int, n+1)
	parent := make([]int, n+1)
	order := make([]int, 0, n)
	stack := []int{1}
	parent[1] = 0
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, v := range g[u] {
			if v == parent[u] {
				continue
			}
			parent[v] = u
			depth[v] = depth[u] + 1
			stack = append(stack, v)
		}
	}
	size := make([]int, n+1)
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		size[u] = 1
		for _, v := range g[u] {
			if v == parent[u] {
				continue
			}
			size[u] += size[v]
		}
	}
	vals := make([]int, n)
	for i := 1; i <= n; i++ {
		vals[i-1] = depth[i] - (size[i] - 1)
	}
	sort.Slice(vals, func(i, j int) bool { return vals[i] > vals[j] })
	ans := 0
	for i := 0; i < k; i++ {
		ans += vals[i]
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 3 {
			fmt.Fprintf(os.Stderr, "case %d invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		k, _ := strconv.Atoi(fields[1])
		expectEdges := n - 1
		if len(fields) != 2+2*expectEdges {
			fmt.Fprintf(os.Stderr, "case %d invalid number of values\n", idx)
			os.Exit(1)
		}
		edges := make([][2]int, expectEdges)
		pos := 2
		for i := 0; i < expectEdges; i++ {
			u, _ := strconv.Atoi(fields[pos])
			v, _ := strconv.Atoi(fields[pos+1])
			edges[i] = [2]int{u, v}
			pos += 2
		}
		expectedAns := expected(n, k, edges)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for _, e := range edges {
			input.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		var ans int
		if _, err := fmt.Sscan(got, &ans); err != nil || ans != expectedAns {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", idx, expectedAns, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
