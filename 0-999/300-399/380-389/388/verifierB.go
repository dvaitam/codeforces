package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Testcases embedded from testcasesB.txt (one k per line).
const rawTestcases = `990085
584054
13706
5143
389994
299330
872337
872618
225038
57062
978309
326955
688419
711381
155253
133621
70643
752812
163996
618060
438342
281471
802465
849724
137088
76496
944571
202729
967038
801531
168470
631275
436994
227115
819471
411650
567354
889506
514036
181528
908447
187494
199513
921436
856910
691451
167533
910938
495302
65880
524576
195396
265880
38917
407593
754130
516235
781849
63710
90187
248336
499516
15126
28723
659648
995663
51914
132955
99248
28404
129089
253826
931925
957724
131256
297559
396342
995049
895486
301274
950126
135342
926885
3014
365784
190212
404578
955689
628816
218830
247258
745819
349290
384048
38786
536168
750302
467880
8803
97086`

const testcaseCount = 100

func loadTestcases() ([]int, error) {
	lines := strings.Split(strings.TrimSpace(rawTestcases), "\n")
	if len(lines) != testcaseCount {
		return nil, fmt.Errorf("expected %d testcases got %d", testcaseCount, len(lines))
	}
	res := make([]int, 0, len(lines))
	for i, l := range lines {
		v, err := strconv.Atoi(strings.TrimSpace(l))
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", i+1, err)
		}
		res = append(res, v)
	}
	return res, nil
}

// countShortestPaths parses the output graph and counts shortest paths from vertex 1 to vertex 2
func countShortestPaths(output string) (int64, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 1 {
		return 0, fmt.Errorf("empty output")
	}
	n, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return 0, fmt.Errorf("invalid n: %v", err)
	}
	if n < 2 || n > 1000 {
		return 0, fmt.Errorf("n=%d out of range [2,1000]", n)
	}
	if len(lines) < n+1 {
		return 0, fmt.Errorf("expected %d rows of adjacency matrix, got %d", n, len(lines)-1)
	}
	// Parse adjacency matrix
	adj := make([][]bool, n)
	for i := 0; i < n; i++ {
		row := strings.TrimSpace(lines[i+1])
		if len(row) != n {
			return 0, fmt.Errorf("row %d has length %d, expected %d", i+1, len(row), n)
		}
		adj[i] = make([]bool, n)
		for j := 0; j < n; j++ {
			if row[j] == 'Y' {
				adj[i][j] = true
			} else if row[j] != 'N' {
				return 0, fmt.Errorf("invalid character '%c' at row %d col %d", row[j], i+1, j+1)
			}
		}
	}
	// Check simple undirected: no self-loops, symmetric
	for i := 0; i < n; i++ {
		if adj[i][i] {
			return 0, fmt.Errorf("self-loop at vertex %d", i+1)
		}
		for j := 0; j < n; j++ {
			if adj[i][j] != adj[j][i] {
				return 0, fmt.Errorf("asymmetric edge between %d and %d", i+1, j+1)
			}
		}
	}
	// BFS from vertex 0 (vertex 1 in 1-indexed)
	dist := make([]int, n)
	ways := make([]int64, n)
	for i := range dist {
		dist[i] = -1
	}
	dist[0] = 0
	ways[0] = 1
	queue := []int{0}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for u := 0; u < n; u++ {
			if !adj[v][u] {
				continue
			}
			if dist[u] == -1 {
				dist[u] = dist[v] + 1
				ways[u] = ways[v]
				queue = append(queue, u)
			} else if dist[u] == dist[v]+1 {
				ways[u] += ways[v]
			}
		}
	}
	if dist[1] == -1 {
		return 0, fmt.Errorf("no path from vertex 1 to vertex 2")
	}
	return ways[1], nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v, output: %s", err, out.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases, err := loadTestcases()
	if err != nil {
		fmt.Printf("failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, k := range testcases {
		input := fmt.Sprintf("%d\n", k)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d (k=%d) failed: %v\n", idx+1, k, err)
			os.Exit(1)
		}
		count, err := countShortestPaths(got)
		if err != nil {
			fmt.Printf("case %d (k=%d) invalid output: %v\n", idx+1, k, err)
			os.Exit(1)
		}
		if count != int64(k) {
			fmt.Printf("case %d (k=%d) failed: graph has %d shortest paths, expected %d\n", idx+1, k, count, k)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
