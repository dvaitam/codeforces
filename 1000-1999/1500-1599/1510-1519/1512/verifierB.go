package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesRaw = `3 .*. ... *..
4 ..*. .... ...* ....
2 *. *.
2 .* *.
5 ....* ..... ..... *.... .....
2 *. *.
3 ..* ... ..*
3 ... .** ...
2 *. *.
3 *.* ... ...
4 ..*. ..*. .... ....
4 .... .... *... .*..
3 *.. *.. ...
4 ..*. .... ..*. ....
5 ..... ..... ...*. ..*.. .....
4 .... .... .... .**.
5 *.... ..... ..*.. ..... .....
4 .... **.. .... ....
5 *.... ..... ..... ..... *....
5 *.... ..... ..... ...*. .....
2 .. **
3 ..* ... ..*
4 .... ..*. *... ....
5 ..... *.... ..... ....* .....
2 .* .*
2 .* *.
4 .... .... .... ..**
3 ... *.* ...
5 ..... ..... ..... *.... .*...
3 ... ..* *..
4 *... ...* .... ....
2 *. *.
3 ... .*. ..*
4 ..*. .... ...* ....
5 ..*.. ..... ..... ....* .....
2 *. .*
3 ... .*. *..
3 *.. ... .*.
5 ..... ..... ..... ..... *.*..
4 .... .... *... *...
3 *.. ... ..*
2 .. **
2 .. **
2 .* *.
2 *. .*
5 ..... ..... ..**. ..... .....
3 ... ..* .*.
3 ... *.. .*.
2 .. **
5 ...*. ..... ..... ..... .*...
4 ...* .... *... ....
5 ..... ..... ..... ..*.. ...*.
2 .* *.
5 *.... ..... ..... ...*. .....
2 .* .*
3 ..* .*. ...
3 ... *.. ..*
4 .... .... .*.. ...*
4 .... .... .... *.*.
3 ... ..* ..*
5 ...*. ..... ...*. ..... .....
3 ... *.. *..
4 ...* .... .... ..*.
2 *. *.
2 .. **
3 .*. ..* ...
4 .... .... *... .*..
4 *... .... ..*. ....
2 .. **
5 ..... ....* .*... ..... .....
5 ..... ..... *.... ..... ...*.
5 ...*. ..... ..... ..... *....
5 ..... .*... ..*.. ..... .....
2 .. **
3 ... ... **.
4 .... .*.* .... ....
3 *.. ..* ...
5 ..... ..... ...*. ...*. .....
5 ..... ..... *.... ....* .....
2 .. **
3 *.. .*. ...
5 ..... ..... ..... ..... *..*.
3 ... .*. *..
3 **. ... ...
4 .... .... *... *...
5 ..... ..... ..... ..*.. .*...
4 .... *..* .... ....
4 ..*. .... .... ...*
5 .*... ..... ..... .*... .....
4 ...* .... ...* ....
5 ..*.. ..... .*... ..... .....
2 *. .*
3 ..* ... ..*
2 .. **
4 ...* .... .*.. ....
3 ... *.* ...
4 .**. .... .... ....
4 .... .... *... ..*.
2 .* .*
4 .... .... .*.* ....`

type testCase struct {
	n    int
	grid []string
}

func solveCase(tc testCase) []string {
	g := make([][]byte, tc.n)
	stars := [][2]int{}
	for i := 0; i < tc.n; i++ {
		g[i] = []byte(tc.grid[i])
		for j := 0; j < tc.n; j++ {
			if g[i][j] == '*' {
				stars = append(stars, [2]int{i, j})
			}
		}
	}
	r1, c1 := stars[0][0], stars[0][1]
	r2, c2 := stars[1][0], stars[1][1]
	if r1 == r2 {
		r3 := r1 + 1
		if r3 >= tc.n {
			r3 = r1 - 1
		}
		g[r3][c1] = '*'
		g[r3][c2] = '*'
	} else if c1 == c2 {
		c3 := c1 + 1
		if c3 >= tc.n {
			c3 = c1 - 1
		}
		g[r1][c3] = '*'
		g[r2][c3] = '*'
	} else {
		g[r1][c2] = '*'
		g[r2][c1] = '*'
	}
	res := make([]string, tc.n)
	for i := 0; i < tc.n; i++ {
		res[i] = string(g[i])
	}
	return res
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	res := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("case %d: malformed line", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", idx+1, err)
		}
		if len(fields)-1 != n {
			return nil, fmt.Errorf("case %d: expected %d rows got %d", idx+1, n, len(fields)-1)
		}
		grid := make([]string, n)
		copy(grid, fields[1:])
		res = append(res, testCase{n: n, grid: grid})
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	return res, nil
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
		expectedLines := solveCase(tc)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for _, row := range tc.grid {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotLines := strings.Split(strings.TrimSpace(got), "\n")
		if len(gotLines) != tc.n {
			fmt.Printf("case %d failed\nexpected %d lines\ngot: %d lines\n", i+1, tc.n, len(gotLines))
			os.Exit(1)
		}
		match := true
		for idx, line := range expectedLines {
			if strings.TrimSpace(gotLines[idx]) != line {
				match = false
				break
			}
		}
		if !match {
			fmt.Printf("case %d failed\nexpected:\n%s\ngot:\n%s\n", i+1, strings.Join(expectedLines, "\n"), strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
