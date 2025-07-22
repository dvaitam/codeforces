package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// solveD replicates the reference solution and returns whether a solution
// exists and one possible assignment if it does.
func solveD(n int, edges [][2]int) (bool, []int) {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	col := make([]int, n+1)
	for i := 1; i <= n; i++ {
		col[i] = -1
	}
	queue := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if col[i] != -1 {
			continue
		}
		col[i] = 0
		queue = queue[:0]
		queue = append(queue, i)
		for qi := 0; qi < len(queue); qi++ {
			u := queue[qi]
			for _, v := range adj[u] {
				if col[v] == -1 {
					col[v] = 1 - col[u]
					queue = append(queue, v)
				}
			}
		}
	}
	cnt := [2]int{}
	for i := 1; i <= n; i++ {
		cnt[col[i]]++
	}
	du := make([]int, n+1)
	for i := 1; i <= n; i++ {
		du[i] = len(adj[i])
	}
	idx := make([]int, n+1)
	tot := 0

	paint := func(x int) {
		tot++
		idx[x] = tot
		vist := make([]bool, n+1)
		for _, v := range adj[x] {
			vist[v] = true
		}
		t := 0
		for i := 1; i <= n && t < 2; i++ {
			if idx[i] == 0 && col[i] != col[x] && !vist[i] {
				tot++
				idx[i] = tot
				t++
			}
		}
	}

	outputIdx := func() []int {
		res := make([]int, n+1)
		for i := 1; i <= n; i++ {
			if idx[i] == 0 && col[i] == 0 {
				tot++
				idx[i] = tot
			}
		}
		for i := 1; i <= n; i++ {
			if idx[i] == 0 && col[i] == 1 {
				tot++
				idx[i] = tot
			}
		}
		for i := 1; i <= n; i++ {
			res[i] = (idx[i]-1)/3 + 1
		}
		return res[1:]
	}

	if cnt[0]%3 == 0 {
		return true, outputIdx()
	}
	if cnt[0]%3 == 2 {
		for i := 1; i <= n; i++ {
			col[i] = 1 - col[i]
		}
		cnt[0], cnt[1] = cnt[1], cnt[0]
	}
	for u := 1; u <= n; u++ {
		if col[u] != 0 {
			continue
		}
		if du[u]+2 <= cnt[1] {
			paint(u)
			return true, outputIdx()
		}
	}
	tcnt := 0
	for u := 1; u <= n; u++ {
		if col[u] != 1 {
			continue
		}
		if du[u]+2 <= cnt[0] {
			paint(u)
			tcnt++
			if tcnt == 2 {
				return true, outputIdx()
			}
		}
	}
	return false, nil
}

func verifyOutput(n int, edges [][2]int, solPossible bool, out string) error {
	rdr := bufio.NewReader(strings.NewReader(strings.TrimSpace(out)))
	first, err := rdr.ReadString('\n')
	if err != nil && err.Error() != "EOF" {
		return fmt.Errorf("read output: %v", err)
	}
	first = strings.TrimSpace(first)
	if first == "NO" {
		if solPossible {
			return fmt.Errorf("expected YES, got NO")
		}
		return nil
	}
	if first != "YES" {
		return fmt.Errorf("first token must be YES or NO")
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(rdr, &arr[i]); err != nil {
			return fmt.Errorf("read assignment: %v", err)
		}
	}
	k := n / 3
	cnt := make([]int, k+1)
	for _, g := range arr {
		if g < 1 || g > k {
			return fmt.Errorf("invalid group id %d", g)
		}
		cnt[g]++
	}
	for g := 1; g <= k; g++ {
		if cnt[g] != 3 {
			return fmt.Errorf("group %d size %d", g, cnt[g])
		}
	}
	for _, e := range edges {
		if arr[e[0]-1] == arr[e[1]-1] {
			return fmt.Errorf("edge %d-%d same group", e[0], e[1])
		}
	}
	if !solPossible {
		return fmt.Errorf("expected NO, but got assignment")
	}
	return nil
}

func runCase(bin string, n int, edges [][2]int) error {
	var input bytes.Buffer
	fmt.Fprintf(&input, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&input, "%d %d\n", e[0], e[1])
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	possible, _ := solveD(n, edges)
	if err := verifyOutput(n, edges, possible, out.String()); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	const tests = 100
	for t := 0; t < tests; t++ {
		k := rng.Intn(5) + 1
		n := 3 * k
		maxEdges := n * (n - 1) / 2
		m := rng.Intn(min(maxEdges, 4*n))
		edges := make([][2]int, 0, m)
		seen := make(map[[2]int]struct{})
		for len(edges) < m {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			if u > v {
				u, v = v, u
			}
			e := [2]int{u, v}
			if _, ok := seen[e]; ok {
				continue
			}
			seen[e] = struct{}{}
			edges = append(edges, e)
		}
		if err := runCase(bin, n, edges); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
