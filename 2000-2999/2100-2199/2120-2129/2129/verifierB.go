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

const refPath = "2000-2999/2100-2199/2120-2129/2129/2129B.go"

type testCase struct {
	p []int
}

func main() {
	if len(os.Args) != 2 {
		if len(os.Args) == 3 && os.Args[1] == "--" {
			// allow: go run verifierB.go -- /path/to/binary
			os.Args = []string{os.Args[0], os.Args[2]}
		} else {
			fmt.Println("usage: go run verifierB.go /path/to/binary")
			os.Exit(1)
		}
	}
	bin := os.Args[1]

	tests := buildTests()
	input := renderInput(tests)

	expectRaw, err := runBinary(refPath, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\ninput:\n%s\n", err, input)
		os.Exit(1)
	}
	actualRaw, err := runBinary(bin, input)
	if err != nil {
		fmt.Printf("runtime error: %v\ninput:\n%s\n", err, input)
		os.Exit(1)
	}

	expVals, err := parseOutputs(expectRaw, len(tests))
	if err != nil {
		fmt.Printf("unable to parse reference output: %v\noutput:\n%s\n", err, expectRaw)
		os.Exit(1)
	}
	actVals, err := parseOutputs(actualRaw, len(tests))
	if err != nil {
		fmt.Printf("unable to parse contestant output: %v\noutput:\n%s\n", err, actualRaw)
		os.Exit(1)
	}

	for i := range tests {
		if expVals[i] != actVals[i] {
			fmt.Printf("case %d mismatch: expected %d got %d\nperm=%v\nfull input:\n%s\n", i+1, expVals[i], actVals[i], tests[i].p, input)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func buildTests() []testCase {
	rand.Seed(time.Now().UnixNano())
	var tests []testCase
	sumN := 0

	// Exhaustive for n up to 4
	for n := 2; n <= 4; n++ {
		for _, perm := range allPerms(n) {
			tests = append(tests, testCase{p: perm})
			sumN += n
		}
	}

	// Random mid-size cases
	for i := 0; i < 40; i++ {
		n := 2 + rand.Intn(30) // 2..31
		if sumN+n > 5000 {
			break
		}
		tests = append(tests, testCase{p: randomPerm(n)})
		sumN += n
	}

	// One larger stress case
	n := 1000
	if sumN+n <= 5000 {
		tests = append(tests, testCase{p: randomPerm(n)})
		sumN += n
	}

	// Fill remaining budget with smaller random cases
	for sumN < 5000 {
		n := 2 + rand.Intn(20)
		if sumN+n > 5000 {
			break
		}
		tests = append(tests, testCase{p: randomPerm(n)})
		sumN += n
	}

	return tests
}

func renderInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", len(tc.p))
		for i, v := range tc.p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutputs(out string, t int) ([]int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	vals := make([]int, 0, t)
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		var v int
		if _, err := fmt.Sscan(ln, &v); err != nil {
			return nil, fmt.Errorf("could not parse line %q: %v", ln, err)
		}
		vals = append(vals, v)
	}
	if len(vals) != t {
		return nil, fmt.Errorf("expected %d values, got %d", t, len(vals))
	}
	return vals, nil
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func allPerms(n int) [][]int {
	used := make([]bool, n)
	cur := make([]int, n)
	var res [][]int
	var dfs func(int)
	dfs = func(pos int) {
		if pos == n {
			tmp := make([]int, n)
			copy(tmp, cur)
			res = append(res, tmp)
			return
		}
		for i := 0; i < n; i++ {
			if used[i] {
				continue
			}
			used[i] = true
			cur[pos] = i + 1
			dfs(pos + 1)
			used[i] = false
		}
	}
	dfs(0)
	return res
}

func randomPerm(n int) []int {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i + 1
	}
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}
	return p
}
