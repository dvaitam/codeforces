package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type edge struct{ to, id int }

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

func expected(n, m, k int, edges [][2]int, pairs [][2]int) int {
	adj := make([][]edge, n+1)
	for i, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], edge{v, i})
		adj[v] = append(adj[v], edge{u, i})
	}
	const logN = 17
	up := make([][]int, logN)
	for i := range up {
		up[i] = make([]int, n+1)
	}
	parentEdge := make([]int, n+1)
	depth := make([]int, n+1)
	var dfs func(int, int)
	dfs = func(v, p int) {
		for _, e := range adj[v] {
			if e.to == p {
				continue
			}
			up[0][e.to] = v
			parentEdge[e.to] = e.id
			depth[e.to] = depth[v] + 1
			for i := 1; i < logN; i++ {
				up[i][e.to] = up[i-1][up[i-1][e.to]]
			}
			dfs(e.to, v)
		}
	}
	dfs(1, 0)
	lca := func(a, b int) int {
		if depth[a] < depth[b] {
			a, b = b, a
		}
		diff := depth[a] - depth[b]
		for i := 0; i < logN; i++ {
			if diff&(1<<i) > 0 {
				a = up[i][a]
			}
		}
		if a == b {
			return a
		}
		for i := logN - 1; i >= 0; i-- {
			if up[i][a] != up[i][b] {
				a = up[i][a]
				b = up[i][b]
			}
		}
		return up[0][a]
	}
	getPath := func(u, v int) []int {
		p := lca(u, v)
		res := make([]int, 0)
		x := u
		for x != p {
			res = append(res, parentEdge[x])
			x = up[0][x]
		}
		tmp := make([]int, 0)
		x = v
		for x != p {
			tmp = append(tmp, parentEdge[x])
			x = up[0][x]
		}
		for i := len(tmp) - 1; i >= 0; i-- {
			res = append(res, tmp[i])
		}
		return res
	}
	edgeTrav := make([][]int, n-1)
	for idx, pr := range pairs {
		path := getPath(pr[0], pr[1])
		for _, e := range path {
			edgeTrav[e] = append(edgeTrav[e], idx)
		}
	}
	pairCount := make(map[[2]int]int)
	ans := 0
	for _, list := range edgeTrav {
		L := len(list)
		for i := 0; i < L; i++ {
			for j := i + 1; j < L; j++ {
				a := list[i]
				b := list[j]
				if a > b {
					a, b = b, a
				}
				key := [2]int{a, b}
				pairCount[key]++
				if pairCount[key] == k {
					ans++
				}
			}
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesF.txt")
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
		m, _ := strconv.Atoi(fields[1])
		k, _ := strconv.Atoi(fields[2])
		if len(fields) != 3+2*(n-1)+2*m {
			fmt.Fprintf(os.Stderr, "case %d invalid number of values\n", idx)
			os.Exit(1)
		}
		pos := 3
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			u, _ := strconv.Atoi(fields[pos])
			v, _ := strconv.Atoi(fields[pos+1])
			edges[i] = [2]int{u, v}
			pos += 2
		}
		pairs := make([][2]int, m)
		for i := 0; i < m; i++ {
			s, _ := strconv.Atoi(fields[pos])
			t, _ := strconv.Atoi(fields[pos+1])
			pairs[i] = [2]int{s, t}
			pos += 2
		}
		expect := expected(n, m, k, edges, pairs)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
		for _, e := range edges {
			input.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		for _, p := range pairs {
			input.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
		}
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		var ans int
		if _, err := fmt.Sscan(got, &ans); err != nil || ans != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
