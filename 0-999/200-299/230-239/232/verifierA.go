package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt.
const testcasesAData = `
906416936
36806112
693395237
85777400
93893823
12121379
302818751
81135646
980872303
686621452
334813583
743689233
973874292
370855069
240834559
512862115
188372850
99430360
643165047
90599696
267838844
45460550
883152694
417049584
407071030
935053888
659039390
243275550
536741515
64852770
985928065
652486957
368731254
246450223
915432900
311640501
332565483
755762563
411550743
190754073
420328525
632759534
832857026
996046336
992107447
77093850
487422242
330073672
711323017
335764258
485827097
805845141
88761208
180702919
35768680
159582351
946102762
799444987
972406321
967626496
572419280
538994442
909851832
936897252
328735085
8591695
142039843
989238713
871823974
395759504
169714263
983315345
909285057
819240908
996679379
133955526
696292674
896141622
590259313
876399194
13148367
765119329
393861987
550243994
725498774
378207075
782678593
978262015
814718349
470610814
313113177
961777074
700996511
161096102
38293397
30247196
424546349
`

// generateGraph implements the logic from 232A.go, returning the constructed graph.
func generateGraph(k int) (int, [][]bool) {
	const maxN = 101
	adj := make([][]bool, maxN)
	for i := range adj {
		adj[i] = make([]bool, maxN)
	}
	n := 1
	for k > 0 {
		m := n
		for m*(m-1)/2 > k {
			m--
		}
		for i := 0; i < m; i++ {
			adj[i][n] = true
			adj[n][i] = true
		}
		k -= m * (m - 1) / 2
		n++
	}
	graph := make([][]bool, n)
	for i := 0; i < n; i++ {
		graph[i] = make([]bool, n)
		copy(graph[i], adj[i][:n])
	}
	return n, graph
}

func countTriangles(adj [][]bool) int {
	n := len(adj)
	cnt := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if !adj[i][j] {
				continue
			}
			for k := j + 1; k < n; k++ {
				if adj[i][k] && adj[j][k] {
					cnt++
				}
			}
		}
	}
	return cnt
}

func runCase(bin string, k int) error {
	input := fmt.Sprintf("%d\n", k)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(out.String())))
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	n64, err := strconv.ParseInt(strings.TrimSpace(scanner.Text()), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse n: %v", err)
	}
	n := int(n64)
	if n < 3 || n > 100 {
		return fmt.Errorf("invalid n %d", n)
	}
	adj := make([][]bool, n)
	for i := 0; i < n; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("not enough lines for adjacency matrix")
		}
		line := strings.TrimSpace(scanner.Text())
		if len(line) != n {
			return fmt.Errorf("line %d length mismatch", i+1)
		}
		adj[i] = make([]bool, n)
		for j := 0; j < n; j++ {
			c := line[j]
			if c != '0' && c != '1' {
				return fmt.Errorf("invalid char at row %d col %d", i+1, j+1)
			}
			if i == j && c != '0' {
				return fmt.Errorf("self loop at row %d", i+1)
			}
			adj[i][j] = c == '1'
		}
	}
	if scanner.Scan() {
		extra := strings.TrimSpace(scanner.Text())
		if extra != "" {
			return fmt.Errorf("extra output: %q", extra)
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if adj[i][j] != adj[j][i] {
				return fmt.Errorf("asymmetry at %d,%d", i, j)
			}
		}
	}
	if countTriangles(adj) != k {
		return fmt.Errorf("expected %d triangles", k)
	}
	nRef, _ := generateGraph(k)
	if n > nRef {
		return fmt.Errorf("n (%d) larger than reference construction (%d)", n, nRef)
	}
	return nil
}

func parseTestcases() ([]int, error) {
	sc := bufio.NewScanner(strings.NewReader(testcasesAData))
	var cases []int
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		val, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("parse testcase: %w", err)
		}
		cases = append(cases, val)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, k := range testcases {
		if err := runCase(bin, k); err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
