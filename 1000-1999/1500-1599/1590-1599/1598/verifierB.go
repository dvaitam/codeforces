package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func expectedB(n int, mat [][]int) string {
	for d1 := 0; d1 < 5; d1++ {
		for d2 := d1 + 1; d2 < 5; d2++ {
			aOnly, bOnly := 0, 0
			valid := true
			for i := 0; i < n; i++ {
				a := mat[i][d1]
				b := mat[i][d2]
				if a == 0 && b == 0 {
					valid = false
					break
				}
				if a == 1 && b == 0 {
					aOnly++
				}
				if a == 0 && b == 1 {
					bOnly++
				}
			}
			if valid && aOnly <= n/2 && bOnly <= n/2 {
				return "YES"
			}
		}
	}
	return "NO"
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(10)*2 + 2 // even 2..20
		mat := make([][]int, n)
		for j := 0; j < n; j++ {
			mat[j] = make([]int, 5)
			ok := false
			for k := 0; k < 5; k++ {
				val := rng.Intn(2)
				if val == 1 {
					ok = true
				}
				mat[j][k] = val
			}
			if !ok {
				mat[j][rng.Intn(5)] = 1
			}
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			for k := 0; k < 5; k++ {
				if k > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprint(mat[j][k]))
			}
			sb.WriteByte('\n')
		}
		exp := expectedB(n, mat)
		cases[i] = testCase{input: sb.String(), expected: exp}
	}
	return cases
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if out != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
