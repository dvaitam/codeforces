package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const K = 45

type pair struct {
	val, idx int
}

func solveE(r io.Reader) string {
	reader := bufio.NewReader(r)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return ""
	}
	a := make([]int, n+2)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	p2 := make([]float64, K)
	p2[0] = 1.0
	for i := 1; i < K; i++ {
		p2[i] = p2[i-1] * 2.0
	}
	v := make([]pair, n)
	for i := 1; i <= n; i++ {
		v[i-1] = pair{a[i], i}
	}
	sort.Slice(v, func(i, j int) bool { return v[i].val < v[j].val })
	tlft := make([]int, n+2)
	trgt := make([]int, n+2)
	for i := 0; i <= n+1; i++ {
		tlft[i] = i - 1
		trgt[i] = i + 1
	}
	tlft[0] = 0
	trgt[n] = n + 1
	ans := 0.0
	for _, pr := range v {
		i := pr.idx
		tl := 0.0
		tr := 0.0
		pl := i
		prr := i
		for j := 0; j < K; j++ {
			if pl != 0 {
				dist := float64(pl - tlft[pl])
				tl += dist / p2[j]
				pl = tlft[pl]
			}
			if prr <= n {
				dist := float64(trgt[prr] - prr)
				tr += dist / p2[j]
				prr = trgt[prr]
			}
		}
		left := tlft[i]
		right := trgt[i]
		trgt[left] = right
		tlft[right] = left
		ans += tl * tr * float64(a[i]) * 0.5
	}
	result := ans / (float64(n) * float64(n))
	return fmt.Sprintf("%.12f\n", result)
}

func runCaseE(bin, in, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expect)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func genCaseE(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(rng.Intn(100) + 1))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCaseE(rng)
		expect := solveE(strings.NewReader(in))
		if err := runCaseE(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
