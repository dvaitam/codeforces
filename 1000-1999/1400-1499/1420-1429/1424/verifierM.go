package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCase struct {
	input      string
	used       []bool
	usedCount  int
	edges      [][2]int
	impossible bool
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierM.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for idx, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := checkOutput(strings.TrimSpace(out), tc); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func checkOutput(out string, tc testCase) error {
	if tc.impossible {
		if out != "IMPOSSIBLE" {
			return fmt.Errorf("expected IMPOSSIBLE but got %q", out)
		}
		return nil
	}
	if out == "IMPOSSIBLE" {
		return fmt.Errorf("expected alphabet ordering, got IMPOSSIBLE")
	}
	if strings.ContainsAny(out, " \t\r\n") {
		return fmt.Errorf("alphabet must be a single contiguous string without whitespace")
	}
	if len(out) != tc.usedCount {
		return fmt.Errorf("alphabet length mismatch: expected %d letters, got %d", tc.usedCount, len(out))
	}
	pos := make([]int, 26)
	for i := range pos {
		pos[i] = -1
	}
	for i, ch := range out {
		if ch < 'a' || ch > 'z' {
			return fmt.Errorf("invalid character %q in alphabet", ch)
		}
		idx := int(ch - 'a')
		if !tc.used[idx] {
			return fmt.Errorf("letter %c not present in input words", ch)
		}
		if pos[idx] != -1 {
			return fmt.Errorf("letter %c appears multiple times", ch)
		}
		pos[idx] = i
	}
	for idx, used := range tc.used {
		if used && pos[idx] == -1 {
			return fmt.Errorf("letter %c missing from alphabet", 'a'+idx)
		}
	}
	for _, e := range tc.edges {
		if pos[e[0]] >= pos[e[1]] {
			return fmt.Errorf("ordering constraint violated: %c must appear before %c", 'a'+e[0], 'a'+e[1])
		}
	}
	return nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeTest(1, 1, [][]string{{"a"}}),
		makeTest(2, 2, [][]string{{"ab", "ac"}, {"ad", "ae"}}),
		makeTest(2, 2, [][]string{{"ab", "a"}, {"b", "ba"}}),
	}
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest())
	}
	return tests
}

func makeTest(n, k int, words [][]string) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d\n", i+1)
		for j := 0; j < k; j++ {
			fmt.Fprintf(&sb, "%s\n", words[i][j])
		}
	}
	return buildCase(sb.String())
}

func randomTest() testCase {
	n := rand.Intn(4) + 1
	k := rand.Intn(3) + 1
	type page struct {
		p     int
		words []string
	}
	pages := make([]page, n)
	usedP := make(map[int]bool)
	for i := 0; i < n; i++ {
		p := rand.Intn(1000)
		for usedP[p] {
			p = rand.Intn(1000)
		}
		usedP[p] = true
		pages[i].p = p
		pages[i].words = make([]string, k)
		for j := 0; j < k; j++ {
			pages[i].words[j] = randomWord()
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for _, pg := range pages {
		fmt.Fprintf(&sb, "%d\n", pg.p)
		for _, w := range pg.words {
			fmt.Fprintf(&sb, "%s\n", w)
		}
	}
	return buildCase(sb.String())
}

func randomWord() string {
	length := rand.Intn(5) + 1
	alphabet := "abcdefghijklmnopqrstuvwxyz"[:rand.Intn(6)+1]
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(alphabet[rand.Intn(len(alphabet))])
	}
	return sb.String()
}

func buildCase(input string) testCase {
	used, usedCount, edges, impossible := analyzeInput(input)
	return testCase{
		input:      input,
		used:       used,
		usedCount:  usedCount,
		edges:      edges,
		impossible: impossible,
	}
}

func analyzeInput(input string) ([]bool, int, [][2]int, bool) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, k int
	fmt.Fscan(reader, &n, &k)
	type page struct {
		p     int
		words []string
	}
	pages := make([]page, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &pages[i].p)
		pages[i].words = make([]string, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(reader, &pages[i].words[j])
		}
	}
	sort.Slice(pages, func(i, j int) bool { return pages[i].p < pages[j].p })
	used := make([]bool, 26)
	var list []string
	for _, pg := range pages {
		for _, w := range pg.words {
			list = append(list, w)
			for _, ch := range w {
				if ch >= 'a' && ch <= 'z' {
					used[ch-'a'] = true
				}
			}
		}
	}
	edges, prefixConflict := buildEdges(list)
	cycle := false
	if !prefixConflict {
		cycle = hasCycle(used, edges)
	}
	usedCount := 0
	for _, v := range used {
		if v {
			usedCount++
		}
	}
	return append([]bool(nil), used...), usedCount, edges, prefixConflict || cycle
}

func buildEdges(words []string) ([][2]int, bool) {
	var edges [][2]int
	hasEdge := [26][26]bool{}
	prefixConflict := false
	for i := 1; i < len(words); i++ {
		prev := words[i-1]
		curr := words[i]
		minLen := len(prev)
		if len(curr) < minLen {
			minLen = len(curr)
		}
		found := false
		for j := 0; j < minLen; j++ {
			if prev[j] != curr[j] {
				u := int(prev[j] - 'a')
				v := int(curr[j] - 'a')
				if u >= 0 && u < 26 && v >= 0 && v < 26 && !hasEdge[u][v] {
					hasEdge[u][v] = true
					edges = append(edges, [2]int{u, v})
				}
				found = true
				break
			}
		}
		if !found && len(prev) > len(curr) {
			prefixConflict = true
		}
	}
	return edges, prefixConflict
}

func hasCycle(used []bool, edges [][2]int) bool {
	adj := make([][]int, 26)
	for _, e := range edges {
		adj[e[0]] = append(adj[e[0]], e[1])
	}
	color := make([]int, 26)
	var dfs func(int) bool
	dfs = func(u int) bool {
		color[u] = 1
		for _, v := range adj[u] {
			if color[v] == 0 {
				if dfs(v) {
					return true
				}
			} else if color[v] == 1 {
				return true
			}
		}
		color[u] = 2
		return false
	}
	for i, u := range used {
		if u && color[i] == 0 {
			if dfs(i) {
				return true
			}
		}
	}
	return false
}
