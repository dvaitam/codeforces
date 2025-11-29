package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `100
2 4 3
8 6 8
1 1 1
4 4 5
4 5 6
7 11 11
8 10 6
5 5 7
10 7 6
1 1 1
6 4 5
3 8 3
7 3 6
3 4 4
11 9 10
4 3 5
9 11 12
8 6 2
6 6 10
4 2 4
12 8 10
4 1 1
5 10 4
5 7 7
9 10 15
7 2 3
6 4 2
12 6 8
6 9 7
10 7 8
2 5 2
5 4 3
5 6 6
9 6 6
6 10 9
6 3 5
2 5 2
3 4 3
6 5 3
6 12 8
7 11 13
5 9 6
11 7 10
5 7 6
9 10 11
5 5 3
4 4 3
5 2 1
7 5 5
2 7 3
4 4 1
6 2 3
12 8 10
3 5 1
12 8 12
1 3 2
2 1 1
2 4 4
5 7 9
2 3 2
9 9 9
3 8 2
5 7 3
8 10 5
10 5 7
5 6 7
8 7 8
3 8 3
5 5 4
0 2 0
5 11 6
0 6 0
9 6 4
4 2 3
9 6 8
6 8 7
8 11 11
3 3 4
2 8 3
6 6 5
4 1 2
11 8 6
2 3 2
4 2 3
7 5 7
4 6 7
9 10 8
2 3 4
9 7 4
5 5 6
5 5 3
1 4 1
5 4 3
3 7 2
2 1 1
8 11 10
3 5 3
0 2 0
3 2 1
11 7 11`

// Embedded reference logic from 1003B.go.
func buildString(a, b, x int) (string, error) {
	var sb0, sb1 strings.Builder
	for i := 0; i <= x; i++ {
		if i%2 == 0 {
			sb0.WriteByte('0')
			sb1.WriteByte('1')
		} else {
			sb0.WriteByte('1')
			sb1.WriteByte('0')
		}
	}
	alt0 := sb0.String()
	alt1 := sb1.String()

	if x%2 == 1 {
		atemp := a - (x+1)/2
		btemp := b - (x+1)/2
		if atemp >= 0 && btemp >= 0 {
			return strings.Repeat("0", atemp) + alt0 + strings.Repeat("1", btemp), nil
		}
	} else {
		atemp := a - 1 - x/2
		btemp := b - x/2
		if atemp >= 0 && btemp >= 0 {
			prefix := alt0[:len(alt0)-1]
			suffix := alt0[len(alt0)-1:]
			return strings.Repeat("0", atemp) + prefix + strings.Repeat("1", btemp) + suffix, nil
		}

		atemp2 := a - x/2
		btemp2 := b - 1 - x/2
		if atemp2 >= 0 && btemp2 >= 0 {
			prefix := alt1[:len(alt1)-1]
			suffix := alt1[len(alt1)-1:]
			return strings.Repeat("1", btemp2) + prefix + strings.Repeat("0", atemp2) + suffix, nil
		}
	}
	return "", fmt.Errorf("no construction possible for a=%d b=%d x=%d", a, b, x)
}

func runCase(bin string, a, b, x int) error {
	input := fmt.Sprintf("%d %d %d\n", a, b, x)
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	want, err := buildString(a, b, x)
	if err != nil {
		return err
	}
	if got != want {
		return fmt.Errorf("unexpected output: expected %q got %q", want, got)
	}
	return nil
}

type testCase struct {
	a, b, x int
}

func parseTestcases(raw string) ([]testCase, error) {
	fields := strings.Fields(raw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcase data provided")
	}

	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid test count: %w", err)
	}
	expectedFields := 1 + t*3
	if len(fields) != expectedFields {
		return nil, fmt.Errorf("expected %d numbers, found %d", expectedFields, len(fields))
	}

	tests := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		base := 1 + i*3
		a, err1 := strconv.Atoi(fields[base])
		b, err2 := strconv.Atoi(fields[base+1])
		x, err3 := strconv.Atoi(fields[base+2])
		if err := firstErr(err1, err2, err3); err != nil {
			return nil, fmt.Errorf("invalid numbers in case %d: %w", i+1, err)
		}
		tests = append(tests, testCase{a: a, b: b, x: x})
	}
	return tests, nil
}

func firstErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Println("invalid test data:", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc.a, tc.b, tc.x); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
