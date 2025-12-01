package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesE.txt.
const embeddedTestcasesE = `2 3 0 1 1 0 0 0
1 4 1 0 0 1
3 2 0 1 0 0 1 1
2 2 1 1 1 0
5 3 1 0 0 0 1 1 0 1 0 1 1 0 1 1 1
1 1 1
2 4 0 1 0 1 1 1 1 0
2 2 0 1 0 1
1 3 1 0 1
2 3 0 1 1 0 0 1
5 6 0 1 0 1 0 0 0 1 0 0 0 0 1 1 1 0 1 0 0 1 1 1 1 1 0 1 0 1 1 1
3 3 1 1 1 0 1 0 1 1 1
1 5 0 0 0 0 1
2 2 1 0 1 1
4 4 1 1 0 1 1 1 1 1 0 0 0 0 1 1 1 0
5 6 1 0 0 0 0 1 1 0 0 0 0 1 1 0 0 0 0 1 1 0 0 0 0 1 0 0 0 1 0 0 0 0
2 2 1 1 1 1
3 6 1 1 1 1 1 1 1 1 1 1 1 1 1 0 1 1 0 1
2 2 1 1 1 1
1 3 1 1 1
3 5 0 0 1 1 1 0 0 0 0 1 1 0 0 1 1
2 5 1 0 1 1 1 0 0 1 0 0
1 5 1 1 1 1 0
1 5 1 0 1 1 0
5 4 1 0 0 1 0 1 1 1 0 0 0 1 1 0 0 1 1 1 0 1
4 4 1 1 1 0 1 1 1 0 1 1 0 1 0 1 1 0
3 5 1 0 0 0 0 0 1 0 1 0 0 0 0 0 0
2 3 0 0 1 1 1 1
3 5 0 1 1 0 0 1 1 1 1 1 1 1 1 0 0
2 3 1 1 0 1 1 1
5 6 0 0 0 1 0 0 1 0 1 1 1 1 0 0 0 0 0 0 1 1 1 1 1 0 1 0 1 0 0 0
5 6 0 0 1 1 1 1 1 1 0 0 0 1 0 1 1 0 0 0 1 1 0 0 1 0 1 0 1 1 1 1
3 5 0 0 1 0 1 0 1 0 1 0 1 0 0 0 0
2 6 1 0 1 1 0 0 0 0 1 0 0 0
2 2 1 0 0 1
3 3 0 0 1 1 1 1 1 1 1
2 3 1 0 0 1 0 1
2 4 0 0 0 0 0 0 1 1
5 3 0 0 0 1 0 1 0 1 0 0 0 0 1 1 0
3 3 0 1 0 0 0 0 0 1 0
3 4 0 0 0 0 0 0 0 0 0 0 0 0
1 5 0 1 1 1 0
1 2 1 0
1 1 1
2 2 0 1 1 0
2 3 0 0 0 1 0 0
2 4 1 0 1 1 1 0 1 1
2 3 1 1 0 1 1 0
1 6 1 1 1 0 1 0
3 5 0 0 1 1 1 0 1 0 0 1 1 1 0 1 0
1 6 0 1 1 0 1 0
2 6 0 1 1 1 1 1 0 1 1 1 1 0
2 2 1 0 0 0
2 5 0 0 1 1 0 1 1 1 0 0
2 2 1 1 1 1
1 2 0 0
1 2 1 1
4 5 1 1 1 0 1 1 0 0 0 1 1 0 1 1 1 1 1 0 1 0 0 0 0
3 5 0 1 0 0 1 1 0 0 0 0 1 1 0 1 0
1 5 0 1 0 1 1
3 6 0 0 0 1 1 0 1 1 0 0 0 0 0 1 1 0 0 1 0
2 4 0 1 0 0 0 0 0 0
3 6 0 1 1 1 0 0 0 1 0 0 1 0 0 0 1 0 1 1
2 4 0 0 0 1 0 0 1 1
5 5 0 0 0 0 1 1 0 0 1 0 1 0 0 0 1 1 1 1 0 1 0 0 1 1 0
4 6 1 1 1 1 1 0 1 1 0 0 0 1 0 0 1 1 1 0 1 1 0 1 1 1
2 5 1 1 0 1 0 1 0 0 0
1 5 0 1 0 1 1
2 2 1 1 1 1
1 5 0 1 0 0 1
5 6 1 1 1 0 0 1 0 0 0 0 1 0 1 0 0 0 0 1 1 0 0 1 0 0 1 1 1 1
5 6 0 1 1 0 0 0 0 0 0 0 0 1 1 1 1 0 1 1 0 1 1 0 0 0 0 0 1 0 0 1
5 5 1 0 0 1 1 1 0 0 0 1 0 0 1 0 1 0 1 1 0 0 1 1 1 0 1
1 5 0 1 0 1 1
5 6 0 0 0 0 0 0 1 0 1 1 0 0 1 1 1 1 1 1 1 0 1 1 1 1 1 1 1 1 1 0
1 5 1 0 1 0 0
3 3 0 0 0 0 0 0 1 1 1
1 2 1 1
4 5 0 1 0 1 1 0 1 0 1 0 1 1 1 1 0 1 1 1 0 1 1 1 0
2 5 1 1 0 1 0 0 1 1 0 0
1 5 0 0 0 0 1
4 6 0 0 0 1 1 1 0 1 1 0 1 0 1 0 0 0 0 0 1 0 0 0 1 1
5 6 1 0 0 1 1 1 1 1 1 0 0 0 0 1 1 1 0 1 1 1 0 0 1 0 0 0 0 0 0 1
5 6 1 1 0 0 1 1 0 0 0 0 1 1 0 0 0 1 1 0 1 1 1 0 0 0 1 1 1 1 1 1
5 6 1 1 0 1 0 1 0 1 0 1 0 0 0 0 0 0 1 1 1 0 1 1 0 1 1 1 1 1 0 0
2 6 0 0 1 0 0 0 0 0 0 1 0 1
2 3 0 1 0 1 0 1
3 3 1 0 0 0 0 0 1 0 0
5 6 0 0 0 1 0 0 0 1 0 0 0 1 1 0 1 0 1 0 0 1 1 1 0 0 1 0 0 1 0 1
5 5 0 0 0 1 0 1 1 0 1 1 1 1 1 1 1 0 0 0 1 1 0 1 0 0 0
2 6 1 1 0 0 0 0 0 1 1 0 0 0
3 2 1 0 0 0 0 1
5 3 1 1 1 0 0 0 1 1 1 1 1 1 1 1 1
2 2 0 0 0 0
4 5 0 0 1 1 1 1 0 1 0 1 0 0 0 1 1 1 1 1 1 0 1 1 1 0
4 5 0 0 0 1 0 1 0 0 0 0 0 0 1 1 1 0 1 0 1 0 1 0 0 0
1 5 1 0 0 0 1
2 5 0 0 1 1 1 1 1 0 1 0
1 4 0 0 0 1
1 5 1 0 1 0 1
4 5 0 1 1 0 1 0 1 1 0 0 0 1 1 1 0 1 1 1 0 0 0 1 1
1 3 1 1 1
1 6 1 1 0 1 1 1
5 5 0 0 0 1 1 1 1 1 0 0 1 1 1 1 1 0 0 0 1 1 0 1 1 1 0
5 5 0 0 1 0 0 0 0 1 1 1 0 1 0 1 1 1 0 0 1 0 1 0 0 1 0
4 6 1 1 0 1 1 1 0 0 0 1 0 0 1 0 1 1 1 0 1 1 1 1
1 6 1 1 1 1 1 0
5 6 0 0 0 0 1 0 0 1 0 0 1 0 1 1 0 1 0 0 0 0 1 0 1 0 1 1 1 0 1 0
5 5 0 0 1 1 1 1 0 0 1 0 1 0 1 0 0 0 1 1 0 1 1 0 0 1 1
4 6 0 1 0 1 0 0 0 1 0 0 0 1 0 1 0 0 0 1 0 1 1 1 1 1
2 3 1 0 1 1 1 0
5 5 1 0 1 0 1 1 1 0 0 1 1 1 0 1 1 1 1 0 1 1 0 1 0 1 1
4 6 0 0 1 0 0 0 1 0 1 0 1 1 0 1 0 1 1 0 1 1 1 0
5 6 1 1 0 1 0 0 0 1 1 0 0 0 0 1 1 1 1 1 0 0 1 1 1 1 1 0 1 0
3 3 1 1 1 0 0 0 0 0 0
1 4 1 1 0 1
5 6 0 0 1 1 1 0 0 0 1 0 1 0 1 0 1 0 1 0 0 0 0 1 1 1 0 1 1 1
1 6 0 0 1 1 0 1
2 2 0 1 1 0
2 2 1 1 1 1
1 6 0 0 1 0 1 0
2 4 0 0 1 1 0 0 0 0
2 4 0 0 0 0 0 1 1 1`

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func solveCase(n, m int, cells []int) string {
	totalOn := 0
	for _, v := range cells {
		if v == 1 {
			totalOn++
		}
	}
	if totalOn == 0 {
		return "-1"
	}

	deg := make([]int, n*m)
	adj := make([][4]int, n*m)
	for i := range adj {
		for k := 0; k < 4; k++ {
			adj[i][k] = -1
		}
	}
	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if cells[i*m+j] == 0 {
				continue
			}
			id := i*m + j
			cnt := 0
			for d := 0; d < 4; d++ {
				ni, nj := i+dirs[d][0], j+dirs[d][1]
				if ni >= 0 && ni < n && nj >= 0 && nj < m && cells[ni*m+nj] == 1 {
					adj[id][cnt] = ni*m + nj
					cnt++
				}
			}
			deg[id] = cnt
			if cnt > 2 {
				return "-1"
			}
		}
	}

	vis := make([]bool, n*m)
	start := -1
	for i := 0; i < n*m; i++ {
		if cells[i] == 1 {
			start = i
			break
		}
	}
	if start == -1 {
		return "-1"
	}
	queue := []int{start}
	vis[start] = true
	seen := 1
	for head := 0; head < len(queue); head++ {
		u := queue[head]
		for k := 0; k < deg[u]; k++ {
			v := adj[u][k]
			if !vis[v] {
				vis[v] = true
				queue = append(queue, v)
				seen++
			}
		}
	}
	if seen != totalOn {
		return "-1"
	}

	edgeCount := 0
	endpoints := []int{}
	for i := 0; i < n*m; i++ {
		if cells[i] == 1 {
			edgeCount += deg[i]
			if deg[i] == 1 {
				endpoints = append(endpoints, i)
			}
		}
	}
	edgeCount /= 2
	if edgeCount == 0 {
		return "-1"
	}

	if len(endpoints) == 0 {
		// cycle; keep existing start
	} else if len(endpoints) == 2 {
		start = endpoints[0]
	} else {
		return "-1"
	}

	path := make([]int, 0, edgeCount+1)
	prev := -1
	curr := start
	path = append(path, curr)
	for steps := 0; steps < edgeCount; steps++ {
		found := false
		for k := 0; k < deg[curr]; k++ {
			v := adj[curr][k]
			if v != prev {
				prev, curr = curr, v
				path = append(path, curr)
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	if len(path) < 2 {
		return "-1"
	}

	prevDir := [2]int{path[1]/m - path[0]/m, path[1]%m - path[0]%m}
	runLen := 1
	G := 0
	for i := 1; i < len(path)-1; i++ {
		dx := path[i+1]/m - path[i]/m
		dy := path[i+1]%m - path[i]%m
		if dx == prevDir[0] && dy == prevDir[1] {
			runLen++
		} else {
			G = gcd(G, runLen)
			runLen = 1
			prevDir = [2]int{dx, dy}
		}
	}
	G = gcd(G, runLen)
	if G <= 1 {
		return "-1"
	}
	divs := []int{}
	for k := 2; k*k <= G; k++ {
		if G%k == 0 {
			divs = append(divs, k)
			if k != G/k {
				divs = append(divs, G/k)
			}
		}
	}
	divs = append(divs, G)
	sort.Ints(divs)
	var sb strings.Builder
	for i, v := range divs {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcasesE), "\n")
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			fmt.Fprintf(os.Stderr, "case %d: invalid line\n", idx+1)
			os.Exit(1)
		}
		n, err1 := strconv.Atoi(fields[0])
		m, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			fmt.Fprintf(os.Stderr, "case %d: parse error\n", idx+1)
			os.Exit(1)
		}
		expectVals := n * m
		if len(fields)-2 != expectVals {
			fmt.Fprintf(os.Stderr, "case %d: expected %d values got %d\n", idx+1, expectVals, len(fields)-2)
			os.Exit(1)
		}
		grid := make([]int, expectVals)
		for i := 0; i < expectVals; i++ {
			v, err := strconv.Atoi(fields[2+i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d: invalid grid value\n", idx+1)
				os.Exit(1)
			}
			grid[i] = v
		}
		input := fmt.Sprintf("%d %d\n%s\n", n, m, strings.Join(fields[2:], " "))
		want := solveCase(n, m, grid)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
