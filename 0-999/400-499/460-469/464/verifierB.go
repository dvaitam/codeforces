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

func parseTests(path string) ([][8][3]int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
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
	tests, err := parseTests("testcasesB.txt")
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
