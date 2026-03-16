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
const testcasesBRaw = `
1 6 -2
4 -1 -2
-4 4 -1
6 -4 -1
1 6 -4
-2 6 -1
1 4 -2
-4 1 4
1 0 5
5 -5 6
6 0 5
-5 10 1
5 1 -5
1 0 10
10 6 0
10 6 -5
3 1 5
10 1 3
6 3 10
10 6 8
6 5 3
5 8 6
8 1 5
10 8 1
-5 5 2
-2 5 5
-2 5 2
2 5 -2
5 -5 2
5 5 -5
-2 2 2
2 2 -5
8 6 3
8 1 3
6 -2 8
6 3 -2
-2 1 8
3 3 1
6 3 3
-2 1 3
3 -4 -4
-4 2 -4
3 -5 -4
-4 3 -5
-5 2 -5
-5 -5 3
-5 -4 2
2 -5 -4
0 4 4
4 3 4
0 1 4
0 1 4
4 1 3
0 1 1
1 3 1
3 4 1
4 5 8
0 4 9
4 9 8
5 4 4
8 0 5
4 9 4
4 5 0
9 8 0
3 6 -5
-3 4 5
6 3 -3
6 -3 5
-5 3 4
4 5 -5
-3 4 3
5 -5 6
7 -3 0
0 7 3
-3 3 7
0 3 4
-3 0 4
0 4 0
4 3 -3
7 0 0
-1 3 0
-4 0 0
-4 0 3
-4 3 3
-1 3 0
3 -1 3
0 0 -1
3 -4 0
-1 3 4
6 -1 4
3 -4 7
7 6 -4
3 -1 7
-4 4 3
6 4 -4
6 -1 7
7 3 5
3 3 1
5 7 -1
1 -1 3
1 -1 7
7 3 1
5 3 3
-1 3 5
2 5 6
2 3 4
4 5 4
4 5 6
2 3 6
4 4 3
2 4 5
6 4 3
3 0 -2
-4 2 3
-2 5 0
2 5 -4
0 3 -4
-2 3 2
2 -2 5
-4 5 0
-3 2 -3
2 -3 -1
0 -1 -1
0 -3 -1
-1 2 -1
2 -1 -3
-3 -3 0
-1 0 -3
7 8 2
8 7 7
2 3 7
3 7 7
3 2 7
2 8 2
2 7 8
2 3 2
-2 6 -2
0 -4 6
4 0 -2
-2 -2 4
6 -2 -4
-4 4 0
-4 -2 4
6 -2 0
0 3 1
1 -2 3
-2 -1 3
5 -1 0
0 1 5
3 -1 0
5 -2 1
-1 5 -2
1 -5 3
0 1 8
1 -5 8
0 -4 8
-4 8 -5
3 -4 -5
0 1 3
-4 3 0
3 1 4
3 4 -2
1 6 1
1 6 4
-2 4 6
1 1 3
6 1 -2
-2 3 1
-2 4 2
5 -1 2
1 5 -1
-1 4 1
-1 2 4
-2 4 1
5 1 -2
-2 2 5
1 -1 6
3 -1 4
1 -1 3
-1 4 6
2 3 1
6 2 4
1 2 6
4 2 3
4 5 1
5 4 5
5 4 1
4 1 1
5 8 5
1 1 8
8 1 5
8 1 5
-1 2 2
7 2 4
7 -3 4
-3 2 -1
2 4 2
-3 -1 7
-1 7 2
4 -3 2
0 0 3
5 -2 -2
0 3 -2
5 -2 0
0 0 5
3 -2 -2
0 3 -2
0 5 -2
4 2 1
1 5 0
4 1 0
1 5 1
2 0 4
4 1 1
5 1 2
0 5 2
5 -1 -1
6 -1 -1
5 0 -1
-2 -1 6
-2 0 6
6 0 -1
5 -1 -2
0 -2 5
6 1 1
6 -3 5
1 2 -3
1 -3 6
5 1 6
-3 5 2
2 1 1
2 5 1
1 -3 7
7 3 -1
5 -3 3
1 -1 7
-3 7 3
-3 1 5
-1 3 5
5 -1 1
2 -4 -2
2 -1 -3
2 -3 -2
-2 1 -3
-4 -1 1
-4 2 -1
1 -1 -3
-4 1 -2
4 6 3
6 3 1
4 3 3
1 3 3
6 6 4
1 3 6
6 1 6
3 4 6
-2 -4 1
-3 -5 1
-3 0 -4
-3 1 -4
-2 0 -5
1 -5 -2
-2 0 -4
0 -3 -5
5 -1 -4
5 0 -4
0 -3 5
6 -1 -3
-1 6 -4
-3 5 -1
-3 0 6
6 -4 0
2 -3 1
1 -2 1
1 2 -2
0 2 -2
-3 0 2
1 1 -3
-3 1 0
1 -2 0
2 4 5
1 4 6
2 1 4
6 5 0
6 5 4
0 2 1
0 1 6
5 2 0
1 -3 2
2 -4 -3
1 -3 -3
1 2 2
-4 -3 -3
2 -4 -3
2 -4 2
2 -3 1
0 2 4
4 2 -2
0 4 4
-2 4 4
2 6 -2
0 6 2
4 6 0
4 -2 6
-1 5 -5
1 7 -3
-1 7 -3
1 -3 5
-5 5 1
7 -1 -5
7 1 -5
-3 -1 5
-1 4 5
9 -5 4
9 -1 8
5 4 -5
5 -5 8
-5 9 8
5 -1 8
-1 4 9
-5 -3 -1
-5 -1 1
-3 -1 1
-5 -1 -1
-1 -3 -3
-5 -3 1
-3 -1 -1
-3 1 -3
-3 0 0
-3 -2 0
-2 0 -1
-1 -2 0
-2 -3 0
-2 -2 -1
-1 0 0
-2 -3 -2
-1 0 6
0 5 -1
0 0 6
0 -1 5
-1 5 -1
5 0 0
0 -1 6
-1 6 -1
3 -1 6
7 -1 6
2 -1 3
6 7 3
2 3 3
2 7 -1
2 7 3
6 3 3
3 7 5
1 5 5
1 7 5
7 3 1
3 7 3
5 5 3
3 5 3
1 3 5
-3 3 -5
0 3 -3
8 0 -3
0 3 2
-5 -3 8
0 2 8
-5 2 3
-5 8 2
5 4 4
9 -1 5
4 4 10
10 4 9
9 4 5
5 -1 4
-1 4 10
10 9 -1
-1 -1 -1
-4 -1 -1
-4 -1 -1
-4 -4 -4
-4 -1 -1
-4 -4 -1
-4 -1 -4
-4 -1 -4
-1 6 -4
1 -1 -4
1 1 -1
1 1 4
-4 6 4
1 -1 6
4 -4 1
4 6 1
-3 1 0
-4 1 5
5 0 -3
5 -3 -4
1 1 -4
1 0 5
-3 -4 1
1 0 1
1 2 1
-4 4 -3
-1 -1 5
5 5 4
-5 3 -5
5 -3 1
3 -4 2
-5 1 4
5 1 -1
0 1 1
4 2 -5
-4 2 -5
5 -5 -5
-4 4 -3
3 3 0
3 -1 4
5 0 2
-2 4 -2
-4 3 0
-3 -4 -5
0 1 0
-1 5 5
-5 4 1
1 1 0
-1 0 2
-2 5 4
3 -3 -5
0 5 -4
3 -3 3
5 5 2
0 -4 4
-5 2 -2
1 5 -3
1 -2 -4
-2 0 0
5 -2 5
2 2 0
2 5 5
-2 1 2
1 3 -4
4 2 -1
-3 -3 -5
1 1 -4
-5 5 -4
-3 2 1
5 3 -1
-3 -3 3
-4 -1 -5
2 1 5
-2 3 1
-5 3 -2
1 -3 5
-3 0 5
-2 -4 3
3 -3 -3
1 4 -5
3 -2 1
-2 -5 3
-2 3 4
5 3 -4
-2 1 2
-4 4 5
-5 1 -4
3 -4 5
2 -5 3
-2 -5 -5
-1 2 -1
1 -3 4
-3 3 0
3 5 2
3 1 3
-3 1 1
-2 2 -1
0 -3 -1
4 -1 -3
4 -4 0
0 -3 -1
-1 -1 0
1 -1 4
2 -5 -3
-3 -1 -2
-2 -4 4
3 4 -2
3 1 -2
4 -3 3
2 1 -2
-4 5 -4
-3 5 -5
-5 1 1
1 5 -3
4 4 -3
5 3 3
-4 -2 1
-3 -1 -2
5 1 0
-3 -2 -1
-3 0 2
3 -1 -4
3 -1 -2
2 -5 -1
4 4 -4
4 0 2
-1 4 -5
-5 0 -3
-3 5 -4
-4 1 5
4 -2 -2
3 3 1
-4 -2 1
5 3 -3
4 -1 -5
-4 -2 4
1 5 2
3 4 -2
-1 -5 5
-3 5 5
3 3 -2
1 -1 5
1 1 -1
2 -4 5
-3 -3 3
-5 2 -5
2 -2 1
3 0 -2
-4 -4 5
-5 1 2
-2 -3 4
3 -2 3
1 3 0
-2 -2 0
5 4 -4
0 -5 2
-5 4 -3
-3 -1 2
-5 4 3
-4 4 1
-4 1 3
4 5 -1
1 -1 0
2 -5 3
2 -5 1
-1 4 0
-3 4 4
3 -1 -4
4 0 1
1 3 -5
4 4 -4
-5 4 3
-5 -4 0
0 0 3
-5 5 0
4 -4 2
5 -4 3
2 0 3
3 -5 -3
0 0 -2
-3 4 -3
4 -4 1
0 3 1
0 0 -1
4 0 -5
-4 5 -2
-1 1 3
-1 4 4
-4 -4 -3
-1 1 -4
-3 -1 3
5 -1 -2
-2 -4 -1
2 -5 3
-1 -2 3
-4 3 0
0 -1 3
-3 -5 2
0 -5 -5
0 1 -3
3 -5 4
5 5 3
1 -3 -2
-2 -4 4
-3 4 3
-4 -1 2
-2 -5 0
2 0 4
0 -2 5
-5 -5 2
-5 -3 -1
3 -5 -5
-2 -4 3
-3 -5 3
-2 -2 2
-1 -2 2
3 0 0
1 5 -4
-2 4 -3
-2 5 4
-1 4 1
4 2 0
-5 2 -5
-4 5 5
4 5 4
1 4 0
0 -4 5
1 -2 3
2 4 4
5 3 3
2 4 5
4 2 4
2 -3 -1
5 3 -1
4 1 4
3 -1 -1
-1 -5 4
-5 2 2
0 -2 3
2 -2 2
0 5 -3
1 1 -5
5 -4 0
-5 -1 3
-5 -1 1
-5 0 0
-1 4 -5
-2 -4 0
-4 5 5
-4 -3 -1
1 4 0
-2 -5 5
-3 3 4
5 0 -1
-1 1 1
3 2 -4
-2 1 -2
4 -5 4
-2 5 -2
-2 1 1
-2 4 -3
-1 0 -5
5 -1 2
2 -3 5
-3 -5 0
1 3 0
3 2 0
4 -4 4
5 -1 3
5 -1 1
-5 -1 -4
5 2 -4
3 -2 4
5 -1 1
0 -2 -5
-4 4 3
3 3 -3
-3 -1 -5
-4 -2 -5
5 -5 1
-5 -4 -5
-5 -5 3
0 0 -5
4 -5 3
-2 2 -2
-1 -1 4
3 3 -1
-2 -3 -2
1 -5 -2
3 2 -5
0 0 1
-4 -5 4
-3 3 5
-4 -3 -2
-2 -3 -1
-4 -5 0
-3 -4 2
-3 -2 -5
-1 0 -5
4 -4 2
-2 -2 5
-3 -4 -5
-2 -5 -4
-4 -2 -1
-1 3 1
-2 -5 -1
-2 0 0
0 2 5
4 1 5
1 -4 1
-2 2 0
-3 4 5
-4 -2 -4
1 -1 3
-1 0 0
1 2 0
0 0 1
2 3 -5
0 -3 -1
-3 -1 4
-3 3 -3
-3 2 5
5 -3 -3
-3 -4 4
-1 -2 0
5 0 -3
-1 2 -1
-4 1 -3
3 0 2
-4 -3 5
0 -4 5
-3 2 3
-5 -5 -2
5 0 0
3 0 3
5 5 0
0 5 -4
-3 1 -5
-1 4 -2
-5 -2 -1
0 4 1
-2 0 -5
-2 -1 4
-5 -2 -4
-3 -2 0
3 -1 -3
-3 -2 -4
-1 4 3
3 3 4
3 1 2
4 3 2
-3 3 0
-2 1 -4
-1 -2 -2
-3 -3 -2
-5 -3 2
0 -3 -5
0 -4 4
-2 5 -2
-4 2 5
5 -2 4
0 -3 4
5 -5 -2
0 2 3
-5 -5 0
2 3 0
-3 2 -4
3 0 5
4 5 -1
4 0 4
-4 2 0
1 -4 -1
-4 5 5
0 -5 -3
0 -2 0
-1 -1 -1
2 1 -5
-1 -3 5
-1 -5 -4
1 1 4
-2 -1 0
5 4 2
4 -1 4
-1 5 -3
0 -3 0
-4 1 0
3 4 -2
1 2 -3
2 -2 -5
5 -2 -4
-4 -5 3
3 2 4
2 0 3
-3 4 2
1 -5 1
3 3 2
-3 4 4
0 -5 0
0 2 -2
5 5 3
-1 -4 2
0 -2 -3
-3 2 -5
0 4 0
-3 4 2
2 -5 4
-2 4 -5
2 5 -3
3 -2 1
2 -4 0
-1 -3 -3
0 -3 -3
4 3 -1
-2 3 1
2 2 3
3 -1 -3
3 4 3
-1 4 -2
-1 5 -3
5 -5 0
-4 1 1
5 3 -3
4 2 2
3 2 0
`


func isCube(points [8][3]int) bool {
	dists := make([]int64, 0, 28)
	for i := 0; i < 8; i++ {
		for j := i + 1; j < 8; j++ {
			var d int64
			for k := 0; k < 3; k++ {
				diff := int64(points[i][k] - points[j][k])
				d += diff * diff
			}
			dists = append(dists, d)
		}
	}
	sort.Slice(dists, func(i, j int) bool { return dists[i] < dists[j] })
	for i := 1; i < 12; i++ {
		if dists[i] != dists[i-1] {
			return false
		}
	}
	if dists[12] != 2*dists[11] {
		return false
	}
	for i := 13; i < 24; i++ {
		if dists[i] != dists[i-1] {
			return false
		}
	}
	if dists[24] != 3*dists[0] {
		return false
	}
	for i := 25; i < 28; i++ {
		if dists[i] != dists[i-1] {
			return false
		}
	}
	return dists[0] != 0
}

func solvePossible(orig [8][3]int) bool {
	var point [8][3]int
	perms := [6][3]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
	var dfs func(int) bool
	dfs = func(k int) bool {
		if k == 8 {
			return isCube(point)
		}
		for pi := 0; pi < 6; pi++ {
			p := perms[pi]
			for j := 0; j < 3; j++ {
				point[k][j] = orig[k][p[j]]
			}
			if dfs(k + 1) {
				return true
			}
		}
		return false
	}
	return dfs(0)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseTests() ([][8][3]int, error) {
	data := []byte(testcasesBRaw)
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	tests := make([][8][3]int, 0)
	for {
		var pts [8][3]int
		for i := 0; i < 8; i++ {
			for j := 0; j < 3; j++ {
				if !scan.Scan() {
					if i == 0 && j == 0 {
						return tests, nil
					}
					return nil, fmt.Errorf("incomplete test case")
				}
				val, _ := strconv.Atoi(scan.Text())
				pts[i][j] = val
			}
		}
		tests = append(tests, pts)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read testcases:", err)
		os.Exit(1)
	}
	for idx, tc := range tests {
		expected := solvePossible(tc)
		var sb strings.Builder
		for i := 0; i < 8; i++ {
			fmt.Fprintf(&sb, "%d %d %d\n", tc[i][0], tc[i][1], tc[i][2])
		}
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		scan := bufio.NewScanner(strings.NewReader(out))
		scan.Split(bufio.ScanWords)
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d: no output\n", idx+1)
			os.Exit(1)
		}
		word := strings.ToUpper(scan.Text())
		if word == "NO" {
			if expected {
				fmt.Fprintf(os.Stderr, "case %d: expected YES got NO\n", idx+1)
				os.Exit(1)
			}
			if scan.Scan() {
				fmt.Fprintf(os.Stderr, "case %d: extra output\n", idx+1)
				os.Exit(1)
			}
			continue
		}
		if word != "YES" {
			fmt.Fprintf(os.Stderr, "case %d: invalid output %s\n", idx+1, word)
			os.Exit(1)
		}
		if !expected {
			fmt.Fprintf(os.Stderr, "case %d: expected NO got YES\n", idx+1)
			os.Exit(1)
		}
		var pts [8][3]int
		for i := 0; i < 8; i++ {
			for j := 0; j < 3; j++ {
				if !scan.Scan() {
					fmt.Fprintf(os.Stderr, "case %d: missing coordinates\n", idx+1)
					os.Exit(1)
				}
				v, err := strconv.Atoi(scan.Text())
				if err != nil {
					fmt.Fprintf(os.Stderr, "case %d: bad number\n", idx+1)
					os.Exit(1)
				}
				pts[i][j] = v
			}
		}
		if scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d: extra output\n", idx+1)
			os.Exit(1)
		}
		for i := 0; i < 8; i++ {
			inRow := []int{tc[i][0], tc[i][1], tc[i][2]}
			outRow := []int{pts[i][0], pts[i][1], pts[i][2]}
			sort.Ints(inRow)
			sort.Ints(outRow)
			for j := 0; j < 3; j++ {
				if inRow[j] != outRow[j] {
					fmt.Fprintf(os.Stderr, "case %d: row %d not a permutation\n", idx+1, i+1)
					os.Exit(1)
				}
			}
		}
		if !isCube(pts) {
			fmt.Fprintf(os.Stderr, "case %d: output is not a cube\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
