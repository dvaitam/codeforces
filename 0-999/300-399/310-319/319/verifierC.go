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

type Line struct{ m, c int64 }

func eval(l Line, x int64) int64 { return l.m*x + l.c }

func isBad(l1, l2, l3 Line) bool {
	return (l2.c-l1.c)*(l2.m-l3.m) >= (l3.c-l2.c)*(l1.m-l2.m)
}

func expected(a, b []int64) int64 {
	n := len(a)
	dp := make([]int64, n)
	hull := []Line{{m: b[0], c: dp[0]}}
	pos := 0
	for i := 1; i < n; i++ {
		x := a[i]
		for pos+1 < len(hull) && eval(hull[pos], x) >= eval(hull[pos+1], x) {
			pos++
		}
		dp[i] = eval(hull[pos], x)
		nl := Line{m: b[i], c: dp[i]}
		for len(hull) >= 2 && isBad(hull[len(hull)-2], hull[len(hull)-1], nl) {
			hull = hull[:len(hull)-1]
			if pos >= len(hull) {
				pos = len(hull) - 1
			}
		}
		hull = append(hull, nl)
	}
	return dp[n-1]
}

func runCase(bin string, a, b []int64) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(a)))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := expected(a, b)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func generateCase(rng *rand.Rand) ([]int64, []int64) {
	n := rng.Intn(50) + 1
	a := make([]int64, n)
	a[0] = 1
	for i := 1; i < n; i++ {
		a[i] = a[i-1] + int64(rng.Intn(5)+1)
	}
	b := make([]int64, n)
	b[n-1] = 0
	cur := int64(rng.Intn(1000) + 2000)
	for i := 0; i < n-1; i++ {
		b[i] = cur
		cur -= int64(rng.Intn(5) + 1)
	}
	return a, b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	a0 := []int64{1}
	b0 := []int64{0}
	casesA := [][]int64{a0}
	casesB := [][]int64{b0}
	for i := 0; i < 99; i++ {
		a, b := generateCase(rng)
		casesA = append(casesA, a)
		casesB = append(casesB, b)
	}
	for i := 0; i < len(casesA); i++ {
		if err := runCase(bin, casesA[i], casesB[i]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
