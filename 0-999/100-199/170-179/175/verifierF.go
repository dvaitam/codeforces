package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type Edge struct{ to int }

type State struct {
	bad, dist int
	path      []int
	node      int
}

type PQ []*State

func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool {
	a, b := pq[i], pq[j]
	if a.bad != b.bad {
		return a.bad < b.bad
	}
	if a.dist != b.dist {
		return a.dist < b.dist
	}
	pa, pb := a.path, b.path
	na, nb := len(pa), len(pb)
	lim := na
	if nb < lim {
		lim = nb
	}
	for i := 0; i < lim; i++ {
		if pa[i] != pb[i] {
			return pa[i] < pb[i]
		}
	}
	return na < nb
}
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(*State)) }
func (pq *PQ) Pop() interface{}   { old := *pq; n := len(old); x := old[n-1]; *pq = old[:n-1]; return x }

func solveCase(n, m int, good []int, shortcuts [][]int, queries []string) []string {
	adj := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		u := good[i]
		v := good[(i+1)%m]
		adj[u] = append(adj[u], Edge{v})
		adj[v] = append(adj[v], Edge{u})
	}
	for _, pts := range shortcuts {
		for j := 0; j+1 < len(pts); j++ {
			u, v := pts[j], pts[j+1]
			adj[u] = append(adj[u], Edge{v})
			adj[v] = append(adj[v], Edge{u})
		}
	}
	for i := 1; i <= n; i++ {
		sort.Slice(adj[i], func(a, b int) bool { return adj[i][a].to < adj[i][b].to })
	}
	weights := make(map[int]int)
	res := []string{}
	for _, line := range queries {
		parts := strings.Fields(line)
		op := parts[0]
		s, _ := strconv.Atoi(parts[1])
		t, _ := strconv.Atoi(parts[2])
		if op == "+" {
			key1 := s*(n+1) + t
			key2 := t*(n+1) + s
			weights[key1]++
			weights[key2]++
		} else {
			dist := make([]int, n+1)
			badc := make([]int, n+1)
			vis := make([]bool, n+1)
			paths := make([][]int, n+1)
			for i := 1; i <= n; i++ {
				dist[i] = 1e9
				badc[i] = 1e9
			}
			pq := &PQ{}
			heap.Init(pq)
			badc[s], dist[s] = 0, 0
			paths[s] = []int{s}
			heap.Push(pq, &State{0, 0, []int{s}, s})
			var ans *State
			for pq.Len() > 0 {
				cur := heap.Pop(pq).(*State)
				u := cur.node
				if vis[u] {
					continue
				}
				vis[u] = true
				if u == t {
					ans = cur
					break
				}
				for _, e := range adj[u] {
					v := e.to
					if vis[v] {
						continue
					}
					key := u*(n+1) + v
					w := weights[key]
					nb := cur.bad + w
					nd := cur.dist + 1
					newPath := append(append([]int{}, cur.path...), v)
					better := false
					if nb < badc[v] || (nb == badc[v] && nd < dist[v]) {
						better = true
					}
					if nb == badc[v] && nd == dist[v] {
						pOld := paths[v]
						for x := 0; x < len(newPath) && x < len(pOld); x++ {
							if newPath[x] != pOld[x] {
								if newPath[x] < pOld[x] {
									better = true
								}
								break
							}
						}
						if !better && len(newPath) < len(pOld) {
							better = true
						}
					}
					if better {
						badc[v], dist[v] = nb, nd
						paths[v] = newPath
						heap.Push(pq, &State{nb, nd, newPath, v})
					}
				}
			}
			if ans == nil {
				res = append(res, "-1")
			} else {
				res = append(res, fmt.Sprintf("%d", ans.bad))
				for j := 1; j < len(ans.path); j++ {
					u := ans.path[j-1]
					v := ans.path[j]
					key1 := u*(n+1) + v
					key2 := v*(n+1) + u
					delete(weights, key1)
					delete(weights, key2)
				}
			}
		}
	}
	return res
}

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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesF.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read testcasesF.txt: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "bad file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expectedOut := make([][]string, t)
	for caseNum := 0; caseNum < t; caseNum++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		good := make([]int, m)
		for i := 0; i < m; i++ {
			scan.Scan()
			good[i], _ = strconv.Atoi(scan.Text())
		}
		shortcuts := make([][]int, m)
		for i := 0; i < m; i++ {
			scan.Scan()
			k, _ := strconv.Atoi(scan.Text())
			pts := make([]int, k)
			for j := 0; j < k; j++ {
				scan.Scan()
				pts[j], _ = strconv.Atoi(scan.Text())
			}
			shortcuts[i] = pts
		}
		scan.Scan()
		q, _ := strconv.Atoi(scan.Text())
		queries := make([]string, q)
		for i := 0; i < q; i++ {
			scan.Scan()
			op := scan.Text()
			scan.Scan()
			s := scan.Text()
			scan.Scan()
			ttt := scan.Text()
			queries[i] = fmt.Sprintf("%s %s %s", op, s, ttt)
		}
		expectedOut[caseNum] = solveCase(n, m, good, shortcuts, queries)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(data)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n%s", err, errBuf.String())
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	outScan.Split(bufio.ScanWords)
	for caseNum := 0; caseNum < t; caseNum++ {
		for _, exp := range expectedOut[caseNum] {
			if !outScan.Scan() {
				fmt.Fprintf(os.Stderr, "missing output for case %d\n", caseNum+1)
				os.Exit(1)
			}
			got := outScan.Text()
			if got != exp {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", caseNum+1, exp, got)
				os.Exit(1)
			}
		}
	}
	if outScan.Scan() {
		fmt.Fprintln(os.Stderr, "extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
