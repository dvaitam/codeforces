package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	desc  string
	input string
	t     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierJ.go /path/to/binary")
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
		exp, err := parseOutputs(refOut, tc.t)
		if err != nil {
			fmt.Printf("Failed to parse reference output on test %d (%s): %v\nOutput:\n%s\n", i+1, tc.desc, err, refOut)
			os.Exit(1)
		}

		out, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, out)
			os.Exit(1)
		}
		got, err := parseOutputs(out, tc.t)
		if err != nil {
			fmt.Printf("Failed to parse output on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, out)
			os.Exit(1)
		}
		for j := 0; j < tc.t; j++ {
			if got[j] != exp[j] {
				fmt.Printf("Wrong answer on test %d (%s) case %d: expected %d got %d\nInput:\n%s", i+1, tc.desc, j+1, exp[j], got[j], tc.input)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	path := "./ref1267J.bin"
	cmd := exec.Command("go", "build", "-o", path, "1267J.go")
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

func parseOutputs(out string, t int) ([]int, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(tokens))
	}
	res := make([]int, t)
	for i, tok := range tokens {
		val, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		res[i] = val
	}
	return res, nil
}

func generateTests() []testCase {
	var tests []testCase
	add := func(desc string, cases [][]int) {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
		for _, arr := range cases {
			sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
			for i, v := range arr {
				if i > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(v))
			}
			sb.WriteByte('\n')
		}
		tests = append(tests, testCase{
			desc:  desc,
			input: sb.String(),
			t:     len(cases),
		})
	}

	add("single-icon", [][]int{{1}})
	add("two-categories", [][]int{{1, 1, 2, 2}})
	add("all-same", [][]int{{3, 3, 3, 3, 3, 3}})

	rng := rand.New(rand.NewSource(123456789))
	for len(tests) < 60 {
		t := rng.Intn(5) + 1
		cases := make([][]int, t)
		for i := 0; i < t; i++ {
			n := rng.Intn(2000) + 1
			arr := make([]int, n)
			for j := 0; j < n; j++ {
				arr[j] = rng.Intn(n) + 1
			}
			cases[i] = arr
		}
		add(fmt.Sprintf("random-%d", len(tests)), cases)
	}

	// Large stress tests
	largeN := 200000
	largeCase := make([]int, largeN)
	for i := 0; i < largeN; i++ {
		largeCase[i] = (i % 1000) + 1
	}
	add("large-single", [][]int{largeCase})

	tMany := 10
	cases := make([][]int, tMany)
	for i := 0; i < tMany; i++ {
		n := 50000
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(n) + 1
		}
		cases[i] = arr
	}
	add("many-large", cases)

	return tests
}
