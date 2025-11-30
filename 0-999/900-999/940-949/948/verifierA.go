package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesA = `1 1 .
1 1 S
1 1 W
1 2 SW
1 2 WS
2 2 SW WS
2 2 S. .W
2 2 SS WW
3 3 ... ... ...
2 3 S.W W.S
3 2 SW SW .S
5 5 S.SWS WSWW. WS.WW .SS.S WSSWS
5 4 .W.W .SS. .W.. .SSW .W.S
1 2 SW
5 3 WS. WSS WW. WSW .SS
3 2 .W WW S.
4 4 SW.W WSWW SWSW W.WS
5 4 W.WW .SS. S..S S..S WW.S
2 3 W.W SWW
2 4 .SWW W...
3 1 S S W
3 3 SWW SW. ..S
2 5 .SWW. S.WWS
3 1 S W W
1 5 WW.W.
1 5 WWSWS
1 1 .
3 5 S..S. WSS.. .W...
4 1 W . W W
5 5 ...S. .W..S S..W. WSSWS WSWWS
5 3 .SW WWS SWS ..W W.W
4 5 SSSS. .S... SWSS. S.WSW
4 2 .S .W .. ..
3 3 SS. W.. .S.
5 2 SW SS .S W. .S
4 4 SS.. .W.. WS.W S..S
3 5 S..W. .SW.W WS.S.
2 4 S.SS .SWW
5 3 W.W SWS WS. SS. WSS
4 2 .W W. WW WW
4 1 S . S S
1 5 ..S.S
3 5 S.WW. W..SW WW.SW
5 3 WW. WS. S.. ..S WWS
3 3 SWW .W. W.W
5 3 W.W ..S .WW ..S .SW
5 3 WWS WS. ... SSS WSS
4 4 .SSS SSSW ...W ...W
2 4 W... WSSW
5 1 S S . W .
5 2 WS W. .W SW .S
4 3 W.. W.S WW. WS.
2 4 .S.W SSWS
4 5 WWSWW .S.WW S..SS W.SWS
1 2 W.
1 2 S.
4 5 .WSSW .WSSS SSS.S ..SSW
5 2 .S .. W. W. WS
2 5 ..W.. SWS.S
2 1 W W
2 2 SW ..
2 4 SWWW SSSW
1 1 .
1 4 .WW.
5 5 WSW.. .S... S..SS W..W. SWWWS
4 2 WW SS SS S.
2 1 S S
4 5 .W.WS .WWWS .S.S. .S.WS
3 4 SSSS WS.. ...W
5 4 S.WW WS.W .WWW SW.W W...
1 4 .SSW
5 2 WW S. .. SW WW
3 3 .SS .SW S..
2 4 WWWS SW..
5 1 W W . S .
2 3 .WS S.S
1 5 SSSS.
2 3 ..W .W.
5 1 W W S . S
1 1 .
1 2 S.
2 3 SWS WW.
3 1 . S .
2 1 W W
2 2 SS ..
5 5 SS.W. WS.S. W..WW SW.SS W.SWW
3 3 .WS .WW .S.
1 4 S.SS
1 3 S.W
5 2 SS .S W. SS SW
5 5 WSW.W .SS.. ..SS. WSWS. .S..S
5 4 WSS. ..WW SWW. SWW. .W..
1 3 .SS
1 2 SS
5 2 S. W. .. S. SW
5 5 ..SWW SWSW. .SWS. WS.WS W..SW
3 4 SWSW .WSW W...
4 1 S S . .
1 5 S.W.W
2 3 WWS ...
3 5 WS.SS W.W.W S.W..
3 3 .SS ..S S..
3 2 SW W. S.
2 2 S. S.
2 1 W W
3 2 S. WS SS
4 3 WS. W.W .S. ..W
5 2 WW .. .. WW SS
4 4 WS.. W.SS .WS. SSS.
2 4 SWS. WSSW`

type testCase struct {
	n, m int
	grid []string
}

func solve(tc testCase) (string, [][]byte) {
	g := make([][]byte, tc.n)
	for i := 0; i < tc.n; i++ {
		g[i] = []byte(tc.grid[i])
	}
	dirs := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			if g[i][j] == 'W' {
				for _, d := range dirs {
					ni, nj := i+d[0], j+d[1]
					if ni >= 0 && ni < tc.n && nj >= 0 && nj < tc.m {
						if g[ni][nj] == 'S' {
							return "No", nil
						}
						if g[ni][nj] == '.' {
							g[ni][nj] = 'D'
						}
					}
				}
			}
		}
	}
	return "Yes", g
}

func parseTests() ([]testCase, error) {
	fields := strings.Fields(testcasesA)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty testcases")
	}
	ptr := 0
	var tests []testCase
	for ptr < len(fields) {
		if ptr+2 > len(fields) {
			return nil, fmt.Errorf("incomplete header")
		}
		n, err := strconv.Atoi(fields[ptr])
		if err != nil {
			return nil, err
		}
		m, err := strconv.Atoi(fields[ptr+1])
		if err != nil {
			return nil, err
		}
		ptr += 2
		if ptr+n > len(fields) {
			return nil, fmt.Errorf("missing rows")
		}
		grid := fields[ptr : ptr+n]
		ptr += n
		tests = append(tests, testCase{n: n, m: m, grid: grid})
	}
	return tests, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTests()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		verdict, expectedGrid := solve(tc)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for _, row := range tc.grid {
			input.WriteString(row)
			input.WriteByte('\n')
		}
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(got), "\n")
		if len(lines) == 0 {
			fmt.Fprintf(os.Stderr, "case %d failed: empty output\n", i+1)
			os.Exit(1)
		}
		if strings.TrimSpace(lines[0]) != verdict {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, verdict, strings.TrimSpace(lines[0]))
			os.Exit(1)
		}
		if verdict == "No" {
			continue
		}
		if len(lines)-1 != tc.n {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d grid lines got %d\n", i+1, tc.n, len(lines)-1)
			os.Exit(1)
		}
		for r := 0; r < tc.n; r++ {
			row := strings.TrimSpace(lines[r+1])
			if len(row) != tc.m {
				fmt.Fprintf(os.Stderr, "case %d failed: row %d length mismatch\n", i+1, r)
				os.Exit(1)
			}
			for c := 0; c < tc.m; c++ {
				if row[c] != expectedGrid[r][c] {
					fmt.Fprintf(os.Stderr, "case %d failed: cell %d,%d expected %c got %c\n", i+1, r, c, expectedGrid[r][c], row[c])
					os.Exit(1)
				}
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
