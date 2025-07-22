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

type state struct{ wall, pos, time int }

func solve(n, k int, left, right string) string {
	l := []byte(" " + left)
	r := []byte(" " + right)
	visited := make([][]bool, 2)
	visited[0] = make([]bool, n+2)
	visited[1] = make([]bool, n+2)
	q := []state{{0, 1, 0}}
	visited[0][1] = true
	for head := 0; head < len(q); head++ {
		cur := q[head]
		w, p, t := cur.wall, cur.pos, cur.time
		moves := []int{p + 1, p - 1, p + k}
		walls := []int{w, w, 1 - w}
		for i, np := range moves {
			nw := walls[i]
			nt := t + 1
			if np > n {
				return "YES"
			}
			if np <= nt || np < 1 {
				continue
			}
			if visited[nw][np] {
				continue
			}
			if (nw == 0 && l[np] == 'X') || (nw == 1 && r[np] == 'X') {
				continue
			}
			visited[nw][np] = true
			q = append(q, state{nw, np, nt})
		}
	}
	return "NO"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	k := rng.Intn(20) + 1
	left := make([]byte, n)
	right := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(4) == 0 {
			left[i] = 'X'
		} else {
			left[i] = '-'
		}
		if rng.Intn(4) == 0 {
			right[i] = 'X'
		} else {
			right[i] = '-'
		}
	}
	left[0] = '-'
	inp := fmt.Sprintf("%d %d\n%s\n%s\n", n, k, string(left), string(right))
	ans := solve(n, k, string(left), string(right))
	return inp, ans
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
