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

const INF int64 = 1 << 62

func solveB(n, m, p int, d []int64, h []int, t []int64) string {
	pos := make([]int64, n+1)
	for i := 2; i <= n; i++ {
		pos[i] = pos[i-1] + d[i-2]
	}
	ski := make([]int64, m+1)
	for i := 1; i <= m; i++ {
		ski[i] = t[i-1] - pos[h[i-1]]
	}
	sort.Slice(ski[1:], func(i, j int) bool { return ski[i+1] < ski[j+1] })
	S := make([]int64, m+1)
	for i := 1; i <= m; i++ {
		S[i] = S[i-1] + ski[i]
	}
	dpPrev := make([]int64, m+1)
	for i := 1; i <= m; i++ {
		dpPrev[i] = INF
	}
	dpPrev[0] = 0
	cost := func(k, i int) int64 {
		return int64(i-k)*ski[i] - (S[i] - S[k])
	}
	for step := 1; step <= p; step++ {
		dpCur := make([]int64, m+1)
		dpCur[0] = 0
		for i := 1; i <= m; i++ {
			best := INF
			for k := 0; k < i; k++ {
				v := dpPrev[k] + cost(k, i)
				if v < best {
					best = v
				}
			}
			dpCur[i] = best
		}
		dpPrev = dpCur
	}
	return fmt.Sprintf("%d", dpPrev[m])
}

func generateCaseB(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 2
	m := rng.Intn(5) + 1
	if m < 1 {
		m = 1
	}
	p := rng.Intn(m) + 1
	d := make([]int64, n-1)
	for i := range d {
		d[i] = int64(rng.Intn(5) + 1)
	}
	h := make([]int, m)
	tarr := make([]int64, m)
	for i := 0; i < m; i++ {
		h[i] = rng.Intn(n) + 1
		tarr[i] = int64(rng.Intn(30))
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, p)
	for i := 0; i < n-1; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", d[i])
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		fmt.Fprintf(&sb, "%d %d\n", h[i], tarr[i])
	}
	expected := solveB(n, m, p, d, h, tarr)
	return sb.String(), expected
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	if outStr != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseB(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
