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

type testCaseB struct {
	input    string
	expected string
}

func runCandidate(bin, input string) (string, error) {
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

func solveB(arr []int) string {
	res := make([]byte, len(arr))
	var first, last int
	have := false
	flag := false
	for i, x := range arr {
		if !have {
			first, last = x, x
			have = true
			res[i] = '1'
			continue
		}
		if !flag {
			if x >= last {
				last = x
				res[i] = '1'
			} else if x <= first {
				last = x
				flag = true
				res[i] = '1'
			} else {
				res[i] = '0'
			}
		} else {
			if x >= last && x <= first {
				last = x
				res[i] = '1'
			} else {
				res[i] = '0'
			}
		}
	}
	return string(res)
}

func generateCaseB(rng *rand.Rand) testCaseB {
	t := rng.Intn(3) + 1
	var in strings.Builder
	var out strings.Builder
	in.WriteString(fmt.Sprintf("%d\n", t))
	for j := 0; j < t; j++ {
		q := rng.Intn(20) + 1
		arr := make([]int, q)
		in.WriteString(fmt.Sprintf("%d\n", q))
		for i := 0; i < q; i++ {
			if i > 0 {
				in.WriteByte(' ')
			}
			arr[i] = rng.Intn(50)
			in.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		in.WriteByte('\n')
		out.WriteString(solveB(arr))
		out.WriteByte('\n')
	}
	return testCaseB{input: in.String(), expected: out.String()}
}

func runCaseB(bin string, tc testCaseB) error {
	got, err := runCandidate(bin, tc.input)
	if err != nil {
		return err
	}
	got = strings.TrimSpace(got)
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCaseB{generateCaseB(rng)}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCaseB(rng))
	}
	for i, tc := range cases {
		if err := runCaseB(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
