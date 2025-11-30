package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesA = `1 25
1 50
1 100
3 25 50 25
3 25 50 100
4 25 25 50 100
4 25 25 25 100
6 25 25 50 50 25 100
5 25 25 50 100 100
2 25 50
13 50 25 50 100 50 50 50 50 50 100 25 100 25
10 25 25 100 50 100 100 100 25 50 25
3 100 50 50
18 25 50 50 50 100 100 25 100 50 50 100 50 25 100 25 25 100 50
1 100
16 50 25 100 50 100 25 25 100 25 25 25 100 50 25 25 50
17 50 25 50 100 50 100 25 100 50 100 25 100 100 100 50 50 25
20 50 50 100 25 50 25 25 25 25 100 100 50 50 25 25 100 25 25 25 25
18 100 50 100 100 50 100 25 25 100 100 50 100 50 50 50 100 100 100
12 25 50 100 25 50 100 100 50 25 25 25 100
9 25 100 25 50 25 50 50 25 25
5 100 25 25 100 100
18 100 100 25 25 25 100 25 100 100 25 50 25 50 25 25 100 25 25
6 100 25 50 25 100 25
1 100
14 100 25 50 25 25 25 100 50 50 50 25 25 100 50
2 100 25
13 25 50 50 100 50 100 25 100 100 25 25 100 25
6 50 100 50 25 100 50
6 25 50 100 50 100 100
10 100 50 50 100 50 25 100 100 25 50
3 50 100 25
18 50 25 25 50 50 100 50 100 50 100 100 100 25 100 50 50 100 50
3 25 100 25
11 25 25 25 100 50 50 100 100 100 50 25
13 100 100 50 100 100 25 25 50 25 50 100 25 50
17 50 100 100 25 25 50 50 50 50 25 50 25 100 100 25 100 25
1 50
14 50 25 25 25 100 25 100 100 100 25 25 25 100 100
7 50 50 100 25 25 50 50
3 25 50 50
4 50 25 100 100
12 25 25 50 25 25 25 25 100 50 100 50 50
19 25 100 100 100 100 50 100 100 50 100 50 50 100 25 25 50 100 50 25
5 25 50 50 50 50
3 50 100 25
2 50 25
5 100 50 50 50 100
5 50 25 50 100 25
2 50 25
17 100 25 50 50 50 50 50 25 25 100 50 50 50 50 25 50 25
16 50 25 50 50 100 100 25 25 100 100 50 100 25 25 25 25
8 25 50 25 25 50 100 100 50
15 50 100 100 100 25 50 25 50 25 50 100 25 50 25 50
4 25 100 25 100
15 100 25 25 50 50 50 25 100 25 25 100 25 25 25 50
13 50 100 25 25 100 50 25 100 50 100 100 50 100
16 100 50 50 50 100 100 25 100 100 25 25 50 100 100 50 50
2 100 25
9 100 25 50 100 50 100 100 50 25
3 100 25 25
8 25 25 50 25 50 50 25 25
15 50 100 50 100 100 25 100 25 100 100 100 25 100 50 25
10 100 100 50 50 50 100 100 25 25 100
1 100
6 50 100 100 50 50 25
16 50 50 50 50 50 25 25 100 25 25 50 100 50 25 25 50
14 25 50 100 100 25 100 100 25 50 50 25 100 50 50
15 25 25 50 25 25 25 100 100 25 100 50 25 100 100 100
12 25 50 50 25 25 100 100 50 100 100 50 100
4 100 100 100 50
5 50 50 100 100 25
17 25 25 50 50 50 25 50 50 100 25 25 25 50 50 25 100 100
19 100 25 25 25 100 100 100 100 50 100 50 25 25 50 100 50 25 25 50
3 25 50 25
1 50
14 100 50 25 50 50 100 50 25 50 25 25 50 100 25
12 50 25 25 25 50 25 100 25 25 25 25 100
4 100 25 50 50
1 100
8 25 25 50 25 50 50 100 50
5 25 25 50 50 50
10 50 100 100 50 25 100 25 25 100 100
10 25 50 25 25 25 50 100 25 25 25
10 50 50 100 25 25 100 25 50 25 100
3 25 50 50
18 50 100 25 100 50 50 100 50 25 25 50 100 50 100 100 25 50 50
1 50
11 50 25 50 50 50 25 25 25 50 25 100
20 100 25 100 50 50 25 50 50 25 25 50 50 100 25 50 50 100 100 25 50
7 50 25 100 50 100 50 25
7 50 25 100 50 100 100 25
14 100 50 25 100 50 25 100 50 50 100 25 25 50 100
2 25 50
13 100 100 100 50 50 25 100 50 50 50 100 25 100
3 25 25 50
10 100 50 25 100 25 25 25 50 50 50
17 25 100 100 100 25 100 25 50 100 100 50 100 50 50 50 50 100
5 25 25 100 25 50
13 100 50 25 100 100 50 50 100 50 100 50 25 50`

func expectedA(bills []int) string {
	count25, count50 := 0, 0
	for _, b := range bills {
		switch b {
		case 25:
			count25++
		case 50:
			if count25 == 0 {
				return "NO"
			}
			count25--
			count50++
		case 100:
			if count50 > 0 && count25 > 0 {
				count50--
				count25--
			} else if count25 >= 3 {
				count25 -= 3
			} else {
				return "NO"
			}
		default:
			return "NO"
		}
	}
	return "YES"
}

func parseCases() ([]([]int), error) {
	lines := strings.Split(strings.TrimSpace(testcasesA), "\n")
	cases := make([][]int, 0, len(lines))
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n", idx+1)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d: expected %d bills got %d", idx+1, n, len(fields)-1)
		}
		bills := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[i+1])
			if err != nil {
				return nil, fmt.Errorf("line %d: bad bill %d", idx+1, i+1)
			}
			bills[i] = val
		}
		cases = append(cases, bills)
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, bills := range cases {
		expect := expectedA(bills)
		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", len(bills))
		for i, b := range bills {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(b))
		}
		input.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		got = strings.ToUpper(got)
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
