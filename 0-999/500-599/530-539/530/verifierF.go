package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runCandidate(bin string, input string) (string, error) {
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

func expected(xs, ys []int) int {
	n := len(xs) - 1
	reachable := make([]bool, n+1)
	for i := 1; i <= 10; i++ {
		visited := make([]bool, n+1)
		q := []int{0}
		visited[0] = true
		for head := 0; head < len(q); head++ {
			u := q[head]
			for v := 1; v <= n; v++ {
				if visited[v] {
					continue
				}
				if xs[u] == xs[v] && abs(ys[u]-ys[v]) == i || ys[u] == ys[v] && abs(xs[u]-xs[v]) == i {
					visited[v] = true
					q = append(q, v)
				}
			}
		}
		for v := 1; v <= n; v++ {
			if visited[v] {
				reachable[v] = true
			}
		}
	}
	maxD := 0
	for i := 1; i <= n; i++ {
		if reachable[i] {
			d := abs(xs[i]) + abs(ys[i])
			if d > maxD {
				maxD = d
			}
		}
	}
	return maxD
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		xs := make([]int, n+1)
		ys := make([]int, n+1)
		used := make(map[[2]int]bool)
		for j := 1; j <= n; j++ {
			for {
				x := rng.Intn(11) - 5
				y := rng.Intn(11) - 5
				if x == 0 && y == 0 {
					continue
				}
				if !used[[2]int{x, y}] {
					used[[2]int{x, y}] = true
					xs[j] = x
					ys[j] = y
					break
				}
			}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 1; j <= n; j++ {
			fmt.Fprintf(&sb, "%d %d\n", xs[j], ys[j])
		}
		input := sb.String()
		xs[0], ys[0] = 0, 0
		want := expected(xs, ys)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil || got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", i+1, want, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
