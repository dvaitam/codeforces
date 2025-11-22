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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func parseInput(input string) (int, [][]int, error) {
	in := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return 0, nil, fmt.Errorf("failed to read t: %w", err)
	}
	tests := make([][]int, t)
	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return 0, nil, fmt.Errorf("failed to read n for test %d: %w", i+1, err)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(in, &arr[j]); err != nil {
				return 0, nil, fmt.Errorf("failed to read value %d for test %d: %w", j+1, i+1, err)
			}
		}
		tests[i] = arr
	}
	return t, tests, nil
}

func parseOutput(output string, t int, ns []int) ([][]int, error) {
	fields := strings.Fields(output)
	results := make([][]int, t)
	pos := 0
	for i := 0; i < t; i++ {
		need := ns[i]
		if pos+need > len(fields) {
			return nil, fmt.Errorf("test %d: expected %d numbers, got %d", i+1, need, len(fields)-pos)
		}
		test := make([]int, need)
		for j := 0; j < need; j++ {
			if _, err := fmt.Sscan(fields[pos+j], &test[j]); err != nil {
				return nil, fmt.Errorf("test %d value %d: %w", i+1, j+1, err)
			}
		}
		pos += need
		results[i] = test
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("unexpected extra output: %d leftover values", len(fields)-pos)
	}
	return results, nil
}

func checkPermutation(input, output []int) bool {
	if len(input) != len(output) {
		return false
	}
	freq := make(map[int]int)
	for _, v := range input {
		freq[v]++
	}
	for _, v := range output {
		freq[v]--
		if freq[v] < 0 {
			return false
		}
	}
	for _, c := range freq {
		if c != 0 {
			return false
		}
	}
	return true
}

func checkGreedy(arr []int) bool {
	used := make([]bool, len(arr))
	current := 0
	for step := 0; step < len(arr); step++ {
		best := -1
		for i, v := range arr {
			if !used[i] {
				g := gcd(current, v)
				if g > best {
					best = g
				}
			}
		}
		chosen := arr[step]
		gChosen := gcd(current, chosen)
		if gChosen != best {
			return false
		}
		// mark chosen index in original array occurrence
		idx := -1
		for i := range arr {
			if !used[i] && arr[i] == chosen {
				idx = i
				break
			}
		}
		if idx == -1 {
			return false
		}
		used[idx] = true
		current = gChosen
	}
	return true
}

func validateOutput(input string, output string) error {
	t, tests, err := parseInput(input)
	if err != nil {
		return err
	}
	ns := make([]int, t)
	for i, arr := range tests {
		ns[i] = len(arr)
	}
	results, err := parseOutput(output, t, ns)
	if err != nil {
		return err
	}
	for i := 0; i < t; i++ {
		if !checkPermutation(tests[i], results[i]) {
			return fmt.Errorf("test %d: output is not a permutation of input", i+1)
		}
		if !checkGreedy(results[i]) {
			return fmt.Errorf("test %d: output does not follow greedy gcd rule", i+1)
		}
	}
	return nil
}

func runCase(bin, input string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, errb.String())
	}
	if err := validateOutput(strings.TrimSpace(input), strings.TrimSpace(out.String())); err != nil {
		return err
	}
	return nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		v := rng.Intn(1000) + 1
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		if err := runCase(bin, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
