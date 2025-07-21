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

type query struct{ l, r int }

func expected(a []int, qs []query) []int64 {
	res := make([]int64, len(qs))
	for i, q := range qs {
		cnt := make(map[int]int)
		var ans int64
		for j := q.l; j <= q.r; j++ {
			v := a[j]
			cnt[v]++
		}
		for v, c := range cnt {
			ans += int64(c*c) * int64(v)
		}
		res[i] = ans
	}
	return res
}

func generateCase(rng *rand.Rand) (string, []int64) {
	n := rng.Intn(30) + 1
	t := rng.Intn(20) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(20) + 1
	}
	qs := make([]query, t)
	for i := 0; i < t; i++ {
		l := rng.Intn(n)
		r := rng.Intn(n-l) + l
		qs[i] = query{l: r}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, t))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for _, q := range qs {
		sb.WriteString(fmt.Sprintf("%d %d\n", q.l+1, q.r+1))
	}
	return sb.String(), expected(a, qs)
}

func runCase(bin, input string, exp []int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Fields(out.String())
	if len(lines) != len(exp) {
		return fmt.Errorf("expected %d lines got %d", len(exp), len(lines))
	}
	for i, txt := range lines {
		var val int64
		if _, err := fmt.Sscan(txt, &val); err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		if val != exp[i] {
			return fmt.Errorf("ans %d expected %d got %d", i+1, exp[i], val)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
