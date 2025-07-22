package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type pair struct{ a, b int }

func expectedOutput(n, t1, t2, k int, arr []pair) string {
	type res struct {
		idx int
		h   int
	}
	out := make([]res, 0, n)
	for i, p := range arr {
		h1 := p.a*t1*(100-k) + p.b*t2*100
		h2 := p.b*t1*(100-k) + p.a*t2*100
		h := h1
		if h2 > h1 {
			h = h2
		}
		out = append(out, res{idx: i + 1, h: h})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].h != out[j].h {
			return out[i].h > out[j].h
		}
		return out[i].idx < out[j].idx
	})
	var sb strings.Builder
	for _, r := range out {
		fmt.Fprintf(&sb, "%d %d.%02d\n", r.idx, r.h/100, r.h%100)
	}
	return sb.String()
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	t1 := rng.Intn(10) + 1
	t2 := rng.Intn(10) + 1
	k := rng.Intn(99) + 1
	arr := make([]pair, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", n, t1, t2, k)
	for i := 0; i < n; i++ {
		a := rng.Intn(10) + 1
		b := rng.Intn(10) + 1
		arr[i] = pair{a, b}
		fmt.Fprintf(&sb, "%d %d\n", a, b)
	}
	input := sb.String()
	exp := expectedOutput(n, t1, t2, k, arr)
	return input, exp
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
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected \n%s\ngot \n%s", expected, got)
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
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
