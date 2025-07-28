package main

import (
	"bytes"
	"container/list"
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
	q := list.New()
	q.PushBack(state{1, 1})
	dist[1][1] = 0
	for q.Len() > 0 {
		e := q.Front()
		q.Remove(e)
		s := e.Value.(state)
		d := dist[s.a][s.b]
		if s.a == n && s.b == m {
			return d
		}
		best := 0
		for i := 1; i <= s.a; i++ {
			for j := 1; j <= s.b; j++ {
				p := i + j
				if good[pair{i, j}] {
					p++
				}
				if p > best {
					best = p
				}
			}
		}
		if s.a < n && best >= s.a+1 && dist[s.a+1][s.b] == -1 {
			dist[s.a+1][s.b] = d + 1
			q.PushBack(state{s.a + 1, s.b})
		}
		if s.b < m && best >= s.b+1 && dist[s.a][s.b+1] == -1 {
			dist[s.a][s.b+1] = d + 1
			q.PushBack(state{s.a, s.b + 1})
		}
	}
	return -1
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
