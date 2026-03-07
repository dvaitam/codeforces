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

// ---------- brute-force checker ----------

var memo map[string]bool

func stateKey(boxes [][]int) string {
	var sb strings.Builder
	for bi, box := range boxes {
		if bi > 0 {
			sb.WriteByte('|')
		}
		for pi, v := range box {
			if pi > 0 {
				sb.WriteByte(',')
			}
			sb.WriteByte(byte('0' + v))
		}
	}
	return sb.String()
}

// Returns true if the current-player-to-move wins.
func wins(boxes [][]int) bool {
	// skip leading empty boxes
	for len(boxes) > 0 {
		empty := true
		for _, v := range boxes[0] {
			if v > 0 {
				empty = false
				break
			}
		}
		if !empty {
			break
		}
		boxes = boxes[1:]
	}
	if len(boxes) == 0 {
		return false // no moves
	}
	key := stateKey(boxes)
	if v, ok := memo[key]; ok {
		return v
	}
	result := false
outer:
	for pi, v := range boxes[0] {
		if v == 0 {
			continue
		}
		for remove := 1; remove <= v; remove++ {
			boxes[0][pi] -= remove
			if !wins(boxes) {
				boxes[0][pi] += remove
				result = true
				break outer
			}
			boxes[0][pi] += remove
		}
	}
	memo[key] = result
	return result
}

func copyBoxes(boxes [][]int) [][]int {
	out := make([][]int, len(boxes))
	for i, b := range boxes {
		out[i] = make([]int, len(b))
		copy(out[i], b)
	}
	return out
}

// countWinning counts Alice-wins partitions for array a using brute force.
func countWinning(a []int) int {
	n := len(a)
	count := 0
	for mask := 0; mask < (1 << (n - 1)); mask++ {
		var boxes [][]int
		cur := []int{a[0]}
		for i := 1; i < n; i++ {
			if mask&(1<<(i-1)) != 0 {
				boxes = append(boxes, cur)
				cur = []int{a[i]}
			} else {
				cur = append(cur, a[i])
			}
		}
		boxes = append(boxes, cur)
		memo = map[string]bool{}
		if wins(copyBoxes(boxes)) {
			count++
		}
	}
	return count % 998244353
}

// ---------- candidate runner ----------

func runProgram(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseAnswers(out string) []string {
	reader := strings.NewReader(out)
	var res []string
	for {
		var token string
		if _, err := fmt.Fscan(reader, &token); err != nil {
			break
		}
		res = append(res, token)
	}
	return res
}

// ---------- test generation ----------

func genRandom(rng *rand.Rand) ([]byte, [][]int) {
	t := rng.Intn(4) + 1
	var cases [][]int
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(5) + 1 // 1..5 for fast brute force
		a := make([]int, n)
		for j := range a {
			a[j] = rng.Intn(3) + 1 // 1..3 for fast brute force
		}
		cases = append(cases, a)
		fmt.Fprintf(&sb, "%d\n", n)
		for j, v := range a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
	}
	return []byte(sb.String()), cases
}

// ---------- main ----------

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	// Fixed sample tests with hardcoded expected answers from the problem statement.
	type fixedCase struct {
		input    string
		expected []string
	}
	fixedTests := []fixedCase{
		{
			input:    "5\n3\n1 2 3\n4\n1 2 3 1\n5\n1 1 1 1 1\n2\n1 1\n10\n1 2 3 4 5 6 7 8 9 10\n",
			expected: []string{"1", "4", "16", "0", "205"},
		},
	}
	for i, ft := range fixedTests {
		candOut, err := runProgram(candidate, []byte(ft.input))
		if err != nil {
			fail("fixed test %d: candidate execution failed: %v", i+1, err)
		}
		candAns := parseAnswers(candOut)
		if len(candAns) != len(ft.expected) {
			fail("fixed test %d: expected %d answers got %d", i+1, len(ft.expected), len(candAns))
		}
		for j := range ft.expected {
			if candAns[j] != ft.expected[j] {
				fail("fixed test %d case %d: expected %s got %s", i+1, j+1, ft.expected[j], candAns[j])
			}
		}
	}

	// Random tests: compare candidate against brute-force.
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for iter := 0; iter < 200; iter++ {
		input, cases := genRandom(rng)

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fail("random test %d: candidate execution failed: %v\ninput:\n%s", iter+1, err, string(input))
		}
		candAns := parseAnswers(candOut)
		if len(candAns) != len(cases) {
			fail("random test %d: expected %d answers got %d\ninput:\n%s", iter+1, len(cases), len(candAns), string(input))
		}
		for j, a := range cases {
			want := countWinning(a)
			wantStr := fmt.Sprintf("%d", want)
			if candAns[j] != wantStr {
				fail("random test %d case %d: expected %s got %s\ninput:\n%s", iter+1, j+1, wantStr, candAns[j], string(input))
			}
		}
	}

	fmt.Println("OK")
}
