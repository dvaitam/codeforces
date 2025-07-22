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

func solveCase(n int, arr []int, ops [][]int) string {
	a := make([]int, n)
	copy(a, arr)
	var sb strings.Builder
	for _, op := range ops {
		t := op[0]
		l := op[1]
		r := op[2]
		if t == 1 {
			sum := 0
			for i := l - 1; i < r; i++ {
				sum += a[i]
			}
			fmt.Fprintf(&sb, "%d\n", sum)
		} else {
			x := op[3]
			for i := l - 1; i < r; i++ {
				a[i] ^= x
			}
		}
	}
	return sb.String()
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(20)
	}
	m := rng.Intn(10) + 1
	ops := make([][]int, m)
	for i := 0; i < m; i++ {
		t := rng.Intn(2) + 1
		l := rng.Intn(n) + 1
		r := l + rng.Intn(n-l+1)
		if t == 1 {
			ops[i] = []int{t, l, r, 0}
		} else {
			x := rng.Intn(50) + 1
			ops[i] = []int{t, l, r, x}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", m)
	for _, op := range ops {
		if op[0] == 1 {
			fmt.Fprintf(&sb, "1 %d %d\n", op[1], op[2])
		} else {
			fmt.Fprintf(&sb, "2 %d %d %d\n", op[1], op[2], op[3])
		}
	}
	out := solveCase(n, arr, ops)
	return sb.String(), out
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	if strings.TrimSpace(buf.String()) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected \n%s\ngot \n%s", exp, buf.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
