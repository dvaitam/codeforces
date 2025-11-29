package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesRaw = `1 2 2
7 5 5 2 5 2 1 3 4 3 1 4 0 1 2
4 5 1 2 1 2 0 0 1
3 2 5 5 0 5 3
6 0 2 1 3 1 1 2 0 0 2 3 3
3 3 5 2 4 4 1
1 0 1
4 2 0 4 4 5 2 4 1
6 3 3 1 2 2 2 2 1 2 1 0 2
7 3 3 2 3 5 0 5 0 4 1 1 3 1 4
3 2 0 1 0 5 2
4 2 5 4 2 5 5 2 1
6 0 4 2 1 0 3 3 0 3 5 5 5
4 0 0 3 5 5 3 1 2
2 4 3 3 0
6 2 5 3 4 2 3 3 3 2 4 2 0
1 0 0
2 5 3 3 0
5 3 5 0 4 1 1 5 3 4 5
3 1 5 2 1 1 3
4 1 0 2 1 3 0 5 1
3 0 5 0 2 2 2
2 4 5 1 3
4 5 4 0 2 5 3 4 0
3 2 4 4 3 2 2
1 3 5
2 2 3 0 1
7 0 0 0 2 3 1 2 0 0 0 2 3 2 5
2 3 2 2 2
4 3 4 4 2 4 3 3 1
6 5 1 5 3 2 1 0 4 1 1 2 2
3 0 2 0 0 1 0
5 5 2 4 0 2 4 5 3 3 0
2 2 5 5 1
2 4 2 1 0
1 0 1
1 5 5
1 3 0
3 1 4 2 5 5 3
2 1 5 4 5
5 4 1 3 3 1 0 4 4 4 3
5 5 3 1 3 4 4 4 5 0 1
5 1 2 4 1 4 3 2 5 4 3
6 3 2 1 3 3 4 3 4 5 1 5 3
1 0 5
6 4 3 4 0 2 2 2 3 5 5 4 2
7 4 4 0 0 0 5 2 1 2 5 1 5 0 0
5 1 0 2 3 2 0 4 4 0 0
6 1 5 1 0 3 1 5 4 5 4 2 0
2 1 3 4 0
7 4 1 4 3 1 3 4 2 1 1 0 2 3 4
6 5 0 3 1 3 3 4 1 3 0 4 0
5 1 4 4 2 0 1 0 4 0 1
6 5 1 3 2 1 0 3 1 5 1 5 1
4 4 3 2 2 1 5 5 2
3 5 3 2 1 1 5
2 3 5 5 1
4 2 1 1 0 2 5 2 5
5 5 0 2 1 3 1 5 0 5 1
5 1 2 3 3 1 4 1 1 5 3
5 4 1 5 0 4 0 3 1 0 1
6 4 3 1 3 0 3 4 4 0 2 3 0
7 2 3 0 3 1 5 1 5 4 2 3 2 5 0
1 1 4
6 5 5 0 0 3 2 1 0 2 1 2 3
5 3 4 1 1 3 0 3 3 0 3
2 0 1 1 2
5 4 1 0 3 4 3 3 1 0 0
1 0 3
5 0 0 1 3 2 2 0 4 5 4
5 5 3 4 4 5 5 2 5 0 5
3 2 0 2 0 5 0
4 2 0 3 4 2 1 1 5
2 4 4 4 1
7 2 5 3 3 1 1 2 1 5 0 5 1 1 0
3 3 3 2 5 0 1
2 2 5 1 0
1 2 1
5 0 3 2 2 4 0 5 5 3 2
3 3 2 3 2 2 1
5 2 5 5 3 5 1 0 0 2 2
7 1 4 5 4 3 2 3 1 0 1 4 2 1 3
3 4 5 2 0 5 2
1 4 2
6 2 4 5 2 5 4 1 5 5 4 1 2
1 4 1
3 5 4 5 1 4 2
3 1 3 2 4 3 2
3 5 3 3 2 1 5
5 5 0 2 0 2 3 1 4 3 4
7 0 2 1 1 0 2 2 4 0 2 4 0 1 2
4 0 4 1 3 5 1 2 0
1 4 2
5 4 5 2 2 0 4 4 0 4 3
2 3 5 5 3
3 1 3 0 1 1 3
7 5 0 0 5 2 0 2 5 1 3 0 0 1 0
2 1 1 4 2
7 0 4 1 5 5 2 5 1 4 1 2 0 1 4
6 5 3 1 0 4 2 0 4 2 5 1 1
`

type Point struct {
	x int
	y int
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// solve replicates the logic from 1066F.go.
func solve(points []Point) int64 {
	n := len(points)
	p := make([]Point, n+1)
	copy(p[1:], points)

	sort.Slice(p[1:], func(i, j int) bool {
		a, b := p[i+1], p[j+1]
		fa, fb := max(a.x, a.y), max(b.x, b.y)
		if fa != fb {
			return fa < fb
		}
		if a.x != b.x {
			return a.x < b.x
		}
		return a.y > b.y
	})

	const INF = int64(1e18)
	dp := make([]int64, n+2)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0
	pl, pr := 0, 0

	le := func(d int) int {
		return max(p[d].x, p[d].y)
	}
	get := func(a, b int) int {
		return abs(p[a].x-p[b].x) + abs(p[a].y-p[b].y)
	}

	for pr < n {
		d := pr + 1
		for d < n && le(d+1) == le(pr+1) {
			d++
		}
		c1 := dp[pl] + int64(get(pl, pr+1))
		c2 := dp[pr] + int64(get(pr, pr+1))
		costSeg := int64(get(pr+1, d))
		if c1 < c2 {
			dp[d] = c1 + costSeg
		} else {
			dp[d] = c2 + costSeg
		}
		c3 := dp[pl] + int64(get(pl, d))
		c4 := dp[pr] + int64(get(pr, d))
		if c3 < c4 {
			dp[pr+1] = c3 + costSeg
		} else {
			dp[pr+1] = c4 + costSeg
		}
		pl = pr + 1
		pr = d
	}
	res := dp[pl]
	if dp[pr] < res {
		res = dp[pr]
	}
	return res
}

type testCase struct {
	points []Point
}

func parseTestcases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	var tests []testCase
	for {
		if !scan.Scan() {
			break
		}
		n, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, err
		}
		pts := make([]Point, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("invalid test file")
			}
			x, _ := strconv.Atoi(scan.Text())
			if !scan.Scan() {
				return nil, fmt.Errorf("invalid test file")
			}
			y, _ := strconv.Atoi(scan.Text())
			pts[i] = Point{x: x, y: y}
		}
		tests = append(tests, testCase{points: pts})
	}
	return tests, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.points)))
	for _, p := range tc.points {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
	}
	return sb.String()
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	expected := strconv.FormatInt(solve(tc.points), 10)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
