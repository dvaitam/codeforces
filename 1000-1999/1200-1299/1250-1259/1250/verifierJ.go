package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// maxRows is copied from 1250J.go.
func maxRows(counts []int64, L int64) int64 {
	if L == 0 {
		return 0
	}
	var rows, rem int64
	for _, x := range counts {
		if rem > 0 {
			need := L - rem
			if x >= need {
				rows++
				x -= need
				rem = 0
			} else {
				rem = 0
			}
		}
		rows += x / L
		rem = x % L
	}
	return rows
}

// solveAll runs the 1250J logic on the provided input and returns all outputs.
func solveAll(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return "", err
	}
	var out strings.Builder
	for ; t > 0; t-- {
		var n int
		var k int64
		if _, err := fmt.Fscan(reader, &n, &k); err != nil {
			return "", err
		}
		counts := make([]int64, n)
		var sum int64
		for i := 0; i < n; i++ {
			if _, err := fmt.Fscan(reader, &counts[i]); err != nil {
				return "", err
			}
			sum += counts[i]
		}
		if sum < k {
			fmt.Fprintln(&out, 0)
			continue
		}
		lo := int64(0)
		hi := sum/k + 1
		for hi-lo > 1 {
			mid := (lo + hi) / 2
			if maxRows(counts, mid) >= k {
				lo = mid
			} else {
				hi = mid
			}
		}
		fmt.Fprintln(&out, lo*k)
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

// Embedded copy of testcasesJ.txt so the verifier is self-contained.
const testcasesRaw = `100
3 18
10 11 19
3 1
1 12 4
4 12
20 16 16 12
5 5
6 1 13 19 8
1 5
7
3 9
15 20 9
2 10
2 17
2 8
20 12
3 19
7 16 4
4 9
7 13 14 2
4 2
13 3 15 7
3 1
5 15 17
4 14
17 2 5 5
2 16
3 2
1 18
1
1 18
11
1 12
2
2 16
18 3
5 7
10 10 1 8 16
2 8
12 2
1 15
11
5 13
2 2 17 1 6
5 3
8 16 18 18 17
2 2
8 15
1 5
2
2 15
8 16
1 11
18
4 19
15 19 8 8
1 9
13
4 13
6 10 10 20
3 1
7 8 8
4 1
3 13 11 3
4 1
12 15 16 2
3 19
9 7 6
5 18
9 9 17 3 4
4 13
17 3 12 20
3 9
14 15 15
4 14
14 9 4 5
5 17
18 20 19 14 1
3 6
13 17 10
1 11
1
1 1
8
3 7
3 7 2
2 20
19 6
3 18
18 9 14
1 16
5
4 20
1 12 6 19
4 11
1 1 9 7
5 19
18 6 13 20 16
5 16
7 8 12 16 19
3 18
14 17 10
4 9
10 16 20 4
1 1
9
4 17
13 7 2 13
5 17
1 11 16 8 12
5 11
19 20 2 2 8
4 14
16 9 10 20
1 1
9
1 4
8
5 19
11 8 6 16 14
1 12
2
1 5
4
1 13
16
4 6
11 16 1 18
2 10
11 19
2 6
13 17
5 10
9 4 13 17 12
1 9
12
2 19
14 12
1 1
16
1 8
3
3 4
19 9 14
1 9
4
5 13
7 18 12 7 5
4 3
18 13 4 17
2 17
4 15
2 6
11 6
5 1
10 14 14 2 7
5 3
16 6 8 15 15
1 13
4
3 7
19 13 7
1 4
4
4 7
10 1 4 12
5 8
19 1 13 2 5
1 5
3
4 12
7 20 6 9
1 14
11
3 11
19 18 4
2 7
13 14
2 20
3 8
2 10
14 10
4 2
2 13 13 9
1 13
4
1 4
19
1 10
2
1 13
18
5 18
5 16 4 3 19
1 14
7
2 8
1 3
1 7
7`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierJ.go /path/to/binary")
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
