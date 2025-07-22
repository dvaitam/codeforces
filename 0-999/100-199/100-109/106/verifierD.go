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

type cmd struct {
	dir  byte
	dist int
}

type testcase struct {
	n, m int
	grid []string
	cmds []cmd
}

func parseTests(path string) ([]testcase, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	var tests []testcase
	for {
		if !scan.Scan() {
			break
		}
		line := strings.TrimSpace(scan.Text())
		if line == "" {
			continue
		}
		var n, m int
		fmt.Sscan(line, &n, &m)
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			grid[i] = strings.TrimSpace(scan.Text())
		}
		scan.Scan()
		kLine := strings.TrimSpace(scan.Text())
		k, _ := strconv.Atoi(kLine)
		cmds := make([]cmd, k)
		for i := 0; i < k; i++ {
			scan.Scan()
			var d string
			var dist int
			fmt.Sscan(scan.Text(), &d, &dist)
			cmds[i] = cmd{dir: d[0], dist: dist}
		}
		// consume blank line
		scan.Scan()
		tests = append(tests, testcase{n, m, grid, cmds})
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return tests, nil
}

func expected(t testcase) string {
	n, m := t.n, t.m
	grid := t.grid
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
				pos[idx] = [2]int{i, j}
				ok[idx] = true
			}
		}
	}
	dmap := map[byte][2]int{'N': {-1, 0}, 'S': {1, 0}, 'W': {0, -1}, 'E': {0, 1}}
	for _, cm := range t.cmds {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTests("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, t := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.m))
		for i := 0; i < t.n; i++ {
			sb.WriteString(t.grid[i] + "\n")
		}
		sb.WriteString(fmt.Sprintf("%d\n", len(t.cmds)))
		for _, c := range t.cmds {
			sb.WriteString(fmt.Sprintf("%c %d\n", c.dir, c.dist))
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		want := expected(t)
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
