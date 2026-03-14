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
	"time"
)

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildOfficial() (string, error) {
	src := os.Getenv("REFERENCE_SOURCE_PATH")
	if src == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH env var not set")
	}
	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("1907G_official_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", tmp, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build official solution: %v %s", err, string(out))
	}
	return tmp, nil
}

// parseTestInput parses one test case from the generated input.
type testCase struct {
	n      int
	states []int // 0 or 1
	a      []int // 1-indexed targets
}

func generateTests(rng *rand.Rand, t int) (string, []testCase) {
	var input strings.Builder
	cases := make([]testCase, t)
	input.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(5) + 2
		cases[i].n = n
		cases[i].states = make([]int, n)
		input.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				input.WriteByte('0')
				cases[i].states[j] = 0
			} else {
				input.WriteByte('1')
				cases[i].states[j] = 1
			}
		}
		input.WriteString("\n")
		cases[i].a = make([]int, n)
		for j := 0; j < n; j++ {
			v := rng.Intn(n) + 1
			if v == j+1 {
				v = (v % n) + 1
			}
			cases[i].a[j] = v
			input.WriteString(fmt.Sprintf("%d ", v))
		}
		input.WriteString("\n")
	}
	return input.String(), cases
}

// parseOutput parses the output for t test cases.
// Each test case is either "-1" on a line, or a count k followed by k space-separated switch indices.
func parseOutput(raw string, t int) ([][]int, error) {
	tokens := strings.Fields(raw)
	idx := 0
	next := func() (string, error) {
		if idx >= len(tokens) {
			return "", fmt.Errorf("unexpected end of output")
		}
		s := tokens[idx]
		idx++
		return s, nil
	}
	results := make([][]int, t)
	for i := 0; i < t; i++ {
		tok, err := next()
		if err != nil {
			return nil, fmt.Errorf("test %d: %v", i+1, err)
		}
		k, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("test %d: bad count %q", i+1, tok)
		}
		if k == -1 {
			results[i] = nil
			continue
		}
		switches := make([]int, k)
		for j := 0; j < k; j++ {
			tok, err = next()
			if err != nil {
				return nil, fmt.Errorf("test %d: %v", i+1, err)
			}
			v, err := strconv.Atoi(tok)
			if err != nil {
				return nil, fmt.Errorf("test %d: bad switch %q", i+1, tok)
			}
			switches[j] = v
		}
		results[i] = switches
	}
	return results, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: verifierG /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	official, err := buildOfficial()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(official)

	rng := rand.New(rand.NewSource(6))
	t := 100
	inputStr, cases := generateTests(rng, t)

	refRaw, err := runBinary(official, inputStr)
	if err != nil {
		fmt.Printf("official solution failed: %v\n", err)
		os.Exit(1)
	}
	candRaw, err := runBinary(bin, inputStr)
	if err != nil {
		fmt.Printf("binary failed: %v\n", err)
		os.Exit(1)
	}

	refResults, err := parseOutput(refRaw, t)
	if err != nil {
		fmt.Printf("parsing official output: %v\n", err)
		os.Exit(1)
	}
	candResults, err := parseOutput(candRaw, t)
	if err != nil {
		fmt.Printf("parsing candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := 0; i < t; i++ {
		tc := cases[i]
		ref := refResults[i]
		cand := candResults[i]

		// Both should agree on impossibility
		if ref == nil && cand == nil {
			continue
		}
		if ref == nil && cand != nil {
			fmt.Printf("test %d: reference says -1 but candidate gave answer\n", i+1)
			os.Exit(1)
		}
		if ref != nil && cand == nil {
			fmt.Printf("test %d: candidate says -1 but reference found solution with %d switches\n", i+1, len(ref))
			os.Exit(1)
		}

		// Check optimality: candidate should use same number of switches
		if len(cand) != len(ref) {
			fmt.Printf("test %d: candidate used %d switches, reference used %d\n", i+1, len(cand), len(ref))
			os.Exit(1)
		}

		// Verify correctness: apply candidate switches and check all lights are off
		state := make([]int, tc.n)
		copy(state, tc.states)
		for _, sw := range cand {
			idx := sw - 1 // convert to 0-indexed
			if idx < 0 || idx >= tc.n {
				fmt.Printf("test %d: switch index %d out of range [1,%d]\n", i+1, sw, tc.n)
				os.Exit(1)
			}
			state[idx] ^= 1
			target := tc.a[idx] - 1 // a is 1-indexed
			state[target] ^= 1
		}
		for j := 0; j < tc.n; j++ {
			if state[j] != 0 {
				fmt.Printf("test %d: light %d still on after applying candidate switches\n", i+1, j+1)
				os.Exit(1)
			}
		}
	}

	fmt.Println("All 100 test cases passed.")
}
