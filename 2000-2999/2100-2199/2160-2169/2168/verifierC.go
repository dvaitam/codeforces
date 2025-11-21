package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const maxValue = 1 << 15

type caseData struct {
	x   int
	set []bool
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}

	candidate := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(2168))

	cases := generateTests(rng)
	firstInput := buildFirstInput(cases)

	firstOut, err := runCandidate(candidate, firstInput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed on first run: %v\n%s\n", err, firstOut)
		os.Exit(1)
	}

	sets, err := parseFirstOutput(firstOut, len(cases))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse first run output: %v\n", err)
		os.Exit(1)
	}

	data := make([]caseData, len(cases))
	for i, x := range cases {
		data[i] = caseData{x: x, set: sets[i]}
	}

	mutated := mutateSets(data, rng)
	secondInput, order := buildSecondInput(mutated, rng)

	secondOut, err := runCandidate(candidate, secondInput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed on second run: %v\n%s\n", err, secondOut)
		os.Exit(1)
	}

	ans, err := parseSecondOutput(secondOut, len(cases))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse second run output: %v\n", err)
		os.Exit(1)
	}

	for i, idx := range order {
		expected := data[idx].x
		got := ans[i]
		if got != expected {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d, got %d\n", i+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(cases))
}

func generateTests(rng *rand.Rand) []int {
	base := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		32768, 16384, 12345, 22222, 30000, 50, 9999, 2,
	}
	for len(base) < 400 {
		base = append(base, rng.Intn(maxValue)+1)
	}
	return base
}

func buildFirstInput(cases []int) string {
	var b strings.Builder
	fmt.Fprintln(&b, "first")
	fmt.Fprintln(&b, len(cases))
	for _, x := range cases {
		fmt.Fprintln(&b, x)
	}
	return b.String()
}

func parseFirstOutput(output string, t int) ([][]bool, error) {
	tokens := strings.Fields(output)
	sets := make([][]bool, t)
	pos := 0
	for i := 0; i < t; i++ {
		if pos >= len(tokens) {
			return nil, fmt.Errorf("missing n for test %d", i+1)
		}
		n64, err := strconv.ParseInt(tokens[pos], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("test %d: invalid n %q", i+1, tokens[pos])
		}
		n := int(n64)
		pos++
		if n < 0 || n > 20 {
			return nil, fmt.Errorf("test %d: n=%d out of range", i+1, n)
		}
		set := make([]bool, 21)
		for j := 0; j < n; j++ {
			if pos >= len(tokens) {
				return nil, fmt.Errorf("test %d: missing element %d", i+1, j+1)
			}
			val64, err := strconv.ParseInt(tokens[pos], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid element %q", i+1, tokens[pos])
			}
			val := int(val64)
			pos++
			if val < 1 || val > 20 {
				return nil, fmt.Errorf("test %d: element %d outside [1,20]", i+1, val)
			}
			if set[val] {
				return nil, fmt.Errorf("test %d: duplicate element %d", i+1, val)
			}
			set[val] = true
		}
		sets[i] = set
	}
	if pos != len(tokens) {
		if pos < len(tokens) {
			return nil, fmt.Errorf("extra output detected starting at token %q", tokens[pos])
		}
		return nil, fmt.Errorf("extra output detected")
	}
	return sets, nil
}

func mutateSets(data []caseData, rng *rand.Rand) [][]int {
	result := make([][]int, len(data))
	for i, tc := range data {
		result[i] = mutateSet(tc.set, rng)
	}
	return result
}

func mutateSet(original []bool, rng *rand.Rand) []int {
	set := make([]bool, len(original))
	copy(set, original)
	present := func() []int {
		var res []int
		for i := 1; i <= 20; i++ {
			if set[i] {
				res = append(res, i)
			}
		}
		return res
	}
	absent := func() []int {
		var res []int
		for i := 1; i <= 20; i++ {
			if !set[i] {
				res = append(res, i)
			}
		}
		return res
	}
	operation := rng.Intn(3)
	switch operation {
	case 1:
		pr := present()
		if len(pr) > 0 {
			idx := pr[rng.Intn(len(pr))]
			set[idx] = false
		} else {
			ab := absent()
			if len(ab) > 0 {
				idx := ab[rng.Intn(len(ab))]
				set[idx] = true
			}
		}
	case 2:
		ab := absent()
		if len(ab) > 0 {
			idx := ab[rng.Intn(len(ab))]
			set[idx] = true
		} else {
			pr := present()
			if len(pr) > 0 {
				idx := pr[rng.Intn(len(pr))]
				set[idx] = false
			}
		}
	}
	var res []int
	for i := 1; i <= 20; i++ {
		if set[i] {
			res = append(res, i)
		}
	}
	return res
}

func buildSecondInput(mutated [][]int, rng *rand.Rand) (string, []int) {
	n := len(mutated)
	order := rng.Perm(n)
	var b strings.Builder
	fmt.Fprintln(&b, "second")
	fmt.Fprintln(&b, n)
	for _, idx := range order {
		arr := mutated[idx]
		fmt.Fprintln(&b, len(arr))
		if len(arr) > 0 {
			for i, v := range arr {
				if i > 0 {
					fmt.Fprint(&b, " ")
				}
				fmt.Fprint(&b, v)
			}
		}
		fmt.Fprintln(&b)
	}
	return b.String(), order
}

func parseSecondOutput(output string, t int) ([]int, error) {
	tokens := strings.Fields(output)
	if len(tokens) < t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(tokens))
	}
	if len(tokens) > t {
		return nil, fmt.Errorf("extra output detected starting at token %q", tokens[t])
	}
	ans := make([]int, t)
	for i := 0; i < t; i++ {
		val64, err := strconv.ParseInt(tokens[i], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("answer %d is not an integer (%q)", i+1, tokens[i])
		}
		val := int(val64)
		if val < 1 || val > maxValue {
			return nil, fmt.Errorf("answer %d = %d is outside [1,%d]", i+1, val, maxValue)
		}
		ans[i] = val
	}
	return ans, nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}
