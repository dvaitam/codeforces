package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type graph struct {
	n   int
	adj [][]bool
}

func genGraph(rng *rand.Rand) graph {
	n := rng.Intn(5) + 2 // 2..6
	adj := make([][]bool, n+1)
	for i := range adj {
		adj[i] = make([]bool, n+1)
	}
	// tree edges
	for i := 2; i <= n; i++ {
		j := rng.Intn(i-1) + 1
		adj[i][j] = true
		adj[j][i] = true
	}
	// extra edges
	for k := 0; k < n; k++ {
		if rng.Float64() < 0.3 {
			a := rng.Intn(n) + 1
			b := rng.Intn(n) + 1
			if a != b && !adj[a][b] {
				adj[a][b] = true
				adj[b][a] = true
			}
		}
	}
	return graph{n: n, adj: adj}
}

func countEdges(sub []int, g graph) int {
	mark := make([]bool, g.n+1)
	for _, v := range sub {
		mark[v] = true
	}
	cnt := 0
	for i := 1; i <= g.n; i++ {
		if !mark[i] {
			continue
		}
		for j := i + 1; j <= g.n; j++ {
			if mark[j] && g.adj[i][j] {
				cnt++
			}
		}
	}
	return cnt
}

func bipartiteCheck(g graph) (bool, []int) {
	color := make([]int, g.n+1)
	for i := range color {
		color[i] = -1
	}
	queue := []int{1}
	color[1] = 0
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for u := 1; u <= g.n; u++ {
			if !g.adj[v][u] {
				continue
			}
			if color[u] == -1 {
				color[u] = color[v] ^ 1
				queue = append(queue, u)
			} else if color[u] == color[v] {
				return false, nil
			}
		}
	}
	return true, color
}

func parseList(s string) []int {
	f := strings.Fields(s)
	if len(f) == 0 {
		return nil
	}
	var k int
	fmt.Sscanf(f[0], "%d", &k)
	res := make([]int, 0, k)
	for i := 1; i < len(f); i++ {
		var v int
		fmt.Sscanf(f[i], "%d", &v)
		res = append(res, v)
	}
	return res
}

func runCase(bin string, g graph) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	fmt.Fprintf(stdin, "%d\n", g.n)
	reader := bufio.NewReader(stdout)
	queryCount := 0
	bip, _ := bipartiteCheck(g)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read error: %v stderr:%s", err, stderr.String())
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "?") {
			queryCount++
			if queryCount > 20000 {
				return fmt.Errorf("too many queries")
			}
			fields := strings.Fields(line)
			subset := make([]int, 0, len(fields)-1)
			for _, s := range fields[1:] {
				var v int
				fmt.Sscanf(s, "%d", &v)
				if v < 1 || v > g.n {
					return fmt.Errorf("bad vertex %d", v)
				}
				subset = append(subset, v)
			}
			ans := countEdges(subset, g)
			fmt.Fprintf(stdin, "%d\n", ans)
		} else if strings.HasPrefix(line, "!") {
			rest := strings.TrimSpace(line[1:])
			if bip {
				if !strings.HasPrefix(strings.ToUpper(rest), "Y") {
					return fmt.Errorf("expected Y for bipartite")
				}
				l1, _ := reader.ReadString('\n')
				l2, _ := reader.ReadString('\n')
				g1 := parseList(l1)
				g2 := parseList(l2)
				if len(g1)+len(g2) != g.n {
					return fmt.Errorf("wrong partition size")
				}
				assign := make([]int, g.n+1)
				for _, v := range g1 {
					assign[v] = 0
				}
				for _, v := range g2 {
					if assign[v] != 0 {
						return fmt.Errorf("dup vertex")
					}
					assign[v] = 1
				}
				for i := 1; i <= g.n; i++ {
					for j := i + 1; j <= g.n; j++ {
						if g.adj[i][j] && assign[i] == assign[j] {
							return fmt.Errorf("edge inside partition")
						}
					}
				}
				stdin.Close()
				return cmd.Wait()
			} else {
				if !strings.HasPrefix(strings.ToUpper(rest), "N") {
					return fmt.Errorf("expected N for non-bipartite")
				}
				l1, _ := reader.ReadString('\n')
				cyc := parseList(l1)
				if len(cyc)%2 == 0 || len(cyc) < 3 {
					return fmt.Errorf("invalid cycle length")
				}
				for i := 0; i < len(cyc); i++ {
					v := cyc[i]
					if v < 1 || v > g.n {
						return fmt.Errorf("bad vertex")
					}
					u := cyc[(i+1)%len(cyc)]
					if !g.adj[v][u] {
						return fmt.Errorf("edge missing")
					}
				}
				stdin.Close()
				return cmd.Wait()
			}
		}
	}
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, string(out))
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		g := genGraph(rng)
		if err := runCase(bin, g); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
