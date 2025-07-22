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

type edge struct{ u, v int }

func bfs(n int, edges []edge, black []bool, a, b int) int {
	q := []int{a}
	dist := make([]int, n+1)
	for i := range dist {
		dist[i] = -1
	}
	dist[a] = 0
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		if v == b {
			return dist[v]
		}
		for idx, e := range edges {
			if !black[idx] {
				continue
			}
			if e.u == v && dist[e.v] == -1 {
				dist[e.v] = dist[v] + 1
				q = append(q, e.v)
			} else if e.v == v && dist[e.u] == -1 {
				dist[e.u] = dist[v] + 1
				q = append(q, e.u)
			}
		}
	}
	return -1
}

func genTree(rng *rand.Rand, n int) []edge {
	edges := make([]edge, 0, n-1)
	if rng.Intn(2) == 0 {
		for i := 2; i <= n; i++ {
			edges = append(edges, edge{i - 1, i})
		}
	} else {
		for i := 2; i <= n; i++ {
			edges = append(edges, edge{1, i})
		}
	}
	return edges
}

func genCase(rng *rand.Rand) (string, int, []edge, []string) {
	n := rng.Intn(6) + 2
	edges := genTree(rng, n)
	m := rng.Intn(20) + 1
	var ops []string
	var b strings.Builder
	fmt.Fprintln(&b, n)
	for _, e := range edges {
		fmt.Fprintf(&b, "%d %d\n", e.u, e.v)
	}
	fmt.Fprintln(&b, m)
	for i := 0; i < m; i++ {
		t := rng.Intn(3) + 1
		if t == 1 {
			id := rng.Intn(len(edges)) + 1
			ops = append(ops, fmt.Sprintf("1 %d", id))
			fmt.Fprintf(&b, "1 %d\n", id)
		} else if t == 2 {
			id := rng.Intn(len(edges)) + 1
			ops = append(ops, fmt.Sprintf("2 %d", id))
			fmt.Fprintf(&b, "2 %d\n", id)
		} else {
			a := rng.Intn(n) + 1
			c := rng.Intn(n) + 1
			ops = append(ops, fmt.Sprintf("3 %d %d", a, c))
			fmt.Fprintf(&b, "3 %d %d\n", a, c)
		}
	}
	return b.String(), n, edges, ops
}

func simulate(n int, edges []edge, ops []string) []int {
	black := make([]bool, len(edges))
	for i := range black {
		black[i] = true
	}
	var ans []int
	for _, line := range ops {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		switch fields[0] {
		case "1":
			var id int
			fmt.Sscan(fields[1], &id)
			black[id-1] = true
		case "2":
			var id int
			fmt.Sscan(fields[1], &id)
			black[id-1] = false
		case "3":
			var a, b int
			fmt.Sscan(fields[1], &a)
			fmt.Sscan(fields[2], &b)
			ans = append(ans, bfs(n, edges, black, a, b))
		}
	}
	return ans
}

func runCase(bin, stringInput string, exp []int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(stringInput)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	var got []int
	for scanner.Scan() {
		var x int
		fmt.Sscan(scanner.Text(), &x)
		got = append(got, x)
	}
	if len(got) != len(exp) {
		return fmt.Errorf("expected %v got %v", exp, got)
	}
	for i := range got {
		if got[i] != exp[i] {
			return fmt.Errorf("expected %v got %v", exp, got)
		}
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
	for i := 0; i < 100; i++ {
		input, n, edges, ops := genCase(rng)
		exp := simulate(n, edges, ops)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
