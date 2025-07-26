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

type Edge struct{ u, v int }

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

func verifySolution(n, m int, c []int, edges []Edge, out string) bool {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return false
	}
	tok := sc.Text()
	if tok != "YES" {
		return false
	}
	weights := make([]int, m)
	for i := 0; i < m; i++ {
		if !sc.Scan() {
			return false
		}
		v, err := strconv.Atoi(sc.Text())
		if err != nil {
			return false
		}
		if v < -2*n*n || v > 2*n*n {
			return false
		}
		weights[i] = v
	}
	sums := make([]int, n+1)
	for i, e := range edges {
		w := weights[i]
		sums[e.u] += w
		sums[e.v] += w
	}
	for i := 1; i <= n; i++ {
		if sums[i] != c[i-1] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(4) + 2
		edges := make([]Edge, 0)
		for i := 2; i <= n; i++ {
			p := rng.Intn(i-1) + 1
			edges = append(edges, Edge{p, i})
		}
		if rng.Intn(2) == 0 {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			for u == v {
				v = rng.Intn(n) + 1
			}
			edges = append(edges, Edge{u, v})
		}
		m := len(edges)
		weights := make([]int, m)
		for i := 0; i < m; i++ {
			weights[i] = rng.Intn(21) - 10
		}
		c := make([]int, n)
		for i := 0; i < m; i++ {
			e := edges[i]
			c[e.u-1] += weights[i]
			c[e.v-1] += weights[i]
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", c[i])
		}
		sb.WriteByte('\n')
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
		}
		input := sb.String()
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", t+1, err, input)
			os.Exit(1)
		}
		if !verifySolution(n, m, c, edges, got) {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid solution\ninput:%s\noutput:%s", t+1, input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
