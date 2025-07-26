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

type edge struct{ u, v int }

func existsValid(n int, edges []edge) bool {
	adj := make([][]int, n)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
	}
	m := 1 << n
	for mask := 0; mask < m; mask++ {
		sel := make([]bool, n)
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				sel[i] = true
			}
		}
		valid := true
		for _, e := range edges {
			if sel[e.u] && sel[e.v] {
				valid = false
				break
			}
		}
		if !valid {
			continue
		}
		reach := make([]bool, n)
		for i := 0; i < n; i++ {
			if sel[i] {
				for _, y := range adj[i] {
					reach[y] = true
					for _, z := range adj[y] {
						reach[z] = true
					}
				}
			}
		}
		ok := true
		for v := 0; v < n; v++ {
			if !sel[v] && !reach[v] {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}

func verify(n int, edges []edge, output string) bool {
	scan := bufio.NewScanner(strings.NewReader(output))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return false
	}
	k, err := strconv.Atoi(scan.Text())
	if err != nil || k < 0 || k > n {
		return false
	}
	chosen := make([]int, 0, k)
	used := make(map[int]bool)
	for i := 0; i < k; i++ {
		if !scan.Scan() {
			return false
		}
		v, err := strconv.Atoi(scan.Text())
		if err != nil || v < 1 || v > n || used[v] {
			return false
		}
		used[v] = true
		chosen = append(chosen, v-1)
	}
	if scan.Scan() { /* extra tokens allowed? we ignore? but let's ensure none*/
		return false
	}
	adj := make([][]int, n)
	mat := make([][]bool, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]bool, n)
	}
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		mat[e.u][e.v] = true
	}
	sel := make([]bool, n)
	for _, v := range chosen {
		sel[v] = true
	}
	for i := 0; i < n; i++ {
		for _, to := range adj[i] {
			if sel[i] && sel[to] {
				return false
			}
		}
	}
	reach := make([]bool, n)
	for _, x := range chosen {
		for _, y := range adj[x] {
			reach[y] = true
			for _, z := range adj[y] {
				reach[z] = true
			}
		}
	}
	for v := 0; v < n; v++ {
		if !sel[v] && !reach[v] {
			return false
		}
	}
	return true
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tc := 0; tc < 100; tc++ {
		for {
			n := rng.Intn(4) + 2
			maxEdges := n * (n - 1)
			m := rng.Intn(maxEdges + 1)
			edges := make([]edge, 0, m)
			seen := make(map[[2]int]bool)
			for len(edges) < m {
				u := rng.Intn(n)
				v := rng.Intn(n)
				if u == v {
					continue
				}
				key := [2]int{u, v}
				if seen[key] {
					continue
				}
				seen[key] = true
				edges = append(edges, edge{u, v})
			}
			if existsValid(n, edges) {
				var sb strings.Builder
				fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
				for _, e := range edges {
					fmt.Fprintf(&sb, "%d %d\n", e.u+1, e.v+1)
				}
				input := sb.String()
				out, err := run(bin, input)
				if err != nil {
					fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tc+1, err, input)
					os.Exit(1)
				}
				if !verify(n, edges, out) {
					fmt.Fprintf(os.Stderr, "case %d failed: invalid output\ninput:\n%s\noutput:\n%s", tc+1, input, out)
					os.Exit(1)
				}
				break
			}
		}
	}
	fmt.Println("All tests passed")
}
