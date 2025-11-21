package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input string
	q     int
	desc  string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]
	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Printf("Reference runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, refOut)
			os.Exit(1)
		}
		exp, err := parseAnswers(refOut, tc.q)
		if err != nil {
			fmt.Printf("Failed to parse reference output on test %d (%s): %v\nOutput:\n%s\n", i+1, tc.desc, err, refOut)
			os.Exit(1)
		}
		out, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, out)
			os.Exit(1)
		}
		got, err := parseAnswers(out, tc.q)
		if err != nil {
			fmt.Printf("Failed to parse output on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, out)
			os.Exit(1)
		}
		for j := 0; j < tc.q; j++ {
			if got[j] != exp[j] {
				fmt.Printf("Wrong answer on test %d (%s): query %d expected %s got %s\nInput:\n%s", i+1, tc.desc, j+1, exp[j], got[j], tc.input)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	path := "./ref1221A.bin"
	cmd := exec.Command("go", "build", "-o", path, "1221A.go")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("build failed: %v\n%s", err, stderr.String())
	}
	return path, nil
}

func runProgram(bin, input string) (string, error) {
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
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func parseAnswers(out string, q int) ([]string, error) {
	tokens := strings.Fields(out)
	if len(tokens) != q {
		return nil, fmt.Errorf("expected %d answers, got %d", q, len(tokens))
	}
	ans := make([]string, q)
	for i, t := range tokens {
		upper := strings.ToUpper(t)
		if upper != "YES" && upper != "NO" {
			return nil, fmt.Errorf("invalid answer %q", t)
		}
		ans[i] = upper
	}
	return ans, nil
}

func generateTests() []testCase {
	var tests []testCase
	add := func(desc string, q int, data [][]int) {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", q)
		for _, arr := range data {
			fmt.Fprintf(&sb, "%d\n", len(arr))
			for i, v := range arr {
				if i > 0 {
					sb.WriteByte(' ')
				}
				fmt.Fprintf(&sb, "%d", v)
			}
			sb.WriteByte('\n')
		}
		tests = append(tests, testCase{
			input: sb.String(),
			q:     q,
			desc:  desc,
		})
	}

	add("single-small-no", 1, [][]int{{1}})
	add("single-yes", 1, [][]int{{2048}})
	add("two-queries", 2, [][]int{{1024, 1024}, {1, 2, 4, 8}})

	rng := rand.New(rand.NewSource(123456789))
	for len(tests) < 80 {
		q := rng.Intn(5) + 1
		data := make([][]int, q)
		for i := 0; i < q; i++ {
			n := rng.Intn(6) + 1
			arr := make([]int, n)
			for j := 0; j < n; j++ {
				pow := rng.Intn(30)
				arr[j] = 1 << uint(pow%30)
				if arr[j] > 1<<29 {
					arr[j] = 1 << 29
				}
			}
			data[i] = arr
		}
		add(fmt.Sprintf("random-small-%d", len(tests)), q, data)
	}

	for len(tests) < 120 {
		q := 100
		data := make([][]int, q)
		for i := 0; i < q; i++ {
			n := 100
			arr := make([]int, n)
			for j := 0; j < n; j++ {
				if rng.Intn(2) == 0 {
					arr[j] = 2048
				} else {
					arr[j] = 1 << uint(rng.Intn(11))
				}
			}
			data[i] = arr
		}
		add(fmt.Sprintf("random-max-%d", len(tests)), q, data)
	}

	return tests
}
