package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runProg(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTest(rng *rand.Rand) (string, []int) {
	n := rng.Intn(6) + 1
	p := rng.Perm(n)
	a := make([]int, n)
	present := make([]bool, n+1)
	mex := 0
	for i := 0; i < n; i++ {
		for present[mex] {
			mex++
		}
		a[i] = mex - p[i]
		present[p[i]] = true
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "1\n%d\n", n)
	for i, v := range a {
		if i > 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprint(&buf, v)
	}
	buf.WriteByte('\n')
	return buf.String(), a
}

func verifyOutput(out string, a []int) error {
	fields := strings.Fields(out)
	if len(fields) != len(a) {
		return fmt.Errorf("expected %d numbers, got %d", len(a), len(fields))
	}

	n := len(a)
	p := make([]int, n)
	used := make([]bool, n)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("non-integer token %q", f)
		}
		if v < 0 || v >= n {
			return fmt.Errorf("value out of range at index %d: %d", i, v)
		}
		if used[v] {
			return fmt.Errorf("duplicate value %d", v)
		}
		used[v] = true
		p[i] = v
	}

	present := make([]bool, n+1)
	mex := 0
	for i := 0; i < n; i++ {
		for present[mex] {
			mex++
		}
		if got := mex - p[i]; got != a[i] {
			return fmt.Errorf("constraint mismatch at index %d: expected a[%d]=%d, got %d", i, i, a[i], got)
		}
		present[p[i]] = true
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	rng := rand.New(rand.NewSource(43))
	for i := 0; i < 100; i++ {
		test, a := genTest(rng)
		got, err := runProg(target, test)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d execution error: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := verifyOutput(got, a); err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\noutput:%s\nreason:%v\n", i+1, test, got, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
