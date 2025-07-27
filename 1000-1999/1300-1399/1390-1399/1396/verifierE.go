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

func bfs(n int, g [][]int, start int) []int {
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	q := []int{start}
	dist[start] = 0
	for head := 0; head < len(q); head++ {
		v := q[head]
		for _, to := range g[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}
	}
	return dist
}

func allDistances(n int, edges [][2]int) [][]int {
	g := make([][]int, n)
	for _, e := range edges {
		u, v := e[0]-1, e[1]-1
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = bfs(n, g, i)
	}
	return dist
}

func possible(dist [][]int, k int) bool {
	n := len(dist)
	dp := make([]map[int]struct{}, 1<<n)
	dp[0] = map[int]struct{}{0: {}}
	for mask := 0; mask < (1 << n); mask++ {
		if dp[mask] == nil {
			continue
		}
		var first int
		for first = 0; first < n; first++ {
			if mask&(1<<first) == 0 {
				break
			}
		}
		if first >= n {
			continue
		}
		for second := first + 1; second < n; second++ {
			if mask&(1<<second) != 0 {
				continue
			}
			newMask := mask | 1<<first | 1<<second
			if dp[newMask] == nil {
				dp[newMask] = make(map[int]struct{})
			}
			for sum := range dp[mask] {
				dp[newMask][sum+dist[first][second]] = struct{}{}
			}
		}
	}
	_, ok := dp[(1<<n)-1][k]
	return ok
}

func verifyOutput(n, k int, edges [][2]int, out string) error {
	dist := allDistances(n, edges)
	reader := bufio.NewReader(strings.NewReader(strings.TrimSpace(out)))
	var firstTok string
	if _, err := fmt.Fscan(reader, &firstTok); err != nil {
		return fmt.Errorf("failed to read first token: %v", err)
	}
	if firstTok == "NO" {
		if possible(dist, k) {
			return fmt.Errorf("matching exists but NO given")
		}
		if _, err := fmt.Fscan(reader); err == nil {
			return fmt.Errorf("extra output")
		}
		return nil
	}
	if firstTok != "YES" {
		return fmt.Errorf("expected YES or NO")
	}
	m := n / 2
	used := make([]bool, n)
	sum := 0
	for i := 0; i < m; i++ {
		var u, v int
		if _, err := fmt.Fscan(reader, &u, &v); err != nil {
			return fmt.Errorf("failed to read pair %d: %v", i+1, err)
		}
		if u < 1 || u > n || v < 1 || v > n || u == v {
			return fmt.Errorf("invalid pair %d", i+1)
		}
		if used[u-1] || used[v-1] {
			return fmt.Errorf("node reused in pair %d", i+1)
		}
		used[u-1] = true
		used[v-1] = true
		sum += dist[u-1][v-1]
	}
	if _, err := fmt.Fscan(reader); err == nil {
		return fmt.Errorf("extra output")
	}
	for _, u := range used {
		if !u {
			return fmt.Errorf("not all nodes matched")
		}
	}
	if sum != k {
		return fmt.Errorf("sum %d does not match k %d", sum, k)
	}
	return nil
}

func genTest(rng *rand.Rand) string {
	n := rng.Intn(4)*2 + 2 // 2,4,6,8,10
	if n > 8 {
		n = 8
	}
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{i, p}
	}
	kVal := rng.Intn(n*n) + 1
	var buf strings.Builder
	fmt.Fprintf(&buf, "%d %d\n", n, kVal)
	for _, e := range edges {
		fmt.Fprintf(&buf, "%d %d\n", e[0], e[1])
	}
	return buf.String()
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genTest(rng)
		reader := bufio.NewReader(strings.NewReader(input))
		var n, kVal int
		fmt.Fscan(reader, &n, &kVal)
		edges := make([][2]int, n-1)
		for j := 0; j < n-1; j++ {
			fmt.Fscan(reader, &edges[j][0], &edges[j][1])
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := verifyOutput(n, kVal, edges, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
