package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `100
8 20
4 13 16 5 3 3 1 13
18 19
2 8 17 18 12 9 6 4 9 7 1 9 9 7 6 10 10 12
3 39
11 13 17
8 12
8 16 9 3 18 10 1 10
19 46
10 17 7 14 14 20 10 14 15 6 8 10 9 2 3 2 15 9 17
18 42
16 11 5 7 3 14 7 15 9 6 12 14 19 11 18 7 11 4
2 46
8 9
19 40
8 4 11 6 10 15 1 2 12 3 10 11 1 11 10 11 5 14 20
3 19
20 7 15
10 9
9 13 20 6 11 19 1 12 2 15
6 24
12 10 19 4 15 7
14 14
4 2 2 2 6 20 5 20 2 18 16 19 8 11
2 8
17 10
14 42
7 16 7 8 15 14 16 2 8 14 15 8 14 7
16 13
2 2 9 9 8 17 7 8 14 9 5 11 2 11 19 4
19 26
2 16 13 3 14 7 19 6 11 10 16 11 14 17 7 9 11 13 16
3 18
7 2 13
20 9
9 2 6 15 19 16 13 13 7 1 7 6 1 20 9 4 13 13 8 18
2 13
6 20
11 36
16 17 15 1 3 2 20 4 16 18 9
20 50
5 2 12 3 17 1 10 12 3 3 18 15 13 7 10 13 8 16 13 4
3 8
20 12 17
14 27
15 3 7 10 16 14 4 18 6 12 6 6 5 11
16 22
9 18 1 6 1 10 4 18 4 16 20 16 17 3 17 8
14 19
12 8 6 1 2 20 11 18 15 19 10 17 15 20
20 29
13 5 9 20 12 11 5 14 3 20 5 20 6 10 12 7 19 12 20 3
3 26
6 11 12
11 12
10 1 20 1 17 3 12 4 6 6 19
16 43
19 3 4 6 16 8 20 10 13 20 8 16 8 10 12 8
11 35
17 15 13 17 13 11 10 15 14 19 1
9 12
18 15 18 20 12 13 13 20 1
5 33
3 15 12 20 10
3 49
9 16 8
16 39
3 5 8 3 10 5 2 6 13 19 20 18 18 9 1 17
8 11
4 7 2 11 3 4 9 2
10 40
6 5 14 5 3 18 12 1 18 5
17 27
5 7 10 16 17 3 13 6 6 9 17 13 18 10 13 11 6
13 4
14 1 9 1 10 5 3 6 4 20 1 8 8
18 1
16 18 6 15 13 11 6 17 19 7 4 16 20 12 10 14 20 2
7 17
19 19 10 16 11 3 8
11 7
2 20 11 17 16 12 3 6 2 16 17
18 39
8 2 7 3 11 20 5 10 4 17 17 1 3 8 20 9 20 1
2 1
16 5
7 24
8 12 10 13 20 13 2
6 27
16 3 17 17 10 18
10 14
14 11 6 7 1 16 14 9 11 14
14 21
20 7 9 9 17 3 1 13 9 20 10 5 20 11
3 11
20 3 7
9 43
7 16 4 11 1 15 19 7 6
3 16
17 4 4
11 19
19 19 5 3 6 1 6 7 9 2 19
11 41
11 20 12 6 6 7 2 2 6 2 9
16 7
7 5 13 16 20 18 16 18 11 12 14 13 2 2 6 5
12 15
10 8 12 19 7 20 12 13 12 18 4 8
7 23
6 18 7 1 1 4 7
10 45
11 18 16 2 10 13 6 5 13 20
15 23
12 1 15 5 15 1 14 7 4 11 10 20 18 2 3
16 11
17 2 6 20 8 20 11 20 1 1 13 14 8 19 5 20
2 9
6 20
10 15
14 11 18 12 6 12 18 19 12 9
20 22
18 2 18 5 12 19 7 16 16 4 3 20 7 6 5 18 5 20 2 18
12 6
2 20 5 18 6 9 20 2 10 16 9 7
16 49
10 17 2 5 10 14 2 1 12 17 20 12 3 7 14 3
14 3
10 7 7 6 7 3 3 1 1 10 17 17 17 14
3 6
17 11 6
7 33
18 2 2 3 5 14 5
4 49
17 11 17 11
9 9
19 19 4 19 20 3 20 1 6
2 50
9 18
14 50
10 19 12 5 17 12 18 15 8 20 15 6 16 19
18 2
19 18 19 7 9 19 19 20 15 11 6 10 6 11 10 11 7 18
16 40
10 9 12 6 11 4 8 11 7 6 15 4 10 2 17 1
2 1
18 6
8 16
12 13 13 8 15 14 10 17
3 22
2 4 6
2 27
1 13
6 28
2 18 8 12 11 17
12 44
5 1 12 12 4 13 19 7 16 12 19 14
18 47
3 13 6 3 15 13 17 7 1 3 10 4 5 5 15 15 12 14
11 26
10 13 4 11 13 17 3 16 14 12 8
19 40
6 14 10 4 11 18 2 17 8 20 6 16 8 17 12 14 15 14 18
4 32
1 14 12 16
9 13
5 11 12 5 9 18 19 16 20
5 32
13 16 12 20 17
2 48
14 9
19 40
6 18 6 1 15 19 10 5 2 1 8 11 5 11 11 15 19 8 9
2 43
3 2
15 8
8 8 16 18 3 18 20 6 6 7 1 7 17 7 8
1 33
16
11 23
9 19 13 19 5 11 4 20 9 9 4
2 22
2 13
3 14
8 2 16
13 24
3 3 5 15 14 17 13 4 6 5 2 6 15
18 34
16 15 12 19 14 8 9 9 19 1 9 13 19 11 10 4 6 17
7 11
13 18 5 6 20 7 18
3 35
3 10 1
3 33
15 4 17
1 18
9
7 50
7 4 9 19 13 9 17
9 1
4 13 12 20 6 16 20 11 7`

func expected(n int, x int64, a []int64) string {
	var maxA int64
	for _, v := range a {
		if v > maxA {
			maxA = v
		}
	}
	lo := int64(1)
	hi := maxA + x + 1
	for lo+1 < hi {
		mid := (lo + hi) / 2
		var need int64
		for _, v := range a {
			if mid > v {
				need += mid - v
				if need > x {
					break
				}
			}
		}
		if need <= x {
			lo = mid
		} else {
			hi = mid
		}
	}
	return fmt.Sprintf("%d\n", lo)
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

func loadCases() ([]string, []string) {
	tokens := strings.Fields(testcasesRaw)
	if len(tokens) == 0 {
		fmt.Println("no embedded testcases")
		os.Exit(1)
	}
	pos := 0
	t, err := strconv.Atoi(tokens[pos])
	if err != nil {
		fmt.Println("invalid testcase count")
		os.Exit(1)
	}
	pos++
	var inputs []string
	var expects []string
	for caseIdx := 1; caseIdx <= t; caseIdx++ {
		if pos+1 >= len(tokens) {
			fmt.Printf("case %d missing header\n", caseIdx)
			os.Exit(1)
		}
		n, errN := strconv.Atoi(tokens[pos])
		x, errX := strconv.ParseInt(tokens[pos+1], 10, 64)
		if errN != nil || errX != nil {
			fmt.Printf("case %d invalid header\n", caseIdx)
			os.Exit(1)
		}
		pos += 2
		if pos+n > len(tokens) {
			fmt.Printf("case %d missing numbers\n", caseIdx)
			os.Exit(1)
		}
		a := make([]int64, n)
		var arrStr strings.Builder
		for i := 0; i < n; i++ {
			val, errV := strconv.ParseInt(tokens[pos+i], 10, 64)
			if errV != nil {
				fmt.Printf("case %d invalid value\n", caseIdx)
				os.Exit(1)
			}
			a[i] = val
			if i > 0 {
				arrStr.WriteByte(' ')
			}
			arrStr.WriteString(tokens[pos+i])
		}
		pos += n
		inputs = append(inputs, fmt.Sprintf("1\n%d %d\n%s\n", n, x, arrStr.String()))
		expects = append(expects, expected(n, x, a))
	}
	return inputs, expects
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	inputs, expects := loadCases()
	for idx, input := range inputs {
		if err := runCase(exe, input, expects[idx]); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}
