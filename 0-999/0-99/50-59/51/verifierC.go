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

func solveC(a []int64) (string, []float64) {
	n := len(a)
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	st := [3]int{}
	test := func(x int64) bool {
		b := 0
		for i := 0; i < 3; i++ {
			if b >= n {
				return false
			}
			st[i] = b
			key := a[b] + x
			pos := sort.Search(n, func(j int) bool { return a[j] > key })
			b = pos
			if b == n {
				return false
			}
		}
		return true
	}
	var l, r int64 = 0, 1000000000
	for l < r {
		mid := l + (r-l)/2
		if test(mid) {
			l = mid + 1
		} else {
			r = mid
		}
	}
	test(l)
	d := float64(l) / 2.0
	pos := make([]float64, 3)
	for i := 0; i < 3; i++ {
		pos[i] = float64(a[st[i]]) + d
	}
	res := fmt.Sprintf("%.6f", d)
	return res, pos
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 3
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = int64(rng.Intn(1000))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	dStr, pos := solveC(append([]int64(nil), arr...))
	var exp strings.Builder
	exp.WriteString(dStr)
	exp.WriteByte('\n')
	for _, p := range pos {
		exp.WriteString(fmt.Sprintf("%.6f ", p))
	}
	exp.WriteByte('\n')
	return sb.String(), strings.TrimSpace(exp.String())
}

func run(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
