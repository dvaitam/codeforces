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

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randomCase(rng *rand.Rand) (string, [][]int) {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	cases := make([][]int, t)
	for c := 0; c < t; c++ {
		n := rng.Intn(49) + 2
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			nums[i] = rng.Intn(1000000000) + 1
		}
		cases[c] = nums
		fmt.Fprintf(&sb, "%d\n", n)
		for i, v := range nums {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
	}
	return sb.String(), cases
}

func validate(a []int, bLine string) error {
	n := len(a)
	fields := strings.Fields(bLine)
	if len(fields) != n {
		return fmt.Errorf("expected %d values, got %d", n, len(fields))
	}
	b := make([]int, n)
	for i := 0; i < n; i++ {
		val, err := strconv.Atoi(fields[i])
		if err != nil {
			return fmt.Errorf("invalid integer %q at position %d", fields[i], i)
		}
		b[i] = val
	}

	// Check 1 <= b[i] <= 10^9
	for i, v := range b {
		if v < 1 || v > 1000000000 {
			return fmt.Errorf("b[%d] = %d out of range [1, 10^9]", i, v)
		}
	}

	// Check divisibility: consecutive b[i], b[i+1] must have one dividing the other
	for i := 0; i+1 < n; i++ {
		if b[i]%b[i+1] != 0 && b[i+1]%b[i] != 0 {
			return fmt.Errorf("b[%d]=%d and b[%d]=%d: neither divides the other", i, b[i], i+1, b[i+1])
		}
	}

	// Check 2 * sum(|a[i]-b[i]|) <= S
	S := int64(0)
	diff := int64(0)
	for i := 0; i < n; i++ {
		S += int64(a[i])
		d := int64(a[i]) - int64(b[i])
		if d < 0 {
			d = -d
		}
		diff += d
	}
	if 2*diff > S {
		return fmt.Errorf("2*sum(|a-b|) = %d > S = %d", 2*diff, S)
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for iter := 0; iter < 100; iter++ {
		input, cases := randomCase(rng)
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", iter+1, err, input)
			os.Exit(1)
		}
		lines := strings.Split(got, "\n")
		if len(lines) < len(cases) {
			fmt.Fprintf(os.Stderr, "test %d: expected %d output lines, got %d\noutput:\n%s\n", iter+1, len(cases), len(lines), got)
			os.Exit(1)
		}
		for c, a := range cases {
			if err := validate(a, lines[c]); err != nil {
				fmt.Fprintf(os.Stderr, "test %d, case %d: %v\ninput:\n%s\noutput:\n%s\n", iter+1, c+1, err, input, got)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}
