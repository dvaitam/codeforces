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

var (
	nGlobal int
	aGlobal int
	bGlobal int
	hGlobal []int
	shoot   []int
)

func dfs(l, balls, last int) bool {
	if balls == 0 {
		for i := 1; i <= nGlobal; i++ {
			if hGlobal[i] >= 0 {
				return false
			}
		}
		return true
	}
	if l <= nGlobal && hGlobal[l] < 0 {
		return dfs(l+1, balls, last)
	}
	lb := last
	if l > 2 && l > lb {
		lb = l
	}
	if lb < 2 {
		lb = 2
	}
	ub := nGlobal - 1
	if l+1 < ub {
		ub = l + 1
	}
	for i := lb; i <= ub; i++ {
		shoot[balls] = i
		hGlobal[i] -= aGlobal
		hGlobal[i-1] -= bGlobal
		hGlobal[i+1] -= bGlobal
		if dfs(l, balls-1, i) {
			return true
		}
		hGlobal[i] += aGlobal
		hGlobal[i-1] += bGlobal
		hGlobal[i+1] += bGlobal
	}
	return false
}

// solve uses BFS to find the optimal sequence of shots
func solve(n, a, b int, h []int) (int, []int) {
	type node struct {
		h   []int
		seq []int
	}
	serialize := func(arr []int) string {
		var sb strings.Builder
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(',')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		return sb.String()
	}
	start := make([]int, n)
	copy(start, h)
	vis := map[string]bool{serialize(start): true}
	q := []node{{h: start}}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		allDead := true
		for _, v := range cur.h {
			if v >= 0 {
				allDead = false
				break
			}
		}
		if allDead {
			return len(cur.seq), cur.seq
		}
		for i := 1; i < n-1; i++ {
			nxt := make([]int, n)
			copy(nxt, cur.h)
			nxt[i] -= a
			nxt[i-1] -= b
			nxt[i+1] -= b
			key := serialize(nxt)
			if !vis[key] {
				vis[key] = true
				seq := append(append([]int(nil), cur.seq...), i+1)
				q = append(q, node{h: nxt, seq: seq})
			}
		}
	}
	return -1, nil
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 3 // 3..5
	a := rng.Intn(5) + 2
	b := rng.Intn(a-1) + 1
	heights := make([]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, a, b))
	for i := 0; i < n; i++ {
		heights[i] = rng.Intn(4) + 1
		sb.WriteString(fmt.Sprintf("%d", heights[i]))
		if i+1 < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	ans, seq := solve(n, a, b, heights)
	var exp strings.Builder
	exp.WriteString(fmt.Sprintf("%d\n", ans))
	for i, v := range seq {
		exp.WriteString(fmt.Sprintf("%d", v))
		if i+1 < len(seq) {
			exp.WriteByte(' ')
		}
	}
	return sb.String(), exp.String()
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
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected\n%s\ngot\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
