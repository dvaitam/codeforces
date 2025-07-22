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

func expected(a, b []int) []int {
	n := len(a)
	invA := make([]int, n+1)
	for i := 1; i <= n; i++ {
		invA[a[i-1]] = i
	}
	res := make([]int, n)
	for i := 1; i <= n; i++ {
		res[i-1] = b[invA[i]-1]
	}
	return res
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

func runCase(bin string, arrA, arrB []int) error {
	n := len(arrA)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arrA {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range arrB {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	exp := expected(arrA, arrB)
	out, err := run(bin, sb.String())
	if err != nil {
		return err
	}
	fields := strings.Fields(out)
	if len(fields) != n {
		return fmt.Errorf("expected %d numbers got %d", n, len(fields))
	}
	for i := 0; i < n; i++ {
		val, err := strconv.Atoi(fields[i])
		if err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		if val != exp[i] {
			return fmt.Errorf("expected %v got %v", exp, fields)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []struct {
		a, b []int
	}{{[]int{1}, []int{1}}, {[]int{1, 2}, []int{2, 1}}}

	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 1
		arrA := rng.Perm(n)
		for j := range arrA {
			arrA[j]++
		}
		arrB := rng.Perm(n)
		for j := range arrB {
			arrB[j]++
		}
		cases = append(cases, struct{ a, b []int }{arrA, arrB})
	}

	for idx, c := range cases {
		if err := runCase(bin, c.a, c.b); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
