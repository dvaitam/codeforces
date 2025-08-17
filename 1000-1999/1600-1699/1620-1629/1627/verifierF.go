package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type node struct{ x, y, w int }

// solveF implements the reference solution for problem 1627F used by the
// verifier.  It mirrors the approach of the official solution: we model the
// grid as nodes on the intersections of grid lines and run Dijkstra from the
// centre to the boundary.  Each step costs the number of given pairs that would
// be cut by moving across the corresponding edge and its symmetric counterpart.
// The answer is the total number of pairs minus this minimal cut.
func solveF(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(in, &t)
	var out strings.Builder
	type edgeKey [4]int
	dirs := [][2]int{{0, 1}, {0, -1}, {-1, 0}, {1, 0}}
	for ; t > 0; t-- {
		var m, k int
		fmt.Fscan(in, &m, &k)

		// map to count how many input pairs cross an edge between
		// two grid intersections
		mp := make(map[edgeKey]int)

		for i := 0; i < m; i++ {
			var a, b, c, d int
			fmt.Fscan(in, &a, &b, &c, &d)
			if a > c {
				a, c = c, a
			}
			if b > d {
				b, d = d, b
			}
			if a == c {
				mp[edgeKey{a - 1, b, a, b}]++
				mp[edgeKey{a, b, a - 1, b}]++
			} else {
				mp[edgeKey{a, b - 1, a, b}]++
				mp[edgeKey{a, b, a, b - 1}]++
			}
		}

		// visited matrix for intersections (0..k)
		vis := make([][]bool, k+1)
		for i := range vis {
			vis[i] = make([]bool, k+1)
		}

		// priority queue for Dijkstra's algorithm
		pq := &nodePQ{}
		heap.Init(pq)
		heap.Push(pq, node{k / 2, k / 2, 0})

		for pq.Len() > 0 {
			cur := heap.Pop(pq).(node)
			x, y, w := cur.x, cur.y, cur.w
			if vis[x][y] {
				continue
			}
			vis[x][y] = true
			vis[k-x][k-y] = true
			if x == 0 || x == k || y == 0 || y == k {
				out.WriteString(fmt.Sprintf("%d\n", m-w))
				break
			}
			for _, d := range dirs {
				nx, ny := x+d[0], y+d[1]
				if nx < 0 || nx > k || ny < 0 || ny > k {
					continue
				}
				if vis[nx][ny] {
					continue
				}
				cost := mp[edgeKey{x, y, nx, ny}] + mp[edgeKey{k - x, k - y, k - nx, k - ny}]
				heap.Push(pq, node{nx, ny, w + cost})
			}
		}
	}
	return strings.TrimSpace(out.String())
}

// nodePQ implements heap.Interface for nodes ordered by weight.
type nodePQ []node

func (pq nodePQ) Len() int            { return len(pq) }
func (pq nodePQ) Less(i, j int) bool  { return pq[i].w < pq[j].w }
func (pq nodePQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *nodePQ) Push(x interface{}) { *pq = append(*pq, x.(node)) }
func (pq *nodePQ) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func runProg(bin, input string) (string, error) {
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

func generateTests() []string {
	rng := rand.New(rand.NewSource(6))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(3) + 1
		k := 2 * (rng.Intn(3) + 1)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for j := 0; j < n; j++ {
			r1 := rng.Intn(k) + 1
			c1 := rng.Intn(k) + 1
			dir := rng.Intn(4)
			r2, c2 := r1, c1
			switch dir {
			case 0:
				if r1 > 1 {
					r2 = r1 - 1
				} else {
					r2 = r1 + 1
				}
			case 1:
				if r1 < k {
					r2 = r1 + 1
				} else {
					r2 = r1 - 1
				}
			case 2:
				if c1 > 1 {
					c2 = c1 - 1
				} else {
					c2 = c1 + 1
				}
			case 3:
				if c1 < k {
					c2 = c1 + 1
				} else {
					c2 = c1 - 1
				}
			}
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", r1, c1, r2, c2))
		}
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	tests := generateTests()
	for i, t := range tests {
		expect := solveF(t)
		got, err := runProg(bin, t)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, t, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
