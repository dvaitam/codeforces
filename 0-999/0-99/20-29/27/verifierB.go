package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type match struct {
	winner int
	loser  int
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parsePair(out string) (int, int, error) {
	fields := strings.Fields(out)
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("expected two integers, got %q", out)
	}
	first, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	second, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid integer %q: %v", fields[1], err)
	}
	return first, second, nil
}

func parseInputData(input string) (int, [][]bool, [2]int, error) {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return 0, nil, [2]int{}, fmt.Errorf("failed to read n: %v", err)
	}
	if n < 3 {
		return 0, nil, [2]int{}, fmt.Errorf("invalid n %d", n)
	}
	total := n*(n-1)/2 - 1
	wins := make([][]bool, n+1)
	for i := range wins {
		wins[i] = make([]bool, n+1)
	}
	for i := 0; i < total; i++ {
		var w, l int
		if _, err := fmt.Fscan(in, &w, &l); err != nil {
			return 0, nil, [2]int{}, fmt.Errorf("failed to read match %d: %v", i+1, err)
		}
		if w < 1 || w > n || l < 1 || l > n || w == l {
			return 0, nil, [2]int{}, fmt.Errorf("invalid match %d: %d %d", i+1, w, l)
		}
		if wins[w][l] {
			return 0, nil, [2]int{}, fmt.Errorf("duplicate result for %d beating %d", w, l)
		}
		if wins[l][w] {
			return 0, nil, [2]int{}, fmt.Errorf("conflicting results for pair %d %d", w, l)
		}
		wins[w][l] = true
	}
	missing := [2]int{}
	missingCnt := 0
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if !wins[i][j] && !wins[j][i] {
				missing = [2]int{i, j}
				missingCnt++
			}
		}
	}
	if missingCnt != 1 {
		return 0, nil, [2]int{}, fmt.Errorf("expected one missing pair, got %d", missingCnt)
	}
	return n, wins, missing, nil
}

// validateCandidate checks that declaring x beat y completes the tournament
// into a valid total order (no cycles, no missing/duplicate pairs).
func validateCandidate(n int, wins [][]bool, x, y int) error {
	if x < 1 || x > n || y < 1 || y > n {
		return fmt.Errorf("players must be between 1 and %d, got %d %d", n, x, y)
	}
	if x == y {
		return fmt.Errorf("players must be distinct, got %d %d", x, y)
	}
	if wins[x][y] || wins[y][x] {
		return fmt.Errorf("match between %d and %d is already recorded", x, y)
	}

	complete := make([][]bool, n+1)
	for i := range wins {
		complete[i] = make([]bool, n+1)
		copy(complete[i], wins[i])
	}
	complete[x][y] = true

	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if complete[i][j] == complete[j][i] {
				if complete[i][j] {
					return fmt.Errorf("conflicting double records for %d and %d", i, j)
				}
				return fmt.Errorf("match between %d and %d is still missing", i, j)
			}
		}
	}

	adj := make([][]int, n+1)
	indeg := make([]int, n+1)
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			if wins[i][j] {
				adj[i] = append(adj[i], j)
				indeg[j]++
			}
		}
	}
	adj[x] = append(adj[x], y)
	indeg[y]++

	queue := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if indeg[i] == 0 {
			queue = append(queue, i)
		}
	}
	processed := 0
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		processed++
		for _, to := range adj[v] {
			indeg[to]--
			if indeg[to] == 0 {
				queue = append(queue, to)
			}
		}
	}
	if processed != n {
		return fmt.Errorf("declared result %d beats %d creates a cycle", x, y)
	}
	return nil
}

func makeCase(n int, ranking []int, missing [2]int) string {
	matches := make([]match, 0, n*(n-1)/2-1)
	a, b := missing[0], missing[1]
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if (i == a && j == b) || (i == b && j == a) {
				continue
			}
			if ranking[i-1] < ranking[j-1] {
				matches = append(matches, match{i, j})
			} else {
				matches = append(matches, match{j, i})
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, m := range matches {
		fmt.Fprintf(&sb, "%d %d\n", m.winner, m.loser)
	}
	return sb.String()
}

func deterministicCases() []string {
	cases := []string{
		makeCase(3, []int{0, 1, 2}, [2]int{1, 3}),
		makeCase(4, []int{0, 3, 1, 2}, [2]int{2, 4}),
		makeCase(5, []int{2, 0, 4, 1, 3}, [2]int{4, 5}),
		makeCase(6, []int{5, 1, 3, 0, 2, 4}, [2]int{1, 2}),
	}

	n := 50
	ranking := make([]int, n)
	for i := range ranking {
		ranking[i] = i
	}
	input := makeCase(n, ranking, [2]int{1, n})
	lines := strings.Split(strings.TrimSpace(input), "\n")
	matches := lines[1:]
	rng := rand.New(rand.NewSource(123456789))
	rng.Shuffle(len(matches), func(i, j int) {
		matches[i], matches[j] = matches[j], matches[i]
	})
	var sb strings.Builder
	sb.WriteString(lines[0])
	sb.WriteByte('\n')
	for _, line := range matches {
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	cases = append(cases, sb.String())
	return cases
}

func generateRandomCase(rng *rand.Rand) string {
	n := rng.Intn(48) + 3 // 3..50
	perm := rng.Perm(n)
	pos := make([]int, n)
	for idx, player := range perm {
		pos[player] = idx
	}
	allMatches := make([]match, 0, n*(n-1)/2)
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if pos[i-1] < pos[j-1] {
				allMatches = append(allMatches, match{i, j})
			} else {
				allMatches = append(allMatches, match{j, i})
			}
		}
	}
	missIdx := rng.Intn(len(allMatches))
	allMatches = append(allMatches[:missIdx], allMatches[missIdx+1:]...)
	rng.Shuffle(len(allMatches), func(i, j int) {
		allMatches[i], allMatches[j] = allMatches[j], allMatches[i]
	})

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, m := range allMatches {
		fmt.Fprintf(&sb, "%d %d\n", m.winner, m.loser)
	}
	return sb.String()
}

func verifyCase(candidate, input string) error {
	n, wins, missing, err := parseInputData(input)
	if err != nil {
		return fmt.Errorf("bad test data: %v", err)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		return err
	}
	cx, cy, err := parsePair(candOut)
	if err != nil {
		return err
	}

	// Check the candidate identified the correct missing pair (either orientation).
	if !((cx == missing[0] && cy == missing[1]) || (cx == missing[1] && cy == missing[0])) {
		return fmt.Errorf("missing pair is {%d,%d}, but candidate returned %d %d",
			missing[0], missing[1], cx, cy)
	}

	// Check that the declared winner/loser is consistent (no cycle).
	if err := validateCandidate(n, wins, cx, cy); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := deterministicCases()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 300; i++ {
		tests = append(tests, generateRandomCase(rng))
	}

	for idx, input := range tests {
		if err := verifyCase(candidate, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
