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

// Embedded testcases from testcasesA.txt.
const embeddedTestcasesA = `100
5 10 2 5 2 8
15 8 7 4 2 8 1 7 7 10 1 8 5 4 10 2
11 1 1 1 9 1 7 4 7 1 9 4
15 8 9 4 6 4 4 8 5 1 7 9 2 3 5 2
11 9 7 9 4 5 5 10 8 9 7 10
2 8 4
13 7 3 6 9 6 2 8 9 2 3 9 7 6
16 1 8 1 5 10 10 10 7 3 3 9 4 1 4 9 9
8 7 9 6 10 6 8 5 9
20 1 7 9 3 9 9 4 7 1 8 6 10 9 4 9 7 8 6 7 6
1 9
18 10 10 6 8 10 1 4 3 9 10 3 2 9 5 1 2 2 1
15 1 5 4 5 2 10 3 6 5 2 3 3 5 9 3
9 5 8 6 8 8 2 1 5 7
11 7 4 5 2 5 9 4 10 7 1 4
1 7
5 1 3 8 9 7
18 4 9 8 4 9 1 7 10 6 7 1 5 3 4 1 5 2 2
10 5 3 7 10 5 3 1 9 1 10
7 10 8 3 10 9 1 7
7 6 2 4 10 7 10 4
16 2 7 5 9 8 1 6 10 7 5 1 3 4 6 10 3
11 7 4 5 2 7 9 6 9 8 9 4
3 1 2 3
6 3 9 4 5 6 10
17 5 6 6 6 2 5 4 10 8 3 10 9 2 6 1 7 2
13 3 3 6 2 10 10 7 2 10 9 4 10 2
9 6 5 10 9 2 8 5 2 1
10 1 10 1 2 7 2 1 4 4 10
14 3 2 8 3 4 3 2 7 7 9 5 9 5 8
11 2 4 6 1 1 1 5 10 6 8 7
11 7 2 2 6 10 8 2 5 4 10 9
16 6 5 3 9 4 5 4 4 6 2 5 2 8 2 10 6
8 7 5 1 6 3 6 10 5
8 6 2 9 10 10 10 2 4
8 1 4 7 2 5 9 2 2
1 1
10 6 8 8 3 2 9 6 2 9 3
6 3 3 6 5 2 9
20 5 3 4 3 9 1 6 10 9 4 3 5 7 9 3 1 4 5 2 8
14 9 5 9 8 9 8 1 7 6 3 5 8 1 7
19 1 1 6 10 3 10 3 3 5 5 7 10 7 3 10 2 4 8 1
6 9 6 9 8 4 4
11 8 8 4 7 6 9 10 5 4 1 2
17 6 3 9 4 5 5 5 9 6 3 8 10 2 2 10 9 10
13 3 3 5 7 4 10 1 8 7 6 7 9 3
18 1 9 2 5 2 5 2 3 10 2 8 4 7 7 7 3 6 8
5 10 8 4 2 7
20 9 7 2 5 5 4 7 9 1 4 9 8 10 1 1 10 4 5 4 3
10 3 9 4 5 5 10 5 8 3 9
12 8 7 2 4 10 7 4 5 2 1 2 10
1 9
10 3 2 9 6 10 5 7 9 6 9
11 1 2 8 8 6 5 9 7 6 10 8
4 7 7 4 9
1 5
20 9 4 8 10 9 7 5 3 8 10 9 4 6 9 1 7 10 7 7 6
20 10 2 8 4 5 1 7 3 7 5 3 2 10 1 6 5 7 9 5 3
15 5 8 3 8 9 1 5 9 2 10 7 2 6 2 8
1 3
17 3 2 7 5 10 5 4 9 4 4 6 5 2 2 9 6 8
17 9 1 3 5 9 5 6 10 4 7 9 7 3 8 5 10 6
8 5 10 4 1 10 7 6 7
8 5 4 2 3 10 8 10 3
20 5 8 9 3 3 3 8 6 5 7 4 2 4 5 2 2 4 7 6 8
4 3 1 1 10
1 4
2 8 9
20 8 6 5 2 10 3 2 4 7 4 8 8 7 3 4 4 5 8 9 10
13 4 8 5 6 8 10 2 4 2 1 1 1 8
11 7 10 5 4 7 3 3 1 1 7 3
18 1 10 7 5 3 2 8 5 1 1 9 1 9 3 1 5 2 7
3 4 1 8
5 5 4 8 7 6
9 5 4 4 1 10 10 3 6 7
20 9 9 1 6 9 7 9 4 9 7 2 5 10 2 5 3 2 3 1 4
14 1 1 2 9 8 9 6 2 6 1 3 9 1 8
5 7 8 1 9 5
3 5 6 2
10 1 7 1 5 6 3 5 7 2 5
4 7 4 9 9
7 6 6 9 7 10 8 2
5 8 9 9 10 9
18 1 5 3 4 6 7 9 6 2 7 6 3 10 2 1 5 9 6
14 5 6 6 5 6 9 9 1 9 2 3 6 6 6
19 2 8 5 8 8 6 7 2 10 1 3 1 9 8 10 5 4 10 6
12 6 7 5 8 10 6 9 9 3 1 3 5
8 10 3 2 3 7 10 1 2
18 5 2 4 5 2 10 9 2 2 4 3 9 7 1 10 6 8 5
8 4 10 8 4 7 8 6 9
7 8 2 5 7 4 1 9
13 9 8 2 7 10 9 10 10 7 1 6 8 1
7 5 1 9 2 5 9 6
18 10 9 5 9 7 9 9 7 10 10 5 8 5 3 9 8 10 3
18 3 5 1 7 10 1 6 7 7 5 1 2 2 1 7 5 8 5
12 8 6 7 8 2 8 6 3 7 3 1 3
9 6 3 10 5 7 5 9 5 7
9 7 6 8 4 8 7 7 2 2
5 4 3 4 1 2
9 3 8 2 7 3 1 2 7 10`

func solve637A(nums []int) int {
	counts := make(map[int]int)
	bestID := 0
	bestCount := 0
	for _, id := range nums {
		counts[id]++
		if counts[id] > bestCount {
			bestCount = counts[id]
			bestID = id
		}
	}
	return bestID
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(embeddedTestcasesA))
	scanner.Split(bufio.ScanWords)
	nextInt := func() int {
		if !scanner.Scan() {
			fmt.Fprintln(os.Stderr, "unexpected EOF in test data")
			os.Exit(1)
		}
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
			os.Exit(1)
		}
		return v
	}

	t := nextInt()
	for caseIdx := 1; caseIdx <= t; caseIdx++ {
		n := nextInt()
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			nums[i] = nextInt()
		}
		want := strconv.Itoa(solve637A(nums))

		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", n)
		for i, v := range nums {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseIdx, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", caseIdx, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
