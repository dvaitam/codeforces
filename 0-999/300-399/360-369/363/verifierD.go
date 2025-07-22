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

func expected(n, m int, a int64, b, p []int64) (int, int64) {
	sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })
	sort.Slice(p, func(i, j int) bool { return p[i] < p[j] })
	maxk := n
	if m < n {
		maxk = m
	}
	B := make([]int64, maxk+1)
	for i := 1; i <= maxk; i++ {
		B[i] = B[i-1] + b[i-1]
	}
	P := make([]int64, maxk+1)
	for i := 1; i <= maxk; i++ {
		P[i] = P[i-1] + p[i-1]
	}
	best := 0
	var spend int64 = 0
	for r := 1; r <= maxk; r++ {
		need := P[r] - a
		if need < 0 {
			need = 0
		}
		if B[r] >= need {
			best = r
			spend = need
		}
	}
	return best, spend
}

func generateCase(rng *rand.Rand) (int, int, int64, []int64, []int64) {
	n := rng.Intn(6) + 1
	m := rng.Intn(6) + 1
	a := rng.Int63n(100)
	b := make([]int64, n)
	for i := range b {
		b[i] = rng.Int63n(100) + 1
	}
	p := make([]int64, m)
	for i := range p {
		p[i] = rng.Int63n(100) + 1
	}
	return n, m, a, b, p
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m, a, b, p := generateCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, a))
		for j, v := range b {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		for j, v := range p {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		wantR, wantS := expected(n, m, a, append([]int64(nil), b...), append([]int64(nil), p...))
		want := fmt.Sprintf("%d %d", wantR, wantS)
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
