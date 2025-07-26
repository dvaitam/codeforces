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

type Func struct {
	x1, x2       int
	y1, a, b, y2 int64
}

func eval(f Func, x int64) int64 {
	if x <= int64(f.x1) {
		return f.y1
	}
	if x <= int64(f.x2) {
		return f.a*x + f.b
	}
	return f.y2
}

func expected(n int, funcs []Func, queries [][3]int64) string {
	var last int64
	const mod int64 = 1000000000
	var sb strings.Builder
	for _, q := range queries {
		l := int(q[0])
		r := int(q[1])
		x := (q[2] + last) % mod
		sum := int64(0)
		for j := l - 1; j < r; j++ {
			sum += eval(funcs[j], x)
		}
		sb.WriteString(fmt.Sprintf("%d\n", sum))
		last = sum
	}
	return strings.TrimSpace(sb.String())
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	funcs := make([]Func, n)
	for i := 0; i < n; i++ {
		x1 := rng.Intn(10)
		x2 := x1 + rng.Intn(5) + 1
		y1 := int64(rng.Intn(20))
		a := int64(rng.Intn(5))
		b := int64(rng.Intn(20))
		y2 := int64(rng.Intn(20))
		funcs[i] = Func{x1, x2, y1, a, b, y2}
	}
	m := rng.Intn(5) + 1
	queries := make([][3]int64, m)
	for i := 0; i < m; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		x := int64(rng.Intn(50))
		queries[i] = [3]int64{int64(l), int64(r), x}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		f := funcs[i]
		sb.WriteString(fmt.Sprintf("%d %d %d %d %d %d\n", f.x1, f.x2, f.y1, f.a, f.b, f.y2))
	}
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i := 0; i < m; i++ {
		q := queries[i]
		sb.WriteString(fmt.Sprintf("%d %d %d\n", q[0], q[1], q[2]))
	}
	return sb.String(), expected(n, funcs, queries)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "test %d failed:\nexpected:\n%s\ngot:\n%s\ninput:\n%s", i+1, exp, strings.TrimSpace(got), in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
