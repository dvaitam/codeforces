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

func runBinary(bin, input string) (string, error) {
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
		return out.String(), fmt.Errorf("%v: %s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase() (string, int, []edge) {
	n := rand.Intn(4) + 2
	edges := make([]edge, n-1)
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		edges[i-2] = edge{p, i}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	return sb.String(), n, edges
}

func verify(n int, edges []edge, out string) error {
	adj := make([]map[int]bool, n+1)
	deg := 0
	for i := 1; i <= n; i++ {
		adj[i] = make(map[int]bool)
	}
	for _, e := range edges {
		adj[e.u][e.v] = true
		adj[e.v][e.u] = true
	}
	for i := 1; i <= n; i++ {
		if len(adj[i]) > deg {
			deg = len(adj[i])
		}
	}
	scan := bufio.NewScanner(strings.NewReader(out))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return fmt.Errorf("missing k")
	}
	k, err := strconv.Atoi(scan.Text())
	if err != nil || k <= 0 {
		return fmt.Errorf("bad k")
	}
	type pair struct{ v, t int }
	seq := make([]pair, 0, k)
	for i := 0; i < k; i++ {
		if !scan.Scan() {
			return fmt.Errorf("bad output")
		}
		v, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			return fmt.Errorf("bad output")
		}
		t, _ := strconv.Atoi(scan.Text())
		seq = append(seq, pair{v, t})
	}
	if scan.Scan() {
		return fmt.Errorf("extra output")
	}
	if len(seq) != k {
		return fmt.Errorf("wrong length")
	}
	if seq[0].v != 1 || seq[0].t != 0 {
		return fmt.Errorf("must start at 1 0")
	}
	if seq[len(seq)-1].v != 1 {
		return fmt.Errorf("must end at 1")
	}
	seenPair := make(map[[2]int]bool)
	visited := make(map[int]bool)
	maxT := 0
	for i := 0; i < len(seq); i++ {
		p := seq[i]
		if seenPair[[2]int{p.v, p.t}] {
			return fmt.Errorf("duplicate pair")
		}
		seenPair[[2]int{p.v, p.t}] = true
		visited[p.v] = true
		if p.t > maxT {
			maxT = p.t
		}
		if i > 0 {
			prev := seq[i-1]
			if p.v == prev.v {
				if p.t >= prev.t {
					return fmt.Errorf("time must decrease when staying")
				}
			} else {
				if p.t != prev.t+1 {
					return fmt.Errorf("wrong time increment")
				}
				if !adj[prev.v][p.v] {
					return fmt.Errorf("edge %d-%d doesn't exist", prev.v, p.v)
				}
			}
		}
	}
	if len(visited) != n {
		return fmt.Errorf("not all nodes visited")
	}
	if maxT != deg {
		return fmt.Errorf("max time %d != expected %d", maxT, deg)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input, n, edges := generateCase()
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := verify(n, edges, out); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
