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

const V = 4
const INF int64 = 4e18

type edge struct {
	c1 int
	v  int64
	c2 int
}

type testcase struct {
	n     int
	edges []edge
}

const testcasesRaw = `100
2
3 93 1
4 20 4
1
1 52 1
5
3 29 1
3 100 3
2 34 1
2 83 1
3 25 3
2
3 81 3
3 78 1
3
4 23 2
2 36 4
1 1 3
3
3 53 2
4 56 3
4 30 2
3
3 11 1
1 81 4
3 90 4
3
2 9 2
4 82 2
4 24 3
3
4 82 3
2 13 3
1 36 2
5
2 43 1
2 59 3
1 46 1
1 95 3
3 42 1
3
3 100 2
4 38 1
2 38 4
2
3 77 4
2 74 3
1
3 59 1
2
3 38 3
1 27 4
4
2 8 1
1 95 1
2 78 2
1 75 4
2
3 16 1
3 84 4
2
4 31 2
4 63 4
1
2 57 4
2
4 64 2
2 5 1
3
3 68 2
2 54 2
3 42 2
1
3 73 1
4
1 50 4
1 27 4
2 38 3
4 54 3
5
2 44 3
4 10 4
3 6 2
4 99 2
3 22 1
6
4 96 4
4 28 4
1 21 2
1 15 3
4 29 4
1 21 2
6
3 68 4
4 11 1
1 63 1
3 6 2
3 99 1
1 45 3
1
1 49 4
2
3 30 4
4 13 4
1
1 66 3
4
4 9 4
2 62 3
4 72 1
2 21 3
2
2 64 3
3 70 3
1
2 83 1
3
1 63 1
4 67 1
2 38 4
3
2 81 2
1 79 1
3 73 4
3
4 51 4
2 99 3
3 18 3
4
1 87 2
2 48 3
2 87 3
1 52 1
6
2 84 3
3 23 3
3 78 1
1 46 1
1 24 2
4 98 1
1
2 87 4
6
2 89 3
4 63 2
2 48 3
2 69 3
4 65 4
4 37 3
4
4 33 1
2 89 4
3 50 4
1 66 2
1
4 78 3
3
1 62 3
2 78 4
1 31 2
1
3 7 2
2
4 3 3
2 97 2
1
2 41 1
1
1 8 3
6
3 94 2
2 17 4
1 89 3
1 65 2
4 28 2
3 65 4
1
4 21 2
3
4 52 3
3 50 2
1 4 4
3
1 20 3
1 15 2
1 30 2
5
1 92 4
2 94 4
4 23 3
2 61 1
3 83 3
4
1 33 2
3 81 4
3 31 1
3 85 1
6
1 68 3
4 11 3
2 64 1
2 26 1
1 78 3
2 15 3
5
1 29 1
3 94 1
1 62 1
2 47 2
2 39 3
4
4 23 1
4 9 4
3 27 3
4 23 3
2
1 56 4
3 97 3
4
4 97 3
2 36 3
1 52 1
3 87 3
6
2 12 3
2 90 1
2 86 3
2 98 4
1 2 3
4 22 2
1
2 15 1
3
3 96 2
1 92 2
1 87 2
2
3 76 1
3 95 3
5
3 22 2
2 7 1
2 87 1
3 14 4
2 51 2
4
4 47 3
4 5 4
1 20 2
3 38 2
2
3 80 2
3 47 4
5
1 25 2
3 69 2
2 4 1
1 38 2
3 6 4
3
4 18 2
4 45 4
3 58 1
6
2 4 4
4 15 2
3 77 3
1 62 1
2 24 1
2 77 3
1
1 56 4
2
2 18 1
2 30 3
4
3 22 3
3 34 3
3 70 1
2 74 3
2
4 16 4
1 21 2
2
2 69 1
3 91 1
1
2 34 2
5
1 61 3
3 83 2
4 68 3
1 37 2
4 1 1
3
3 26 1
4 55 1
1 26 3
6
2 26 2
1 2 1
1 67 3
4 12 1
3 81 2
2 89 1
1
1 56 2
2
1 66 3
3 18 3
5
1 78 1
1 6 2
3 100 4
3 98 3
2 69 3
4
2 22 4
4 75 1
2 76 3
4 21 3
3
2 39 3
3 69 2
4 35 3
3
2 13 3
2 82 3
2 57 2
1
3 67 1
1
1 84 1
6
2 32 2
3 50 4
2 94 4
4 65 3
1 6 3
1 8 2
4
1 21 4
4 72 1
2 44 3
3 1 2
3
3 98 1
4 63 2
3 71 4
6
1 21 4
1 50 4
2 12 1
3 17 1
2 57 4
3 41 4
4
3 94 4
1 49 3
1 56 4
3 74 2
5
2 37 4
1 69 3
1 79 2
2 30 4
3 57 4
4
1 2 4
4 61 3
3 19 2
3 20 3
3
4 63 2
4 45 4
1 33 4
5
2 3 2
4 17 3
1 30 1
3 41 2
3 74 4
2
3 86 1
1 57 1
1
2 61 2
5
1 23 2
2 25 1
2 2 2
4 45 3
3 74 4
2
3 79 1
3 13 3
1
3 100 1
4
1 92 2
2 63 1
4 9 3
1 93 2
4
4 90 4
1 19 2
1 81 2
4 57 4
3
4 82 2
3 74 3
1 51 3`

var testcases = mustParseTestcases(testcasesRaw)

func bitsCount(x int) int {
	cnt := 0
	for x != 0 {
		cnt++
		x &= x - 1
	}
	return cnt
}

// solve replicates the reference logic from 1038E.go.
func solve(tc testcase) int64 {
	dist := make([][]int64, V)
	for i := 0; i < V; i++ {
		dist[i] = make([]int64, V)
		for j := 0; j < V; j++ {
			if i == j {
				dist[i][j] = 0
			} else {
				dist[i][j] = INF
			}
		}
	}

	bMask := 0
	var total int64
	for _, e := range tc.edges {
		u := e.c1 - 1
		w := e.c2 - 1
		total += e.v
		bMask ^= 1 << u
		bMask ^= 1 << w
		if e.v < dist[u][w] {
			dist[u][w] = e.v
			dist[w][u] = e.v
		}
	}

	for k := 0; k < V; k++ {
		for i := 0; i < V; i++ {
			for j := 0; j < V; j++ {
				if dist[i][k]+dist[k][j] < dist[i][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
				}
			}
		}
	}

	best := int64(0)
	for fMask := 0; fMask < (1 << V); fMask++ {
		bc := bitsCount(fMask)
		if bc != 0 && bc != 2 {
			continue
		}
		xMask := bMask ^ fMask
		tcnt := bitsCount(xMask)
		var cost int64
		switch tcnt {
		case 0:
			cost = 0
		case 2:
			var nodes [2]int
			idx := 0
			for i := 0; i < V; i++ {
				if xMask&(1<<i) != 0 {
					nodes[idx] = i
					idx++
				}
			}
			cost = dist[nodes[0]][nodes[1]]
		case 4:
			c1 := dist[0][1] + dist[2][3]
			c2 := dist[0][2] + dist[1][3]
			c3 := dist[0][3] + dist[1][2]
			cost = c1
			if c2 < cost {
				cost = c2
			}
			if c3 < cost {
				cost = c3
			}
		default:
			continue
		}
		if cost >= INF/2 {
			continue
		}
		cur := total - cost
		if cur > best {
			best = cur
		}
	}
	return best
}

func mustParseTestcases(raw string) []testcase {
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(raw)))
	scanner.Split(bufio.ScanWords)

	nextInt := func() int {
		if !scanner.Scan() {
			panic("unexpected EOF while reading testcases")
		}
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic(fmt.Sprintf("invalid integer %q: %v", scanner.Text(), err))
		}
		return v
	}

	t := nextInt()
	cases := make([]testcase, 0, t)
	for i := 0; i < t; i++ {
		n := nextInt()
		edges := make([]edge, n)
		for j := 0; j < n; j++ {
			c1 := nextInt()
			v := nextInt()
			c2 := nextInt()
			edges[j] = edge{c1: c1, v: int64(v), c2: c2}
		}
		cases = append(cases, testcase{n: n, edges: edges})
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("scanner error: %v", err))
	}

	if len(cases) == 0 {
		panic("no testcases parsed")
	}
	return cases
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func parseCandidateOutput(out string) (int64, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return 0, fmt.Errorf("no output")
	}
	val, err := strconv.ParseInt(scanner.Text(), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse output: %v", err)
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scanner error: %v", err)
	}
	return val, nil
}

func checkCase(bin string, idx int, tc testcase) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.c1, e.v, e.c2))
	}

	expected := solve(tc)
	out, err := runCandidate(bin, sb.String())
	if err != nil {
		return err
	}
	got, err := parseCandidateOutput(out)
	if err != nil {
		return err
	}
	if got != expected {
		return fmt.Errorf("case %d: expected %d got %d", idx+1, expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	for i, tc := range testcases {
		if err := checkCase(exe, i, tc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
