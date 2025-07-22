package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type TestCase struct {
	n   int
	s   string
	ans int
}

func compute(n int, s string) int {
	for i := 0; i < n; i++ {
		if s[i] == '0' {
			return i + 1
		}
	}
	return n
}

func genCases() []TestCase {
	var cases []TestCase
	for n := 1; n <= 7; n++ {
		for mask := 0; mask < (1 << n); mask++ {
			b := make([]byte, n)
			for i := 0; i < n; i++ {
				if mask&(1<<i) != 0 {
					b[i] = '1'
				} else {
					b[i] = '0'
				}
			}
			s := string(b)
			cases = append(cases, TestCase{n, s, compute(n, s)})
		}
	}
	return cases
}

func runCase(bin string, tc TestCase) error {
	input := fmt.Sprintf("%d\n%s\n", tc.n, tc.s)
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	output := strings.TrimSpace(out.String())
	val, err := strconv.Atoi(output)
	if err != nil {
		return fmt.Errorf("invalid output: %s", output)
	}
	if val != tc.ans {
		return fmt.Errorf("expected %d got %d", tc.ans, val)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) >= 3 {
		bin = os.Args[2]
	}
	cases := genCases()
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%d\n%s\n", i+1, err, tc.n, tc.s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
