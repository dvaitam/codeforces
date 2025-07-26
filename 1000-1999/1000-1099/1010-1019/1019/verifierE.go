package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type edge struct {
	to   int
	a, b int64
}

type form struct{ a, b int64 }

func (f form) eval(t int64) int64 { return f.a*t + f.b }

func dfs(u, p int, aSum, bSum int64, adj [][]edge, res []form) {
	res[u] = form{aSum, bSum}
	for _, e := range adj[u] {
		if e.to != p {
			dfs(e.to, u, aSum+e.a, bSum+e.b, adj, res)
		}
	}
}

func compute(n int, m int64, adj [][]edge) []int64 {
	forms := make([][]form, n)
	for i := 0; i < n; i++ {
		forms[i] = make([]form, n)
	}
	for i := 0; i < n; i++ {
		dfs(i, -1, 0, 0, adj, forms[i])
	}
	ans := make([]int64, m)
	for t := int64(0); t < m; t++ {
		best := forms[0][1].eval(t)
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				val := forms[i][j].eval(t)
				if val > best {
					best = val
				}
			}
		}
		ans[t] = best
	}
	return ans
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
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tc := 0; tc < 100; tc++ {
		n := rng.Intn(5) + 2
		m := int64(rng.Intn(5) + 1)
		adj := make([][]edge, n)
		edgesList := make([][4]int64, 0, n-1)
		for i := 0; i < n-1; i++ {
			u := i + 1
			v := rng.Intn(u)
			a := int64(rng.Intn(3))
			b := int64(rng.Intn(5))
			adj[u] = append(adj[u], edge{v, a, b})
			adj[v] = append(adj[v], edge{u, a, b})
			edgesList = append(edgesList, [4]int64{int64(u), int64(v), a, b})
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for _, e := range edgesList {
			fmt.Fprintf(&sb, "%d %d %d %d\n", e[0]+1, e[1]+1, e[2], e[3])
		}
		input := sb.String()
		expected := compute(n, m, adj)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tc+1, err, input)
			os.Exit(1)
		}
		scan := bufio.NewScanner(strings.NewReader(out))
		scan.Split(bufio.ScanWords)
		vals := make([]int64, 0, m)
		for scan.Scan() {
			v, err := strconv.ParseInt(scan.Text(), 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d bad output\n", tc+1)
				os.Exit(1)
			}
			vals = append(vals, v)
		}
		if int64(len(vals)) != m {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d values got %d\n", tc+1, m, len(vals))
			os.Exit(1)
		}
		for i := int64(0); i < m; i++ {
			if vals[i] != expected[i] {
				fmt.Fprintf(os.Stderr, "case %d failed at t=%d: expected %d got %d\ninput:\n%s", tc+1, i, expected[i], vals[i], input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
