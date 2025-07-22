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

const mod = 95542721

func cube(x int) int {
	return int((int64(x) * int64(x) % mod * int64(x)) % mod)
}

type Query struct {
	t int
	l int
	r int
}

func solveD(n int, arr []int, queries []Query) string {
	a := append([]int(nil), arr...)
	var sb strings.Builder
	for _, q := range queries {
		if q.t == 1 {
			sum := 0
			for i := q.l - 1; i < q.r; i++ {
				sum += a[i]
				if sum >= mod {
					sum %= mod
				}
			}
			fmt.Fprintf(&sb, "%d\n", sum%mod)
		} else {
			for i := q.l - 1; i < q.r; i++ {
				a[i] = cube(a[i])
			}
		}
	}
	return strings.TrimSpace(sb.String())
}

func generateCaseD(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(20)
	}
	q := rng.Intn(10) + 1
	queries := make([]Query, q)
	for i := 0; i < q; i++ {
		t := rng.Intn(2) + 1
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		queries[i] = Query{t, l, r}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", q)
	for _, qu := range queries {
		fmt.Fprintf(&sb, "%d %d %d\n", qu.t, qu.l, qu.r)
	}
	expected := solveD(n, arr, queries)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseD(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
