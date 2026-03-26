package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// ---------- embedded reference solver (from CF-accepted solution) ----------

type Point struct {
	x, y int
}

func referenceSolve(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var buf bytes.Buffer
	out := bufio.NewWriter(&buf)

	var n int
	fmt.Fscan(in, &n)
	p := make([]Point, n)
	to := make([]int, n)
	ord := make([]int, 0, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &p[i].x, &p[i].y, &to[i])
		to[i]--
		if to[i] != i {
			ord = append(ord, i)
		}
	}
	if len(ord) == 0 {
		return "0"
	}
	mn := ord[0]
	for _, v := range ord[1:] {
		if p[v].x < p[mn].x || (p[v].x == p[mn].x && p[v].y < p[mn].y) {
			mn = v
		}
	}
	idx := 0
	for i, v := range ord {
		if v == mn {
			idx = i
			break
		}
	}
	ord = append(ord[idx:], ord[:idx]...)
	sort.Slice(ord[1:], func(i, j int) bool {
		a := ord[1+i]
		b := ord[1+j]
		dx1, dy1 := p[a].x-p[mn].x, p[a].y-p[mn].y
		dx2, dy2 := p[b].x-p[mn].x, p[b].y-p[mn].y
		return int64(dx1)*int64(dy2)-int64(dy1)*int64(dx2) > 0
	})
	ops := make([][2]int, 0, n)
	makeOp := func(i, j int) {
		to[i], to[j] = to[j], to[i]
		ops = append(ops, [2]int{i, j})
	}
	par := make([]int, n)
	for i := range par {
		par[i] = i
	}
	var find func(int) int
	find = func(v int) int {
		if par[v] != v {
			par[v] = find(par[v])
		}
		return par[v]
	}
	used := make([]bool, n)
	for i := 0; i < n; i++ {
		if used[i] {
			continue
		}
		for j := i; !used[j]; j = to[j] {
			par[find(j)] = find(i)
			used[j] = true
		}
	}
	for i := 1; i < len(ord)-1; i++ {
		u := ord[i]
		v := ord[i+1]
		r1, r2 := find(u), find(v)
		if r1 == r2 {
			continue
		}
		makeOp(u, v)
		par[r1] = r2
	}
	for i := range used {
		used[i] = false
	}
	path := make([]int, 0, n)
	for i := mn; !used[i]; i = to[i] {
		path = append(path, i)
		used[i] = true
	}
	for j := 1; j < len(path); j++ {
		makeOp(mn, path[j])
	}
	fmt.Fprintln(out, len(ops))
	for _, pr := range ops {
		fmt.Fprintln(out, pr[0]+1, pr[1]+1)
	}
	out.Flush()
	return strings.TrimSpace(buf.String())
}

// ---------- test generator ----------

func genCase(r *rand.Rand) string {
	n := r.Intn(8) + 2 // 2..9
	to := make([]int, n)
	diff := 0
	for i := 0; i < n; i++ {
		to[i] = r.Intn(n) + 1
		if to[i] != i+1 {
			diff++
		}
	}
	if diff == 0 {
		idx := r.Intn(n)
		val := r.Intn(n) + 1
		for val == idx+1 {
			val = r.Intn(n) + 1
		}
		to[idx] = val
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		x := r.Intn(100)
		y := r.Intn(100)
		fmt.Fprintf(&sb, "%d %d %d\n", x, y, to[i])
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		expect := referenceSolve(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
