package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type state struct {
	pos  int
	mask uint64
	dist int
}

func solveD(n, a, b int, g, k []int) int {
	colorID := make(map[int]int)
	nextID := 0
	for _, x := range g {
		if _, ok := colorID[x]; !ok {
			colorID[x] = nextID
			nextID++
		}
	}
	for _, x := range k {
		if _, ok := colorID[x]; !ok {
			colorID[x] = nextID
			nextID++
		}
	}
	if nextID > 20 {
		return -1
	}
	gID := make([]int, len(g))
	for i, x := range g {
		gID[i] = colorID[x]
	}
	kID := make([]int, len(k))
	for i, x := range k {
		kID[i] = colorID[x]
	}
	startMask := uint64(1) << kID[a-1]
	q := []state{{a - 1, startMask, 0}}
	visited := make(map[[2]int]bool)
	visited[[2]int{a - 1, int(startMask)}] = true
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		mask := cur.mask | (1 << kID[cur.pos])
		if cur.pos == b-1 {
			return cur.dist
		}
		if cur.pos > 0 && (mask&(1<<gID[cur.pos-1])) != 0 {
			np := cur.pos - 1
			nm := mask | (1 << kID[np])
			key := [2]int{np, int(nm)}
			if !visited[key] {
				visited[key] = true
				q = append(q, state{np, nm, cur.dist + 1})
			}
		}
		if cur.pos < n-1 && (mask&(1<<gID[cur.pos])) != 0 {
			np := cur.pos + 1
			nm := mask | (1 << kID[np])
			key := [2]int{np, int(nm)}
			if !visited[key] {
				visited[key] = true
				q = append(q, state{np, nm, cur.dist + 1})
			}
		}
	}
	return -1
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2 // 2..9
	a := rng.Intn(n) + 1
	b := rng.Intn(n-1) + 1
	if b >= a {
		b++
	}
	colors := rng.Intn(6) + 1
	g := make([]int, n-1)
	for i := range g {
		g[i] = rng.Intn(colors) + 1
	}
	k := make([]int, n)
	for i := range k {
		k[i] = rng.Intn(colors) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, a, b))
	for i, v := range g {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range k {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	ans := solveD(n, a, b, g, k)
	return sb.String(), fmt.Sprintf("%d\n", ans)
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expected), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
