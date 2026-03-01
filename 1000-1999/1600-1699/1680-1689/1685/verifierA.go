package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
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

func generateTests() []string {
	r := rand.New(rand.NewSource(1))
	tests := make([]string, 0, 100)
	for len(tests) < 100 {
		n := r.Intn(8) + 3 // 3..10
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(r.Intn(20)))
		}
		sb.WriteByte('\n')
		tests = append(tests, sb.String())
	}
	return tests
}

func isZigzagCircle(b []int) bool {
	n := len(b)
	for i := 0; i < n; i++ {
		prev := b[(i-1+n)%n]
		next := b[(i+1)%n]
		if !((b[i] > prev && b[i] > next) || (b[i] < prev && b[i] < next)) {
			return false
		}
	}
	return true
}

func existsArrangement(a []int) bool {
	n := len(a)
	if n%2 == 1 {
		return false
	}
	b := make([]int, n)
	used := make([]bool, n)
	sorted := append([]int(nil), a...)
	sort.Ints(sorted)

	var dfs func(pos int) bool
	dfs = func(pos int) bool {
		if pos == n {
			return isZigzagCircle(b)
		}
		for i := 0; i < n; i++ {
			if used[i] {
				continue
			}
			if i > 0 && sorted[i] == sorted[i-1] && !used[i-1] {
				continue
			}
			b[pos] = sorted[i]
			if pos >= 2 {
				x, y, z := b[pos-2], b[pos-1], b[pos]
				if !((y > x && y > z) || (y < x && y < z)) {
					continue
				}
			}
			used[i] = true
			if dfs(pos + 1) {
				return true
			}
			used[i] = false
		}
		return false
	}

	return dfs(0)
}

func parseCase(input string) ([]int, error) {
	fields := strings.Fields(input)
	if len(fields) < 3 {
		return nil, fmt.Errorf("bad generated input")
	}
	idx := 0
	t, err := strconv.Atoi(fields[idx])
	if err != nil || t != 1 {
		return nil, fmt.Errorf("bad t in generated input")
	}
	idx++
	n, err := strconv.Atoi(fields[idx])
	if err != nil {
		return nil, fmt.Errorf("bad n in generated input")
	}
	idx++
	if len(fields) < idx+n {
		return nil, fmt.Errorf("bad array in generated input")
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i], err = strconv.Atoi(fields[idx+i])
		if err != nil {
			return nil, fmt.Errorf("bad number in generated input")
		}
	}
	return a, nil
}

func validateOutput(input, output string) error {
	a, err := parseCase(input)
	if err != nil {
		return err
	}
	possible := existsArrangement(a)

	tokens := strings.Fields(output)
	if len(tokens) == 0 {
		return fmt.Errorf("empty output")
	}
	ans := strings.ToUpper(tokens[0])
	if ans != "YES" && ans != "NO" {
		return fmt.Errorf("first token must be YES or NO")
	}

	if ans == "NO" {
		if possible {
			return fmt.Errorf("reported NO, but valid arrangement exists")
		}
		return nil
	}

	if !possible {
		return fmt.Errorf("reported YES, but no valid arrangement exists")
	}

	n := len(a)
	if len(tokens) < 1+n {
		return fmt.Errorf("YES output must contain %d numbers", n)
	}
	b := make([]int, n)
	for i := 0; i < n; i++ {
		b[i], err = strconv.Atoi(tokens[1+i])
		if err != nil {
			return fmt.Errorf("invalid integer in arrangement")
		}
	}

	need := make(map[int]int)
	for _, v := range a {
		need[v]++
	}
	for _, v := range b {
		need[v]--
	}
	for v, cnt := range need {
		if cnt != 0 {
			return fmt.Errorf("arrangement is not a permutation (value %d mismatch)", v)
		}
	}

	if !isZigzagCircle(b) {
		return fmt.Errorf("arrangement does not satisfy zigzag circle constraints")
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := generateTests()
	for i, input := range tests {
		candOut, cErr := runBinary(candidate, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "case %d: candidate error: %v\n", i+1, cErr)
			os.Exit(1)
		}
		if err := validateOutput(input, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%serror:%v\nactual:%s", i+1, input, err, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
