package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func expected(n, m int, a, b []int) int {
	oddA, evenA := 0, 0
	for _, v := range a {
		if v%2 == 1 {
			oddA++
		} else {
			evenA++
		}
	}
	oddB, evenB := 0, 0
	for _, v := range b {
		if v%2 == 1 {
			oddB++
		} else {
			evenB++
		}
	}
	min := func(x, y int) int {
		if x < y {
			return x
		}
		return y
	}
	return min(oddA, evenB) + min(evenA, oddB)
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

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	m := rng.Intn(100) + 1
	a := make([]int, n)
	b := make([]int, m)
	for i := range a {
		a[i] = rng.Intn(1000) + 1
	}
	for i := range b {
		b[i] = rng.Intn(1000) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	exp := strconv.Itoa(expected(n, m, a, b))
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// include some deterministic edge cases
	fixed := []struct {
		n, m int
		a, b []int
	}{
		{1, 1, []int{1}, []int{2}},
		{2, 2, []int{1, 2}, []int{3, 4}},
		{3, 1, []int{2, 4, 6}, []int{3}},
		{4, 3, []int{1, 3, 5, 7}, []int{2, 4, 6}},
	}
	idx := 0
	for ; idx < len(fixed); idx++ {
		f := fixed[idx]
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", f.n, f.m))
		for i, v := range f.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for i, v := range f.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expectedOut := strconv.Itoa(expected(f.n, f.m, f.a, f.b))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, input)
			os.Exit(1)
		}
		if got != expectedOut {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", idx+1, expectedOut, got, input)
			os.Exit(1)
		}
	}
	for ; idx < 100; idx++ {
		input, expectedOut := generateCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, input)
			os.Exit(1)
		}
		if got != expectedOut {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", idx+1, expectedOut, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
