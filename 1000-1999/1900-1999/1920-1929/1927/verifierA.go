package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcases = `13 WBWWWWWWBBWBB
20 WBWBBWWBWWWBWWWBBBWB
20 WWBWBBBBBWBBWWBWWBWB
20 WWBWWBWBBBBWWBBBBBBW
17 WBBWWWWWBWBWWBBBW
4 BWBW
14 BBBBBBBBBBWBWB
2 BB
6 BWBBBW
20 BWBBBWWWBBWBBWBWWWBB
2 BB
11 WBWBBWWWWWW
5 BWBWB
18 WBBWWWWBWWWBBBWBBB
15 WWBWWBBWBWBWWBB
16 WWWBWBBBBWWWBBBB
17 BBBBWWBBWWBBWWBWB
17 WBBWBBBBWWWBWWWWB
7 WWBBBWW
11 WBWBBWBBWWW
18 BWBWBBWBBWWWWWBBWW
11 WBWBWWBWWBB
19 WBBBBBBWBBWWWWBWBWB
9 BWBWBBBWB
4 WWWB
2 BB
4 BWWW
18 BBWBWWWWWWBBBWWWBB
9 BWWWBBBBB
5 BWBWW
6 BWWWBB
17 BWBWWWWBBBBWWWBWW
10 WWWBBBBWWB
2 WB
5 WBBWW
2 WB
15 BBWBBBBWBWBWWBB
20 WWBWBWWBBBWWWBWWBBBW
9 BBBBWWBBW
19 WBBWBBWBBWWWBWWWBWB
4 WBWW
6 BBWBBB
1 B
4 BWWB
20 BBBWBWWWBBBWWWWWWBBB
18 WBWBBBWBBBWWWBBBWB
3 BWW
18 WBWWWBBWWBWWBWWWBW
10 WBBBWBBWWB
14 WBBWWBWWBWBWBW
20 WBBWBWBWWBWBWWBBWBBW
13 WWBWWWBBBBWWW
4 BBBW
10 WBBBWWWWWW
19 BBBBWWWBWWWWBWWWWBW
10 BWBBBWWWBB
3 WBW
2 BB
9 WBBWBBWWW
13 BWWWBWBWBBBWW
13 WWWBBWBBWBBBB
12 BBBWBWWWWBWW
16 WWBWWBBBBWBWBBBW
8 WBWBWBBW
11 WWWBBBWBBWB
7 WWWWWWB
8 BWBBWWWB
14 WWWBBBBBBWBWBW
18 WBBBWBBWWWWBWBWWWW
9 BBBBBWBBW
1 B
2 BW
9 BWBWBBBWW
8 WWWBBBWW
19 BBWBBBBBWBBBBWWBWBB
17 WWWWBBBWWWBWBWBWW
13 BBBWBWWBWWBWB
12 BBWBWBBBBBBW
3 BWW
5 BWWBW
10 BWWWBWWWBB
14 WBBBBBWBBBBBBB
14 BWBBWWWBBWBWWB
15 BWBWWWWWBWWWBWB
15 WBWWWBBWBWWBWWW
7 BWBWBBW
1 B
19 BWWWWWWWWBBBBWBWBWW
20 BBWWBBBBBWWWWWBWWWWW
7 WBBBWWW
2 WB
1 B
17 BBWWBBBWBWBWBWBBB
15 BWWWWWWWBWBWBWB
13 WWWBWWBWWBWWB
2 BW
12 BBBBWBWBWWBW
14 BBWWBBBBWBBWBW
19 BBBWBWWWWWBWBBBBWBB
15 BWBWBWBBWBBBWWW`

func runProg(bin, input string) (string, error) {
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

func solve(n int, s string) int {
	first, last := -1, -1
	for i := 0; i < n; i++ {
		if s[i] == 'B' {
			if first == -1 {
				first = i
			}
			last = i
		}
	}
	return last - first + 1
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	lines := strings.Split(testcases, "\n")
	count := 0
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		count++
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "invalid testcase line %d\n", idx+1)
			os.Exit(1)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid n on line %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		s := parts[1]
		want := solve(n, s)

		input := fmt.Sprintf("1\n%d\n%s\n", n, s)
		gotStr, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(gotStr))
		if err != nil {
			fmt.Printf("test %d output parse error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("test %d failed\nexpected: %d\ngot: %d\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", count)
}
