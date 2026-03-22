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
	name  string
	input string
}

// Embedded solver for 1866I
func solveI(input string) string {
	bpos := 0
	buffer := []byte(input)

	readInt := func() int {
		for bpos < len(buffer) && buffer[bpos] <= ' ' {
			bpos++
		}
		if bpos >= len(buffer) {
			return 0
		}
		res := 0
		for bpos < len(buffer) && buffer[bpos] > ' ' {
			res = res*10 + int(buffer[bpos]-'0')
			bpos++
		}
		return res
	}

	N := readInt()
	if N == 0 {
		return ""
	}
	M := readInt()
	K := readInt()

	specials := make([][]int, N+1)
	for i := 0; i < K; i++ {
		r := readInt()
		c := readInt()
		specials[r] = append(specials[r], c)
	}

	in_S := make([]bool, M+1)
	for i := 1; i <= M; i++ {
		in_S[i] = true
	}

	max_S := M

	for r := N; r >= 1; r-- {
		s_1 := 0
		for _, c := range specials[r] {
			if c > s_1 {
				s_1 = c
			}
		}

		for max_S > 0 && !in_S[max_S] {
			max_S--
		}

		c := max_S
		if c > s_1 && c > 0 {
			in_S[c] = false
			if r == 1 && c == 1 {
				return "Bhinneka"
			}
		}

		for _, sp := range specials[r] {
			in_S[sp] = false
		}
	}

	return "Chaneka"
}

func parseWinner(out string) (string, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return "", fmt.Errorf("empty output")
	}
	w := fields[0]
	if w != "Chaneka" && w != "Bhinneka" {
		return "", fmt.Errorf("invalid winner %q", w)
	}
	return w, nil
}

func verifyCase(candidate string, tc testCase) error {
	expected := solveI(tc.input)
	expWinner, err := parseWinner(expected)
	if err != nil {
		return fmt.Errorf("embedded solver error: %v", err)
	}

	var cmd *exec.Cmd
	if strings.HasSuffix(candidate, ".go") {
		cmd = exec.Command("go", "run", candidate)
	} else {
		cmd = exec.Command(candidate)
	}
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("candidate error: %v\n%s", err, out.String())
	}

	got, err := parseWinner(strings.TrimSpace(out.String()))
	if err != nil {
		return fmt.Errorf("invalid candidate output: %v", err)
	}
	if got != expWinner {
		return fmt.Errorf("expected %s, got %s", expWinner, got)
	}
	return nil
}

func formatInput(n, m, k int, cells [][2]int) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for _, c := range cells {
		fmt.Fprintf(&sb, "%d %d\n", c[0], c[1])
	}
	return sb.String()
}

func manualTests() []testCase {
	return []testCase{
		{name: "immediate_row_win", input: formatInput(3, 3, 1, [][2]int{{1, 3}})},
		{name: "immediate_col_win", input: formatInput(4, 4, 1, [][2]int{{4, 1}})},
		{name: "no_special", input: formatInput(2, 2, 0, nil)},
		{name: "sample_like", input: formatInput(3, 3, 1, [][2]int{{2, 2}})},
	}
}

func randomTest(name string, rng *rand.Rand, maxN int) testCase {
	n := rng.Intn(maxN-1) + 1
	m := rng.Intn(maxN-1) + 1
	if n == 1 && m == 1 {
		n = 2
	}
	maxCells := n*m - 1
	if maxCells < 0 {
		maxCells = 0
	}
	k := rng.Intn(myMin(maxCells, 10) + 1)
	cells := make([][2]int, 0, k)
	used := make(map[[2]int]bool)
	for len(cells) < k {
		x := rng.Intn(n) + 1
		y := rng.Intn(m) + 1
		if x == 1 && y == 1 {
			continue
		}
		key := [2]int{x, y}
		if used[key] {
			continue
		}
		used[key] = true
		cells = append(cells, key)
	}
	return testCase{name: name, input: formatInput(n, m, k, cells)}
}

func myMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func generateTests() []testCase {
	tests := manualTests()
	seeds := []int64{1, 2, 3, 4, 5}
	for idx, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		tests = append(tests, randomTest(fmt.Sprintf("deterministic_%d", idx+1), rng, 10))
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 80 {
		tests = append(tests, randomTest(fmt.Sprintf("random_%d", len(tests)+1), rng, 200))
	}
	return tests
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	tests := generateTests()
	for i, tc := range tests {
		if err := verifyCase(candidate, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d (%s) failed: %v\ninput:\n%s", i+1, tc.name, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
