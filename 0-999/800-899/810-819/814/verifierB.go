package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func validate(n int, a, b []int, output string) error {
	fields := strings.Fields(output)
	if len(fields) != n {
		return fmt.Errorf("expected %d values, got %d", n, len(fields))
	}
	p := make([]int, n)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid integer %q", f)
		}
		p[i] = v
	}

	freq := make([]int, n+1)
	for _, v := range p {
		if v < 1 || v > n {
			return fmt.Errorf("value %d out of range [1,%d]", v, n)
		}
		freq[v]++
	}
	for v := 1; v <= n; v++ {
		if freq[v] != 1 {
			return fmt.Errorf("value %d appears %d times (not a permutation)", v, freq[v])
		}
	}

	da, db := 0, 0
	for i := 0; i < n; i++ {
		if p[i] != a[i] {
			da++
		}
		if p[i] != b[i] {
			db++
		}
	}
	if da != 1 {
		return fmt.Errorf("differs from a in %d positions (need exactly 1)", da)
	}
	if db != 1 {
		return fmt.Errorf("differs from b in %d positions (need exactly 1)", db)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	idx := 0
	for {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			break
		}
		a := make([]int, n)
		for i := range a {
			fmt.Fscan(reader, &a[i])
		}
		b := make([]int, n)
		for i := range b {
			fmt.Fscan(reader, &b[i])
		}
		idx++

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
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
		input := sb.String()

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var outBuf strings.Builder
		var errBuf strings.Builder
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}

		got := strings.TrimSpace(outBuf.String())
		if err := validate(n, a, b, got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sgot: %s\nreason: %v\n", idx, input, got, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
