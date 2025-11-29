package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getChar(x int) byte {
	if x < 26 {
		return byte('a' + x)
	} else if x < 52 {
		return byte('A' + (x - 26))
	}
	return byte('0' + (x - 52))
}

// solveAll implements the logic from 1254A.go for all testcases in a single input string.
func solveAll(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return "", err
	}
	var out strings.Builder
	for ; T > 0; T-- {
		var n, m, k int
		if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
			return "", err
		}
		grid := make([][]byte, n)
		rCount := 0
		for i := 0; i < n; i++ {
			var s string
			if _, err := fmt.Fscan(reader, &s); err != nil {
				return "", err
			}
			row := []byte(s)
			for j := 0; j < m; j++ {
				if row[j] == 'R' {
					rCount++
				}
			}
			grid[i] = row
		}
		counts := make([]int, k)
		for i := 0; i < rCount; i++ {
			counts[i%k]++
		}
		res := make([][]byte, n)
		for i := range res {
			res[i] = make([]byte, m)
		}
		chi := 0
		for i := 0; i < n; i++ {
			if i%2 == 0 {
				for j := 0; j < m; j++ {
					for chi < k-1 && counts[chi] == 0 {
						chi++
					}
					if grid[i][j] == 'R' {
						counts[chi]--
					}
					res[i][j] = getChar(chi)
				}
			} else {
				for jj := m - 1; jj >= 0; jj-- {
					for chi < k-1 && counts[chi] == 0 {
						chi++
					}
					if grid[i][jj] == 'R' {
						counts[chi]--
					}
					res[i][jj] = getChar(chi)
				}
			}
		}
		for i := 0; i < n; i++ {
			out.Write(res[i])
			out.WriteByte('\n')
		}
	}
	return strings.TrimSpace(out.String()), nil
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

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesRaw = `100
4 4 1
....
..RR
.RR.
R.RR
3 4 5
R...
R...
RRR.
1 4 3
R.RR
2 2 2
.R
R.
4 1 3
.
R
.
R
3 4 1
..R.
RRRR
..RR
2 2 1
R.
.R
2 4 5
....
R.R.
3 2 2
R.
RR
.R
3 4 1
RRRR
RRRR
R.R.
1 1 1
R
2 1 2
R
R
1 4 1
.RRR
3 3 4
RR.
RR.
R..
4 2 2
RR
R.
.R
.R
1 4 4
....
2 1 2
R
.
1 3 1
R..
3 3 5
R..
.RR
R.R
2 2 4
..
R.
4 1 2
.
R
.
R
4 4 5
RR..
..R.
RRRR
...R
2 1 1
R
R
1 2 2
.R
1 4 4
RR..
1 3 1
.RR
3 1 1
R
R
.
3 3 5
R..
..R
R..
1 2 1
..
3 3 1
.RR
.RR
...
2 3 1
.RR
.RR
3 4 3
..RR
....
R.R.
4 1 3
.
R
R
.
1 1 1
R
2 1 2
R
R
4 3 4
.R.
R.R
.R.
R.R
1 1 1
R
1 4 4
.RRR
2 1 1
.
.
3 2 1
.R
..
..
4 4 2
RR..
.RR.
R...
RRRR
2 2 1
.R
..
2 2 4
..
RR
1 4 2
....
2 1 1
R
.
3 3 1
...
...
RRR
2 3 3
RR.
.R.
1 2 2
.R
4 4 4
RR.R
RRR.
R.R.
.RR.
3 1 3
.
R
.
3 1 3
R
R
.
4 3 2
..R
RR.
.RR
RR.
3 1 1
.
.
R
1 4 1
R.RR
4 4 4
R...
R.RR
.R..
RRR.
1 2 1
RR
1 1 1
.
1 2 1
R.
1 4 3
.RRR
3 3 4
...
RRR
.R.
2 2 2
.R
RR
3 3 4
RRR
.RR
R..
4 2 5
..
.R
R.
.R
4 4 1
...R
...R
RR.R
R..R
4 4 2
R..R
..R.
R.R.
.RR.
1 3 3
R..
1 4 2
..RR
3 1 1
.
.
.
4 1 3
.
.
R
R
1 1 1
.
3 1 3
R
R
R
4 3 3
RRR
...
...
RRR
1 4 4
.R..
4 4 2
....
R..R
.RRR
...R
1 1 1
R
3 1 1
R
.
.
2 2 3
RR
..
3 4 2
...R
.R.R
RR..
4 4 4
.RR.
RR.R
RRR.
RRR.
1 3 2
..R
4 4 4
..R.
.RRR
R.R.
RRR.
2 4 5
R.R.
RR..
4 3 3
RRR
.RR
.RR
...
4 3 4
.RR
.RR
...
R..
4 4 1
RRRR
R.R.
R..R
RR.R
2 3 3
..R
.R.
4 4 5
..RR
RRR.
RR.R
.RR.
3 1 3
.
R
.
2 2 1
..
R.
4 3 5
RRR
..R
R.R
RRR
2 4 5
RRRR
..R.
1 1 1
.
4 3 1
RR.
..R
.R.
R..
4 1 2
R
.
R
.
4 2 3
.R
.R
.R
R.
1 4 1
RRRR
2 3 1
R..
RR.
4 2 3
.R
..
.R
..
4 1 1
.
.
R
R
2 2 1
.R
RR`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}

	input := strings.TrimSpace(testcasesRaw) + "\n"
	expect, err := solveAll(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "solver failed: %v\n", err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	got, err := runCandidate(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if strings.TrimSpace(got) != strings.TrimSpace(expect) {
		fmt.Println("output mismatch")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
