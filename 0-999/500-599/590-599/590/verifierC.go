package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", oracle, "590C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

const testcasesCRaw = `100
1 1
3
5 2
#1
#2
..
#3
#.
5 3
113
.3.
.#2
#22
212
3 2
21
#3
##
4 5
3#33.
2..#2
.3.##
3..1#
2 3
2#3
.13
2 4
#3#1
3121
1 5
13#21
5 2
32
21
.1
13
32
2 1
3
1
3 3
322
#1.
#12
1 3
#23
1 3
.32
2 1
3
1
4 5
32333
#.1#2
13122
21.2#
1 2
31
5 5
33.3#
121..
21#12
11122
121#.
4 3
#.2
3..
#1#
#1.
4 3
1#1
#33
231
.11
1 4
2.3.
1 3
312
4 2
1#
3.
.2
3.
1 3
312
1 1
3
2 3
2.#
.13
3 5
.1##1
132#2
.1112
3 4
23.2
#312
#.13
5 2
#3
22
.2
.1
#2
1 5
.23#1
4 3
..3
#31
2##
2..
1 3
123
5 5
.21.#
23##2
.#21#
..3##
2.213
1 4
#1#3
1 1
3
2 4
..32
13#.
1 3
23.
2 4
#313
2333
5 2
.#
#2
#1
31
2.
5 2
#3
11
1.
32
33
1 3
.32
2 4
2332
12.2
3 2
32
22
31
4 3
3##
#3.
3##
123
2 3
3..
123
4 2
11
32
31
.1
5 3
2#.
#3.
233
2.1
223
3 2
.3
31
32
2 2
21
33
2 1
.
3
2 1
2
3
2 4
.3.1
1.#2
1 4
3.23
2 2
12
.3
1 4
23.2
5 5
112#1
#3#22
3.12#
121.2
1#2.3
5 4
##2#
1222
.233
.2.2
.331
5 1
.
3
1
2
.
5 1
1
3
#
2
2
5 1
2
1
3
.
#
1 2
3#
1 4
1223
4 4
1.31
312.
33#.
##2.
2 2
.2
31
4 2
.#
33
12
3#
4 2
##
23
2#
31
1 1
3
1 4
.23#
2 2
#3
12
4 2
2#
32
2#
31
2 2
#1
32
2 5
..1..
...32
2 2
2#
#3
4 2
#2
13
#3
#.
2 3
13.
122
2 5
1#231
##12.
2 1
2
3
5 3
1#3
33.
3##
211
3.1
3 5
11###
...22
2#31#
2 3
132
#2.
1 5
322#2
4 1
1
3
#
2
4 2
2.
2.
31
#2
5 4
.#.2
.#23
2132
2.#3
...2
4 4
3231
.#1.
31#.
#332
4 4
3##.
2##2
.33.
1#2.
1 4
1.32
2 2
22
13
3 5
.3331
.#12.
11133
5 3
.23
1##
..#
2.1
#13
1 3
.31
4 4
#1.#
22#3
233.
##..
3 1
3
2
#
5 5
22#.2
.3#2#
21..3
22#2.
#1222
1 5
1.223
4 5
11313
1#3..
11.22
.132.
5 1
#
3
3
1
2
1 4
21#3
1 3
.13
3 3
122
213
131
3 1
#
3
2
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	reader := bufio.NewReader(strings.NewReader(testcasesCRaw))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		fmt.Fprintln(os.Stderr, "invalid test file")
		os.Exit(1)
	}
	for caseIdx := 1; caseIdx <= t; caseIdx++ {
		var n, m int
		if _, err := fmt.Fscan(reader, &n, &m); err != nil {
			fmt.Fprintf(os.Stderr, "bad test file at case %d: %v\n", caseIdx, err)
			os.Exit(1)
		}
		rows := make([]string, n)
		for i := 0; i < n; i++ {
			if _, err := fmt.Fscan(reader, &rows[i]); err != nil {
				fmt.Fprintf(os.Stderr, "bad grid at case %d: %v\n", caseIdx, err)
				os.Exit(1)
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			sb.WriteString(rows[i])
			sb.WriteByte('\n')
		}
		input := sb.String()
		exp, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on case %d: %v\n", caseIdx, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", caseIdx, err, got)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed:\nexpected: %s\n got: %s\n", caseIdx, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
