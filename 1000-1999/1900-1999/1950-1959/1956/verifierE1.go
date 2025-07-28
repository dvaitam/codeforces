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

func step(a []int) bool {
	changed := false
	n := len(a)
	for i := 0; i < n-1; i++ {
		v := a[i+1] - a[i]
		if v < 0 {
			v = 0
		}
		if v != a[i+1] {
			a[i+1] = v
			changed = true
		}
	}
	v := a[0] - a[n-1]
	if v < 0 {
		v = 0
	}
	if v != a[0] {
		a[0] = v
		changed = true
	}
	return changed
}

func solveCase(a []int) string {
	for iter := 0; iter < 200000; iter++ {
		if !step(a) {
			break
		}
	}
	var indices []int
	for i, v := range a {
		if v > 0 {
			indices = append(indices, i+1)
		}
	}
	var out bytes.Buffer
	fmt.Fprintf(&out, "%d", len(indices))
	if len(indices) > 0 {
		out.WriteByte('\n')
		for i, idx := range indices {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprintf(&out, "%d", idx)
		}
	}
	return strings.TrimSpace(out.String())
}

func genCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(3) + 1
	var in bytes.Buffer
	var out bytes.Buffer
	fmt.Fprintf(&in, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(5) + 1
		fmt.Fprintf(&in, "%d\n", n)
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(11)
			if j > 0 {
				in.WriteByte(' ')
			}
			fmt.Fprintf(&in, "%d", a[j])
		}
		in.WriteByte('\n')
		out.WriteString(solveCase(append([]int(nil), a...)))
		if i+1 < t {
			out.WriteByte('\n')
		}
	}
	return in.String(), strings.TrimSpace(out.String())
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
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
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
		if got != exp {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, in, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
