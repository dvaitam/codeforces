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

type pipe struct{ from, to, d int }

func solveA(n int, pipes []pipe) string {
	to := make([][2]int, n+1)
	in := make([]bool, n+1)
	out := make([]bool, n+1)
	for _, p := range pipes {
		to[p.from][0] = p.to
		to[p.from][1] = p.d
		out[p.from] = true
		in[p.to] = true
	}
	tanks := []int{}
	for i := 1; i <= n; i++ {
		if out[i] && !in[i] {
			tanks = append(tanks, i)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tanks)))
	for _, t := range tanks {
		at := to[t][0]
		md := to[t][1]
		for out[at] {
			if md > to[at][1] {
				md = to[at][1]
			}
			at = to[at][0]
		}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", t, at, md))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String() + stderr.String(), err
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		usedIn := make([]bool, n+1)
		usedOut := make([]bool, n+1)
		p := rng.Intn(n)
		pipes := make([]pipe, 0, p)
		for j := 0; j < p; j++ {
			a := rng.Intn(n) + 1
			if usedOut[a] {
				continue
			}
			b := rng.Intn(n) + 1
			for b == a || usedIn[b] {
				b = rng.Intn(n) + 1
			}
			d := rng.Intn(1000) + 1
			usedOut[a] = true
			usedIn[b] = true
			pipes = append(pipes, pipe{a, b, d})
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, len(pipes)))
		for _, pp := range pipes {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", pp.from, pp.to, pp.d))
		}
		input := sb.String()
		expect := solveA(n, pipes)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, input, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
