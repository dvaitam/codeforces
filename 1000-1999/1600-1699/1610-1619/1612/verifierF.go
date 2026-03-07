package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type pair struct{ a, b int }

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

func bfs(n, m int, good map[pair]bool) int {
	type state struct{ a, b int }
	dist := make([][]int, n+1)
	for i := range dist {
		dist[i] = make([]int, m+1)
		for j := range dist[i] {
			dist[i][j] = -1
		}
	}
	queue := []state{{1, 1}}
	dist[1][1] = 0
	for len(queue) > 0 {
		s := queue[0]
		queue = queue[1:]
		d := dist[s.a][s.b]
		if s.a == n && s.b == m {
			return d
		}
		// Best power: max over all (i,j) with i<=a, j<=b
		best := s.a + s.b
		for i := 1; i <= s.a; i++ {
			for j := 1; j <= s.b; j++ {
				if good[pair{i, j}] && i+j+1 > best {
					best = i + j + 1
				}
			}
		}
		// Jump directly to the maximum reachable armor index
		na := best
		if na > n {
			na = n
		}
		if na > s.a && dist[na][s.b] == -1 {
			dist[na][s.b] = d + 1
			queue = append(queue, state{na, s.b})
		}
		// Jump directly to the maximum reachable weapon index
		nb := best
		if nb > m {
			nb = m
		}
		if nb > s.b && dist[s.a][nb] == -1 {
			dist[s.a][nb] = d + 1
			queue = append(queue, state{s.a, nb})
		}
	}
	return dist[n][m]
}

func parse(output string) (int, error) {
	output = strings.TrimSpace(output)
	var val int
	_, err := fmt.Sscanf(output, "%d", &val)
	return val, err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		n := rand.Intn(5) + 2
		m := rand.Intn(5) + 2
		q := rand.Intn(n*m + 1)
		good := make(map[pair]bool)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n%d\n", n, m, q)
		for i := 0; i < q; i++ {
			a := rand.Intn(n) + 1
			b := rand.Intn(m) + 1
			good[pair{a, b}] = true
			fmt.Fprintf(&sb, "%d %d\n", a, b)
		}
		input := sb.String()
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		got, err := parse(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed to parse output: %v\n", t+1, err)
			os.Exit(1)
		}
		want := bfs(n, m, good)
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", t+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
