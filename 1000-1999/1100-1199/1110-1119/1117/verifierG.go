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

func maxPos(arr []int, l, r int) int {
	idx := l
	for i := l + 1; i <= r; i++ {
		if arr[i] > arr[idx] {
			idx = i
		}
	}
	return idx
}

func solve(arr []int, l, r int) int64 {
	if l > r {
		return 0
	}
	m := maxPos(arr, l, r)
	return int64(r-l+1) + solve(arr, l, m-1) + solve(arr, m+1, r)
}

func expectedG(arr []int, l, r int) int64 {
	return solve(arr, l, r)
}

func generateCase(rng *rand.Rand) (string, []int64) {
	n := rng.Intn(7) + 1
	q := rng.Intn(5) + 1
	perm := rng.Perm(n)
	for i := range perm {
		perm[i]++
	}
	l := make([]int, q)
	r := make([]int, q)
	for i := 0; i < q; i++ {
		a := rng.Intn(n)
		b := rng.Intn(n)
		if a > b {
			a, b = b, a
		}
		l[i] = a + 1
		r[i] = b + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	for i, v := range perm {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range l {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range r {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	answers := make([]int64, q)
	for i := 0; i < q; i++ {
		answers[i] = expectedG(perm, l[i]-1, r[i]-1)
	}
	return sb.String(), answers
}

func runCase(bin, input string, exp []int64) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotFields := strings.Fields(out.String())
	if len(gotFields) != len(exp) {
		return fmt.Errorf("expected %d numbers got %d", len(exp), len(gotFields))
	}
	for i, gf := range gotFields {
		var g int64
		fmt.Sscan(gf, &g)
		if g != exp[i] {
			return fmt.Errorf("at %d expected %d got %d", i+1, exp[i], g)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
