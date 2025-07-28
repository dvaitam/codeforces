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

func expectedAnswer(A, B int64, a, b []int64) string {
	var total int64
	var maxAttack int64
	for i := range a {
		hits := (b[i] + A - 1) / A
		total += hits * a[i]
		if a[i] > maxAttack {
			maxAttack = a[i]
		}
	}
	if B > total-maxAttack {
		return "YES\n"
	}
	return "NO\n"
}

func generateCase(rng *rand.Rand) (string, string) {
	A := int64(rng.Intn(1_000_000) + 1)
	B := int64(rng.Intn(1_000_000) + 1)
	n := rng.Intn(5) + 1
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = int64(rng.Intn(1_000_000) + 1)
	}
	for i := 0; i < n; i++ {
		b[i] = int64(rng.Intn(1_000_000) + 1)
	}

	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.FormatInt(A, 10))
	sb.WriteByte(' ')
	sb.WriteString(strconv.FormatInt(B, 10))
	sb.WriteByte(' ')
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')

	expected := expectedAnswer(A, B, a, b)
	return sb.String(), expected
}

func runCase(exe string, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if outStr != exp {
		return fmt.Errorf("expected %q got %q", exp, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type caseData struct {
		A int64
		B int64
		a []int64
		b []int64
	}
	edges := []caseData{
		{A: 1, B: 1, a: []int64{1}, b: []int64{1}},
		{A: 5, B: 5, a: []int64{3}, b: []int64{5}},
		{A: 5, B: 4, a: []int64{5}, b: []int64{25}},
		{A: 10, B: 100, a: []int64{1, 100}, b: []int64{50, 100}},
		{A: 1000000, B: 1000000, a: []int64{1000000}, b: []int64{1000000}},
	}
	for i, c := range edges {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.FormatInt(c.A, 10))
		sb.WriteByte(' ')
		sb.WriteString(strconv.FormatInt(c.B, 10))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(len(c.a)))
		sb.WriteByte('\n')
		for j, v := range c.a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for j, v := range c.b {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := expectedAnswer(c.A, c.B, c.a, c.b)
		if err := runCase(exe, input, expected); err != nil {
			fmt.Fprintf(os.Stderr, "edge case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}

	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "random case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
