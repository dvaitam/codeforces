package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `4 3
2 8
4 1
1 9
7 2
5 2
4 2
5 6
7 3
1 9
8 1
2 7
4 5
6 8
3 4
1 3
3 6
9 5
2 8
3 1
8 7
9 5
6 7
5 3
9 1
8 2
6 1
9 5
3 4
8 6
5 6
3 5
7 7
2 1
4 6
3 4
4 8
7 7
1 7
7 1
3 8
2 5
3 8
9 8
9 1
1 8
6 5
8 1
7 4
9 2
3 1
7 7
6 1
4 1
1 9
2 4
2 4
5 5
3 2
8 7
2 1
5 8
2 5
3 9
6 2
3 5
1 1
1 4
5 9
6 6
1 8
8 7
6 9
3 4
7 5
1 3
3 5
6 6
6 2
6 1
1 5
3 3
5 6
7 9
3 5
2 8
4 1
5 3
9 2
5 7
6 5
7 2
2 9
8 8
6 6
2 8
2 8
7 1
5 6
3 3`

var (
	n, m   int
	ans    int
	rem    int
	p      [9][9]byte
	ansp   [9][9]byte
	shapes = [3][12]byte{
		{'A', 'A', 'A', '.', 'A', '.', 'A', '.', '.', '.', '.', 'A'},
		{'.', 'A', '.', '.', 'A', '.', 'A', 'A', 'A', 'A', 'A', 'A'},
		{'.', 'A', '.', 'A', 'A', 'A', 'A', '.', '.', '.', '.', 'A'},
	}
)

func copyAns() {
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			ansp[i][j] = p[i][j]
		}
	}
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func saya(x, y, move int) {
	if x >= n-2 {
		if move > ans {
			ans = move
			copyAns()
		}
		return
	}
	if y >= m-2 {
		if p[x][y] == '.' {
			rem--
		}
		if p[x][y+1] == '.' {
			rem--
		}
		saya(x+1, 0, move)
		if p[x][y] == '.' {
			rem++
		}
		if p[x][y+1] == '.' {
			rem++
		}
		return
	}
	if rem/5 <= ans-move {
		return
	}
	if p[x][y] == '.' {
		rem--
	}
	for d := 0; d < 12; d += 3 {
		flag := false
		for i := 0; i < 3 && !flag; i++ {
			for j := 0; j < 3; j++ {
				if shapes[i][d+j] == 'A' && p[x+i][y+j] != '.' {
					flag = true
					break
				}
			}
		}
		if !flag {
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					if shapes[i][d+j] == 'A' {
						p[x+i][y+j] = 'A' + byte(move)
					}
				}
			}
			rem -= 5
			saya(x, y+1, move+1)
			rem += 5
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					if shapes[i][d+j] == 'A' {
						p[x+i][y+j] = '.'
					}
				}
			}
		}
	}
	saya(x, y+1, move)
	if p[x][y] == '.' {
		rem++
	}
}

func solveCase(nn, mm int) int {
	n, m = nn, mm
	ans = 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			p[i][j] = '.'
			ansp[i][j] = '.'
		}
	}
	rem = n * m
	saya(0, 0, 0)
	return ans
}

func parseTestcases() ([][2]int, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases [][2]int
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 numbers got %d", idx+1, len(parts))
		}
		n, err1 := strconv.Atoi(parts[0])
		m, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("line %d: parse error: %v %v", idx+1, err1, err2)
		}
		cases = append(cases, [2]int{n, m})
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
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
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for i, tc := range cases {
		input := fmt.Sprintf("%d %d\n", tc[0], tc[1])
		gotStr, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expect := solveCase(tc[0], tc[1])
		var got int
		if _, err := fmt.Sscan(gotStr, &got); err != nil {
			fmt.Printf("test %d: could not parse output %q\n", i+1, gotStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
