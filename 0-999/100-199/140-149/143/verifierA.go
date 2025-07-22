package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	r1, r2, c1, c2, d1, d2 int
	solvable               bool
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("time limit")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func hasSolution(r1, r2, c1, c2, d1, d2 int) bool {
	for a := 1; a <= 9; a++ {
		for b := 1; b <= 9; b++ {
			if a == b || a+b != r1 {
				continue
			}
			for c := 1; c <= 9; c++ {
				if c == a || c == b || a+c != c1 {
					continue
				}
				d := r2 - c
				if d < 1 || d > 9 || d == a || d == b || d == c {
					continue
				}
				if b+d != c2 {
					continue
				}
				if a+d != d1 {
					continue
				}
				if b+c != d2 {
					continue
				}
				return true
			}
		}
	}
	return false
}

func generateCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 100)
	for i := range cases {
		if rng.Intn(4) == 0 { // generate unsolvable case 25%
			for {
				r1 := rng.Intn(20) + 1
				r2 := rng.Intn(20) + 1
				c1 := rng.Intn(20) + 1
				c2 := rng.Intn(20) + 1
				d1 := rng.Intn(20) + 1
				d2 := rng.Intn(20) + 1
				if !hasSolution(r1, r2, c1, c2, d1, d2) {
					cases[i] = testCase{r1, r2, c1, c2, d1, d2, false}
					break
				}
			}
		} else { // solvable case
			digits := rng.Perm(9)[:4]
			a, b, c, d := digits[0]+1, digits[1]+1, digits[2]+1, digits[3]+1
			r1 := a + b
			r2 := c + d
			c1 := a + c
			c2 := b + d
			d1 := a + d
			d2 := b + c
			cases[i] = testCase{r1, r2, c1, c2, d1, d2, true}
		}
	}
	return cases
}

func checkOutput(out string, tc testCase) error {
	out = strings.TrimSpace(out)
	if out == "-1" {
		if tc.solvable {
			return fmt.Errorf("expected solution but got -1")
		}
		if hasSolution(tc.r1, tc.r2, tc.c1, tc.c2, tc.d1, tc.d2) {
			return fmt.Errorf("solution exists but program output -1")
		}
		return nil
	}
	lines := strings.Split(out, "\n")
	if len(lines) != 2 {
		return fmt.Errorf("expected 2 lines, got %d", len(lines))
	}
	var a, b, c, d int
	if _, err := fmt.Sscan(lines[0], &a, &b); err != nil {
		return fmt.Errorf("cannot parse first line: %v", err)
	}
	if _, err := fmt.Sscan(lines[1], &c, &d); err != nil {
		return fmt.Errorf("cannot parse second line: %v", err)
	}
	vals := []int{a, b, c, d}
	used := map[int]bool{}
	for _, v := range vals {
		if v < 1 || v > 9 {
			return fmt.Errorf("value %d out of range", v)
		}
		if used[v] {
			return fmt.Errorf("values not distinct")
		}
		used[v] = true
	}
	if a+b != tc.r1 || c+d != tc.r2 || a+c != tc.c1 || b+d != tc.c2 || a+d != tc.d1 || b+c != tc.d2 {
		return fmt.Errorf("values do not satisfy equations")
	}
	if !tc.solvable {
		return fmt.Errorf("expected -1 for unsolvable case")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCases()
	for i, tc := range cases {
		input := fmt.Sprintf("%d %d %d %d %d %d\n", tc.r1, tc.r2, tc.c1, tc.c2, tc.d1, tc.d2)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if err := checkOutput(out, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%soutput:%s\n", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
