package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runExe(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return edges
}

func computeMatrix(n int, edges [][2]int) []string {
	g := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	res := make([]string, n)
	for root := 1; root <= n; root++ {
		parent := make([]int, n+1)
		order := make([]int, 0, n)
		stack := []int{root}
		parent[root] = -1
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			order = append(order, v)
			for _, to := range g[v] {
				if to == parent[v] {
					continue
				}
				parent[to] = v
				stack = append(stack, to)
			}
		}
		win := make([]bool, n+1)
		for i := len(order) - 1; i >= 0; i-- {
			v := order[i]
			good := false
			for _, to := range g[v] {
				if to == parent[v] {
					continue
				}
				if !win[to] {
					good = true
				}
			}
			win[v] = good
		}
		var sb strings.Builder
		for start := 1; start <= n; start++ {
			if win[start] {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('0')
			}
		}
		res[root-1] = sb.String()
	}
	return res
}

func genCase(rng *rand.Rand) (string, []string) {
	n := rng.Intn(4) + 2
	edges := generateTree(rng, n)
	matrix := computeMatrix(n, edges)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		sb.WriteString(matrix[i])
		sb.WriteByte('\n')
	}
	return sb.String(), matrix
}

func parseEdges(out string, n int) ([][2]int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("no output")
	}
	if strings.ToUpper(strings.TrimSpace(lines[0])) != "YES" {
		return nil, fmt.Errorf("expected YES")
	}
	if len(lines)-1 != n-1 {
		return nil, fmt.Errorf("expected %d edges", n-1)
	}
	edges := make([][2]int, 0, n-1)
	for i := 1; i < len(lines); i++ {
		fields := strings.Fields(lines[i])
		if len(fields) < 2 {
			return nil, fmt.Errorf("edge format")
		}
		u, err1 := strconv.Atoi(fields[0])
		v, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil || u < 1 || u > n || v < 1 || v > n {
			return nil, fmt.Errorf("invalid edge")
		}
		edges = append(edges, [2]int{u, v})
	}
	return edges, nil
}

func isTree(n int, edges [][2]int) bool {
	if len(edges) != n-1 {
		return false
	}
	g := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	vis := make([]bool, n+1)
	stack := []int{1}
	vis[1] = true
	count := 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		count++
		for _, to := range g[v] {
			if !vis[to] {
				vis[to] = true
				stack = append(stack, to)
			}
		}
	}
	return count == n
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, matrix := genCase(rng)
		out, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		nLines := strings.Count(input, "\n") - 1
		candEdges, err := parseEdges(out, nLines)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, input, out)
			os.Exit(1)
		}
		if !isTree(nLines, candEdges) {
			fmt.Printf("case %d failed: output is not a tree\ninput:\n%soutput:\n%s\n", i+1, input, out)
			os.Exit(1)
		}
		want := computeMatrix(nLines, candEdges)
		valid := true
		for j := 0; j < nLines; j++ {
			if matrix[j] != want[j] {
				valid = false
				break
			}
		}
		if !valid {
			fmt.Printf("case %d failed: tree does not match matrix\ninput:\n%soutput:\n%s\n", i+1, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
