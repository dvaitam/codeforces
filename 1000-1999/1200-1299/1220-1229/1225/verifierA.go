package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `3 2
5 2
8 8
8 7
4 2
8 1
7 7
1 8
5 4
2 6
1 1
1 9
1 7
4 7
1 9
4 8
8 9
4 6
4 4
8 5
4 7
7 5
8 8
8 1
3 2
6 2
7 3
3 3
4 2
5 2
1 9
8 8
6 3
8 2
7 6
6 4
1 8
8 5
7 4
5 2
3 9
3 4
2 7
2 3
8 4
2 1
1 3
9 9
6 1
7 4
3 5
4 1
1 7
3 6
1 9
7 1
7 6
6 3
3 6
3 2
9 1
1 9
5 1
6 8
2 4
2 9
6 2
1 3
2 2
8 7
2 5
4 3
7 2
8 9
4 1
2 4
4 1
8 5
6 6
1 4
1 2
5 1
1 8
6 9
7 7
5 1
3 7
6 8
7 9
8 9
9 2
6 4
3 1
1 7
5 9
7 4
3 6
4 2
2 9
6 4`

// solve mirrors 1225A.go logic.
func solve(da, db int64) string {
	switch {
	case db-da == 1:
		return fmt.Sprintf("%d %d", db*10-1, db*10)
	case da == db:
		return fmt.Sprintf("%d %d", da*10, db*10+1)
	case da == 9 && db == 1:
		return "99 100"
	default:
		return "-1"
	}
}

type testCase struct {
	da int64
	db int64
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testCase
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid test line: %q", line)
		}
		da, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, err
		}
		db, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, err
		}
		cases = append(cases, testCase{da: da, db: db})
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("%d %d\n", tc.da, tc.db)
}

func runBinary(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	expected := solve(tc.da, tc.db)
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
