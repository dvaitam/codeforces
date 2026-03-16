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

const testcasesCRaw = `100
8 3
2 1
3 1
4 3
5 4
6 2
7 4
8 1
7 3
2 1
3 2
4 2
5 3
6 4
7 1
6 3
2 1
3 2
4 2
5 2
6 3
5 3
2 1
3 1
4 1
5 1
6 2
2 1
3 1
4 3
5 1
6 3
9 3
2 1
3 2
4 1
5 2
6 5
7 4
8 1
9 3
9 2
2 1
3 1
4 2
5 2
6 2
7 4
8 3
9 7
6 2
2 1
3 1
4 1
5 2
6 5
4 2
2 1
3 2
4 3
5 2
2 1
3 1
4 3
5 4
4 3
2 1
3 1
4 2
6 2
2 1
3 1
4 2
5 3
6 1
4 2
2 1
3 2
4 3
4 1
2 1
3 1
4 3
7 1
2 1
3 1
4 1
5 1
6 4
7 4
5 1
2 1
3 1
4 1
5 1
6 3
2 1
3 2
4 1
5 1
6 3
5 1
2 1
3 2
4 2
5 1
7 3
2 1
3 2
4 1
5 1
6 3
7 2
9 3
2 1
3 1
4 1
5 4
6 4
7 6
8 4
9 5
9 2
2 1
3 1
4 1
5 4
6 4
7 6
8 3
9 2
9 3
2 1
3 1
4 2
5 4
6 3
7 2
8 7
9 4
7 1
2 1
3 1
4 1
5 2
6 5
7 5
9 3
2 1
3 1
4 3
5 3
6 1
7 6
8 5
9 6
8 2
2 1
3 1
4 3
5 4
6 4
7 4
8 5
7 3
2 1
3 1
4 3
5 1
6 4
7 3
8 1
2 1
3 2
4 3
5 2
6 3
7 3
8 5
9 1
2 1
3 2
4 3
5 1
6 3
7 6
8 2
9 7
6 3
2 1
3 1
4 2
5 1
6 2
4 1
2 1
3 2
4 3
7 2
2 1
3 2
4 3
5 3
6 4
7 3
4 2
2 1
3 1
4 3
7 2
2 1
3 2
4 2
5 4
6 2
7 6
7 3
2 1
3 1
4 2
5 1
6 4
7 3
7 1
2 1
3 1
4 3
5 4
6 1
7 1
7 2
2 1
3 2
4 2
5 2
6 1
7 6
5 1
2 1
3 1
4 3
5 2
9 1
2 1
3 2
4 3
5 3
6 3
7 4
8 6
9 6
8 2
2 1
3 1
4 3
5 4
6 2
7 1
8 7
7 2
2 1
3 1
4 2
5 2
6 5
7 2
7 1
2 1
3 1
4 1
5 4
6 2
7 1
4 2
2 1
3 1
4 3
7 2
2 1
3 2
4 3
5 4
6 4
7 3
9 2
2 1
3 1
4 2
5 4
6 5
7 4
8 6
9 2
7 1
2 1
3 2
4 2
5 2
6 1
7 6
7 3
2 1
3 1
4 3
5 3
6 1
7 1
8 1
2 1
3 2
4 1
5 3
6 4
7 6
8 2
5 1
2 1
3 2
4 3
5 4
9 2
2 1
3 1
4 1
5 3
6 2
7 1
8 5
9 8
9 2
2 1
3 2
4 1
5 4
6 5
7 5
8 5
9 7
9 2
2 1
3 2
4 1
5 1
6 2
7 1
8 5
9 7
5 2
2 1
3 2
4 1
5 4
9 3
2 1
3 2
4 1
5 4
6 4
7 5
8 3
9 4
8 3
2 1
3 1
4 2
5 1
6 2
7 4
8 5
6 2
2 1
3 2
4 1
5 1
6 2
8 3
2 1
3 1
4 2
5 1
6 1
7 5
8 4
8 1
2 1
3 1
4 3
5 3
6 5
7 3
8 3
6 2
2 1
3 1
4 2
5 4
6 2
9 2
2 1
3 2
4 3
5 2
6 3
7 6
8 2
9 6
4 3
2 1
3 2
4 3
7 3
2 1
3 1
4 3
5 2
6 4
7 6
7 1
2 1
3 2
4 1
5 3
6 2
7 5
8 1
2 1
3 1
4 1
5 3
6 1
7 4
8 2
6 2
2 1
3 1
4 1
5 3
6 1
5 3
2 1
3 2
4 3
5 4
9 3
2 1
3 1
4 2
5 3
6 5
7 4
8 4
9 1
4 1
2 1
3 1
4 1
9 3
2 1
3 1
4 2
5 2
6 2
7 5
8 7
9 6
5 3
2 1
3 2
4 1
5 1
6 1
2 1
3 2
4 3
5 2
6 4
6 3
2 1
3 2
4 3
5 4
6 1
6 3
2 1
3 2
4 3
5 3
6 5
9 1
2 1
3 1
4 2
5 2
6 1
7 5
8 3
9 3
5 2
2 1
3 2
4 2
5 2
6 3
2 1
3 2
4 2
5 3
6 3
6 1
2 1
3 2
4 2
5 2
6 5
7 3
2 1
3 1
4 2
5 4
6 3
7 1
7 3
2 1
3 2
4 2
5 1
6 3
7 5
5 2
2 1
3 2
4 1
5 2
6 2
2 1
3 2
4 2
5 2
6 3
9 2
2 1
3 2
4 3
5 2
6 5
7 4
8 1
9 2
8 3
2 1
3 1
4 1
5 3
6 3
7 6
8 3
9 3
2 1
3 1
4 1
5 3
6 1
7 6
8 2
9 5
8 2
2 1
3 1
4 3
5 3
6 5
7 4
8 4
5 3
2 1
3 2
4 2
5 4
9 1
2 1
3 1
4 3
5 1
6 4
7 2
8 6
9 4
5 1
2 1
3 1
4 2
5 1
5 1
2 1
3 1
4 3
5 4
6 1
2 1
3 1
4 2
5 2
6 2
4 1
2 1
3 1
4 2
5 3
2 1
3 1
4 1
5 1
7 1
2 1
3 1
4 3
5 3
6 3
7 1
6 1
2 1
3 2
4 3
5 3
6 2
5 3
2 1
3 1
4 2
5 3
4 1
2 1
3 2
4 2
4 1
2 1
3 1
4 2
7 1
2 1
3 2
4 3
5 1
6 3
7 6
6 3
2 1
3 2
4 1
5 4
6 2
8 2
2 1
3 1
4 3
5 3
6 3
7 6
8 5
5 3
2 1
3 1
4 2
5 3
`

var (
	n, k       int
	adj        [][]int
	removed    []bool
	size       []int
	totalFloor int64
	totalDiv0  int64
)

func dfsSize(u, p int) int {
	size[u] = 1
	for _, v := range adj[u] {
		if v != p && !removed[v] {
			size[u] += dfsSize(v, u)
		}
	}
	return size[u]
}

func findCentroid(u, p, tot int) int {
	for _, v := range adj[u] {
		if v != p && !removed[v] && size[v] > tot/2 {
			return findCentroid(v, u, tot)
		}
	}
	return u
}

func collect(u, p, dist int, cnt, sum []int64) {
	r := dist % k
	cnt[r]++
	sum[r] += int64(dist / k)
	for _, v := range adj[u] {
		if v != p && !removed[v] {
			collect(v, u, dist+1, cnt, sum)
		}
	}
}

func decompose(u int) {
	tot := dfsSize(u, -1)
	c := findCentroid(u, -1, tot)
	removed[c] = true
	globalCnt := make([]int64, k)
	globalSum := make([]int64, k)
	globalCnt[0] = 1
	for _, v := range adj[c] {
		if removed[v] {
			continue
		}
		cnt := make([]int64, k)
		sum := make([]int64, k)
		collect(v, c, 1, cnt, sum)
		for r := 0; r < k; r++ {
			if cnt[r] == 0 {
				continue
			}
			for b := 0; b < k; b++ {
				if globalCnt[b] == 0 {
					continue
				}
				carry := int64(0)
				if r+b >= k {
					carry = 1
				}
				pairCnt := cnt[r] * globalCnt[b]
				totalFloor += cnt[r]*globalSum[b] + globalCnt[b]*sum[r] + carry*pairCnt
			}
			totalDiv0 += cnt[r] * globalCnt[(k-r)%k]
		}
		for i := 0; i < k; i++ {
			globalCnt[i] += cnt[i]
			globalSum[i] += sum[i]
		}
	}
	for _, v := range adj[c] {
		if !removed[v] {
			decompose(v)
		}
	}
}

func solveCase(nVal, kVal int, edges [][2]int) int64 {
	n = nVal
	k = kVal
	adj = make([][]int, n+1)
	for _, e := range edges {
		a, b := e[0], e[1]
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	removed = make([]bool, n+1)
	size = make([]int, n+1)
	totalFloor = 0
	totalDiv0 = 0
	decompose(1)
	totalPairs := int64(n) * int64(n-1) / 2
	return totalPairs + totalFloor - totalDiv0
}

func runCandidate(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("execution failed: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	scan := bufio.NewScanner(strings.NewReader(testcasesCRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		scan.Scan()
		nVal, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		kVal, _ := strconv.Atoi(scan.Text())
		edges := make([][2]int, nVal-1)
		for i := 0; i < nVal-1; i++ {
			scan.Scan()
			a, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			b, _ := strconv.Atoi(scan.Text())
			edges[i] = [2]int{a, b}
		}
		expected := solveCase(nVal, kVal, edges)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", nVal, kVal))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		out, err := runCandidate(os.Args[1], []byte(sb.String()))
		if err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil || got != expected {
			fmt.Printf("case %d failed: expected %d got %s\n", caseIdx+1, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
