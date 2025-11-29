package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesData = `
5 4
####
#.##
#.A#
#..#
####
2
N 2
E 1

3 5
#####
#AB.#
#####
5
E 1
E 1
S 2
N 1
S 2

6 4
####
##.#
#.##
##.#
#A.#
####
2
N 2
E 1

6 4
####
#.A#
#..#
####
#B##
####
4
W 2
E 2
S 2
N 1

6 3
###
#.#
#C#
#A#
#B#
###
1
E 2

6 4
####
#C.#
#..#
#..#
#AB#
####
5
S 2
S 1
W 1
E 2
E 1

3 5
#####
#CBA#
#####
5
S 2
E 2
S 1
W 1
W 1

4 4
####
#.A#
#CB#
####
5
E 1
E 1
N 2
W 1
S 1

6 4
####
#A##
#..#
#..#
#..#
####
1
W 2

3 5
#####
#.A##
#####
4
E 1
N 1
W 1
E 1

4 4
####
#AB#
#..#
####
4
E 2
E 2
S 2
E 1

4 6
######
#.#.##
#A#.##
######
1
W 1

4 3
###
#.#
#A#
###
3
N 2
S 1
W 1

4 3
###
#A#
#B#
###
1
S 1

6 5
#####
#.A.#
#...#
#...#
#.BC#
#####
2
W 2
S 1

3 3
###
#A#
###
5
W 1
W 2
W 2
E 1
S 1

3 4
####
#BA#
####
4
W 1
S 1
E 1
W 1

5 5
#####
#..A#
#...#
##CB#
#####
5
N 2
W 2
W 1
S 2
N 1

6 4
####
#.B#
#A##
#..#
#C.#
####
3
E 2
N 1
W 1

3 3
###
#A#
###
3
W 2
W 2
S 1

6 4
####
#A.#
#..#
#..#
#..#
####
1
S 2

6 5
#####
#..##
#...#
#...#
#AB.#
#####
3
N 1
E 2
W 2

3 6
######
##A.B#
######
3
E 2
N 2
E 1

5 4
####
#A.#
#..#
##.#
####
4
W 1
S 1
W 2
N 2

6 5
#####
#...#
#...#
#...#
#.BA#
#####
4
E 2
E 2
N 2
N 1

6 6
######
#.#.C#
##..B#
##.A##
#...##
######
4
E 2
E 2
S 1
E 2

4 4
####
#B##
#.A#
####
3
W 2
W 1
N 2

3 4
####
#A.#
####
1
N 2

6 3
###
#.#
#.#
#.#
#A#
###
2
E 2
W 2

5 3
###
#B#
#A#
#.#
###
5
W 2
S 1
E 1
E 1
N 2

6 5
#####
#...#
#...#
###A#
#..##
#####
1
E 1

3 5
#####
#ABC#
#####
3
W 2
W 2
S 1

6 5
#####
#.#.#
#C..#
#A..#
#.B.#
#####
3
N 2
S 1
E 1

3 3
###
#A#
###
5
E 1
S 1
W 2
S 2
W 1

6 3
###
#A#
#B#
#.#
#.#
###
5
N 1
W 1
W 1
N 1
S 1

6 5
#####
#.B.#
#...#
##..#
#..A#
#####
5
W 1
N 1
N 2
E 2
E 1

6 5
#####
#.#.#
#.A.#
#.#.#
#...#
#####
2
N 2
E 1

5 3
###
#.#
#.#
#A#
###
4
W 2
W 1
W 2
S 1

3 6
######
#A.CB#
######
1
E 2

4 6
######
#.#A.#
#..###
######
4
E 1
E 2
W 1
W 2

6 6
######
#.B.##
##.A##
#.#..#
#....#
######
2
S 2
N 2

5 6
######
#.CA.#
#.#..#
#.B#.#
######
3
S 2
S 2
E 1

6 6
######
#.AB.#
#....#
#..C##
#....#
######
1
W 2

3 3
###
#A#
###
5
W 1
S 2
N 2
E 2
W 2

5 5
#####
#..A#
#...#
#.B.#
#####
4
E 1
N 1
S 2
S 1

4 4
####
#..#
#A.#
####
3
W 2
W 1
N 2

4 4
####
#..#
#BA#
####
3
E 2
S 2
W 1

5 3
###
###
#A#
#B#
###
3
N 1
E 2
N 1

5 6
######
#.AC.#
#....#
#.#.B#
######
5
N 1
S 2
W 1
W 2
E 2

5 4
####
##.#
#A.#
#..#
####
3
S 1
W 2
W 1

4 4
####
#A.#
#B.#
####
3
W 2
E 1
N 1

3 5
#####
#A..#
#####
4
E 1
W 1
W 1
S 1

5 5
#####
#.B.#
#..A#
#...#
#####
3
S 1
S 1
W 2

6 6
######
#....#
#....#
#..B.#
#...A#
######
4
S 1
S 1
N 2
E 2

6 5
#####
#...#
#A###
#...#
##..#
#####
4
S 1
N 2
E 1
N 2

3 5
#####
#BCA#
#####
2
N 1
N 1

4 4
####
#A.#
#.B#
####
5
S 1
E 2
N 2
E 1
W 2

5 3
###
#.#
###
#A#
###
4
E 2
S 2
N 1
W 2

4 4
####
#A.#
#.B#
####
1
N 2

4 6
######
#....#
#.A#.#
######
1
E 1

3 6
######
#..A.#
######
1
E 2

5 3
###
#.#
#B#
#A#
###
2
N 1
W 1

3 3
###
#A#
###
1
S 2

6 3
###
#A#
#B#
#.#
#.#
###
5
S 2
W 2
W 2
W 1
S 1

4 6
######
#..A##
#....#
######
5
S 2
E 1
S 2
E 2
E 1

3 6
######
#CB.A#
######
4
S 2
W 2
N 2
W 2

6 3
###
#A#
#B#
#.#
###
###
3
W 1
N 1
E 1

3 4
####
#BA#
####
2
W 1
W 2

4 4
####
#.A#
##.#
####
3
E 2
E 1
N 2

3 6
######
#.A..#
######
2
S 2
N 2

5 5
#####
#...#
#A###
#.B.#
#####
2
W 2
S 2

3 3
###
#A#
###
4
E 2
W 2
W 1
E 2

4 5
#####
#...#
#.A.#
#####
1
S 2

6 4
####
##A#
#.##
#.##
#.##
####
5
E 1
N 2
W 1
N 1
N 1

4 4
####
#.A#
#.B#
####
4
E 1
N 1
S 1
E 1

4 3
###
#.#
#A#
###
5
S 1
E 2
S 2
S 2
S 2

4 4
####
#.A#
#..#
####
5
E 1
N 1
W 1
E 1
W 1

6 6
######
#...A#
#...##
#...B#
#....#
######
3
S 1
S 2
N 2

5 3
###
#A#
#.#
#B#
###
3
W 2
S 1
S 2

6 5
#####
##..#
##..#
##AB#
#C..#
#####
2
W 2
E 2

6 3
###
#B#
#.#
#A#
#.#
###
1
W 2

5 6
######
#.#.##
#...##
#.A.##
######
1
W 2

5 3
###
#C#
#B#
#A#
###
5
S 2
N 2
S 2
N 1
W 2

6 4
####
##.#
#.##
#A.#
#.B#
####
5
N 1
E 1
N 1
N 2
E 1

3 6
######
#A#..#
######
3
S 2
N 2
W 2

4 6
######
#.##.#
#A.#.#
######
1
W 1

5 5
#####
#.###
#..##
#A#.#
#####
5
N 2
N 1
E 1
N 1
E 2

4 3
###
#A#
#.#
###
2
S 1
W 1

6 4
####
##A#
#.B#
#.##
#.C#
####
2
N 2
E 2

5 3
###
###
#.#
#A#
###
3
W 2
E 2
S 2

3 4
####
#AB#
####
3
W 2
E 2
N 1

5 4
####
#C.#
#BA#
#..#
####
2
N 2
E 2

4 3
###
#B#
#A#
###
1
W 1

6 3
###
#A#
#C#
#B#
###
###
2
N 1
N 2

6 3
###
#.#
###
#A#
#.#
###
2
E 1
E 2

4 6
######
#BC#.#
#.A..#
######
1
W 1

3 6
######
#.A#B#
######
2
W 2
W 1

3 6
######
#..BA#
######
5
S 1
E 1
E 2
S 2
W 2

6 4
####
#.##
#..#
#A.#
#.##
####
4
N 1
W 1
N 2
E 2

5 3
###
#C#
#A#
#B#
###
2
E 2
S 1
`

type cmd struct {
	dir  byte
	dist int
}

type testCase struct {
	n, m int
	grid []string
	cmds []cmd
}

func parseTestcases() ([]testCase, error) {
	scanner := bufio.NewScanner(strings.NewReader(testcasesData))
	scanner.Split(bufio.ScanWords)
	nextInt := func() (int, error) {
		if !scanner.Scan() {
			return 0, io.EOF
		}
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return 0, err
		}
		return v, nil
	}

	var cases []testCase
	for {
		n, err := nextInt()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		m, err := nextInt()
		if err != nil {
			return nil, err
		}
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			if !scanner.Scan() {
				return nil, fmt.Errorf("unexpected EOF reading grid")
			}
			grid[i] = scanner.Text()
		}
		k, err := nextInt()
		if err != nil {
			return nil, err
		}
		cmds := make([]cmd, k)
		for i := 0; i < k; i++ {
			if !scanner.Scan() {
				return nil, fmt.Errorf("unexpected EOF reading direction")
			}
			dir := scanner.Text()
			dist, err := nextInt()
			if err != nil {
				return nil, err
			}
			cmds[i] = cmd{dir: dir[0], dist: dist}
		}
		cases = append(cases, testCase{n: n, m: m, grid: grid, cmds: cmds})
	}
	return cases, nil
}

// solve mirrors 106D.go so the verifier is self-contained.
func solve(tc testCase) string {
	n, m := tc.n, tc.m
	grid := tc.grid
	sum := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		sum[i] = make([]int, m+1)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			add := 0
			if grid[i-1][j-1] == '#' {
				add = 1
			}
			sum[i][j] = sum[i-1][j] + sum[i][j-1] - sum[i-1][j-1] + add
		}
	}
	pos := make([][2]int, 26)
	ok := make([]bool, 26)
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			c := grid[i-1][j-1]
			if c >= 'A' && c <= 'Z' {
				idx := int(c - 'A')
				pos[idx][0] = i
				pos[idx][1] = j
				ok[idx] = true
			}
		}
	}
	dmap := map[byte][2]int{
		'N': {-1, 0},
		'S': {1, 0},
		'W': {0, -1},
		'E': {0, 1},
	}
	for _, cm := range tc.cmds {
		delta := dmap[cm.dir]
		dx, dy := delta[0], delta[1]
		for idx := 0; idx < 26; idx++ {
			if !ok[idx] {
				continue
			}
			x0, y0 := pos[idx][0], pos[idx][1]
			x1 := x0 + dx*cm.dist
			y1 := y0 + dy*cm.dist
			if x1 < 1 || x1 > n || y1 < 1 || y1 > m {
				ok[idx] = false
			} else {
				xa, xb := x0, x1
				if xa > xb {
					xa, xb = xb, xa
				}
				ya, yb := y0, y1
				if ya > yb {
					ya, yb = yb, ya
				}
				if sum[xb][yb]-sum[xa-1][yb]-sum[xb][ya-1]+sum[xa-1][ya-1] > 0 {
					ok[idx] = false
				}
			}
			pos[idx][0], pos[idx][1] = x1, y1
		}
	}
	var res []byte
	for idx := 0; idx < 26; idx++ {
		if ok[idx] {
			res = append(res, byte('A'+idx))
		}
	}
	if len(res) == 0 {
		return "no solution"
	}
	return string(res)
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, row := range tc.grid {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.cmds)))
	for _, c := range tc.cmds {
		sb.WriteString(fmt.Sprintf("%c %d\n", c.dir, c.dist))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		input := buildInput(tc)
		want := solve(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
