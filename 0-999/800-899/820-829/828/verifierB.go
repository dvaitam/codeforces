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

func expected(grid []string, n, m int) int {
	minR, maxR := n, -1
	minC, maxC := m, -1
	count := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 'B' {
				count++
				if i < minR {
					minR = i
				}
				if i > maxR {
					maxR = i
				}
				if j < minC {
					minC = j
				}
				if j > maxC {
					maxC = j
				}
			}
		}
	}
	if count == 0 {
		return 1
	}
	height := maxR - minR + 1
	width := maxC - minC + 1
	side := height
	if width > side {
		side = width
	}
	if side > n || side > m {
		return -1
	}
	return side*side - count
}

func runCase(exe, input, expect string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect = strings.TrimSpace(expect)
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	if exe == "--" && len(os.Args) == 3 {
		exe = os.Args[2]
	}
	const testcasesRaw = `100
2 5
BWBWW
WWBBW
1 4
WBWW
2 5
BWBBB
BWBWB
5 2
WW
BW
BB
WW
BW
5 1
B
W
B
W
W
5 2
WW
WW
BW
BW
WB
3 5
WBWBB
WWWBW
BWWBB
5 2
BB
BW
WW
WW
BW
5 2
BW
BW
WB
WW
WW
3 1
W
W
B
2 2
BB
WB
1 1
B
4 1
W
B
W
B
5 2
WW
BB
BW
BW
WW
3 4
WBBW
WWWB
WBWB
5 4
BBBW
BBBW
WBWB
BWWW
BWBB
1 3
BBW
3 2
WW
BB
BB
5 4
BBWB
WBBW
BWBW
WWBW
WWBB
2 3
BWW
BWB
4 5
WWBBB
BBBBB
WWWWW
WBWBW
2 5
BWBWB
WBBWB
5 5
WBBBW
WWBWW
BBWBB
BWBBB
BWBBW
2 2
BB
WW
5 3
WWW
BBW
BBB
WWW
WWW
1 1
W
5 4
BWBW
WWBB
WBBW
BWBW
BWBW
3 1
W
B
W
5 3
BWB
BBB
BBW
BWB
BBB
3 3
WWB
BWB
BBB
2 3
WBW
BBB
5 1
W
B
B
W
W
5 2
BB
WB
WW
WW
WB
4 3
BWW
BWB
BWB
BBW
3 4
WBBB
WBBW
WBBW
4 4
BWWW
BBBW
BBWW
WWBW
5 1
B
W
B
B
W
4 2
BW
WW
WB
BB
3 1
W
B
B
5 1
W
B
W
W
W
2 3
WBW
BBW
5 5
WBWWB
WBBWB
BBWBB
WBBWW
WWBWW
4 1
B
W
B
W
1 1
B
5 1
W
B
B
W
W
4 5
WWBBW
WWWWW
WBWWB
BWBWW
3 2
WB
WB
WW
4 3
BWB
WBW
BWW
BBB
3 3
WWB
WWW
BWB
3 5
BWBWB
WBBBB
WWWBB
2 3
WBB
WWB
2 3
WWB
WWB
4 3
WBW
BBW
WWB
WBB
2 5
WBWWB
BBWWW
4 2
BB
WB
BB
WW
4 1
B
B
B
B
2 1
W
W
3 3
BBB
BWB
WWW
2 2
BW
WW
2 4
WWWB
BBBB
1 4
WWWB
4 2
BB
BW
BB
WW
2 1
W
W
1 1
B
5 2
BW
BW
BB
BW
BW
2 4
WWWW
BBBB
3 4
BWWB
WBWB
WBBB
1 2
WB
1 1
W
5 3
BWB
BBW
BWW
BWB
WWB
3 1
W
B
W
3 2
WW
BW
BW
2 5
BWWWW
BBWBW
2 2
WW
WB
4 3
BBB
WWW
WWW
WWB
5 1
B
W
W
W
B
4 3
WWW
WBB
BBW
WBW
3 3
WWW
WBB
BWB
5 2
BB
WB
BW
BB
WB
5 5
BBBBW
BWWWB
BWBWW
WBWBW
WBBWW
1 4
WBWW
1 2
WB
5 1
W
W
W
W
W
5 5
WWWBW
BBWBW
BWWWW
BBBBW
WWWWW
3 4
WBWW
BWBB
BWWB
5 3
WWW
WWW
WWB
WWW
BBB
2 2
BB
BW
2 4
BWBB
BWBB
5 4
WBBW
BWWB
WBBB
WBWW
WWBW
1 2
WB
2 5
WWWWB
WWBBW
3 1
W
W
B
2 5
WWBWB
BWBWB
5 5
BWWWW
WWWBW
WWBWB
BWBBW
BBWWW
3 1
B
B
W
2 1
B
B
1 2
BW
4 3
BBW
BBB
WWW
WBB
3 3
BBW
BBW
BBB
3 3
BWB
WBB
BWW`

	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			grid[i] = scan.Text()
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			sb.WriteString(grid[i])
			sb.WriteByte('\n')
		}
		input := sb.String()
		exp := fmt.Sprintf("%d\n", expected(grid, n, m))
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
