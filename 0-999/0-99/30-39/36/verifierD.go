package main

import (
	"bytes"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesD.txt.
const embeddedTestcasesD = `100 3
8 19
18 5
12 20
16 19
3 20
1 16
9 18
8 7
16 18
18 16
13 5
8 5
17 13
1 3
6 19
2 10
1 9
16 20
13 14
13 19
15 5
12 4
2 5
16 7
9 14
10 14
17 13
19 12
18 19
14 19
8 11
1 9
20 6
11 18
19 19
4 7
19 9
10 4
3 16
16 3
12 3
14 5
1 10
14 14
4 2
20 20
2 13
19 11
18 9
17 8
2 10
1 3
4 20
18 2
7 14
10 20
9 5
2 11
11 12
5 13
13 15
17 13
20 18
4 20
17 9
14 8
10 14
9 17
10 18
11 1
14 19
11 1
13 20
19 5
2 11
15 12
12 20
9 16
1 19
2 1
12 9
15 10
19 20
11 6
12 6
11 12
20 9
10 13
4 1
19 5
10 17
8 9
8 11
6 14
4 4
20 11
11 8
15 6
3 11
7 19`

func solve36DCase(n, m, k int64) string {
	mod := k + 1
	an, rn := n/mod, n%mod
	bm, rm := m/mod, m%mod
	var win bool
	if rn == rm {
		if rn == k {
			win = true
		} else if rn%2 == 1 {
			win = false
		} else {
			win = (an%2 == 1)
		}
	} else {
		d := rn - rm
		if d < 0 {
			d = -d
		}
		win = (d%2 == 1)
	}
	if win {
		return "+"
	}
	return "-"
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data := strings.TrimSpace(embeddedTestcasesD)
	if data != "" {
		data += "\n"
	}

	scan := bufio.NewScanner(strings.NewReader(data))
	scan.Split(bufio.ScanWords)
	next := func() string {
		if !scan.Scan() {
			fmt.Fprintln(os.Stderr, "bad test data")
			os.Exit(1)
		}
		return scan.Text()
	}
	t, err := strconv.Atoi(next())
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid t value")
		os.Exit(1)
	}
	kVal, err := strconv.ParseInt(next(), 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid k value")
		os.Exit(1)
	}

	expected := make([]string, t)
	for i := 0; i < t; i++ {
		n, err1 := strconv.ParseInt(next(), 10, 64)
		m, err2 := strconv.ParseInt(next(), 10, 64)
		if err1 != nil || err2 != nil {
			fmt.Fprintf(os.Stderr, "invalid case %d\n", i+1)
			os.Exit(1)
		}
		expected[i] = solve36DCase(n, m, kVal)
	}

	gotRaw, err := runCandidate(bin, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate execution failed: %v\n", err)
		os.Exit(1)
	}
	got := strings.Fields(gotRaw)
	if len(got) != len(expected) {
		fmt.Fprintf(os.Stderr, "output length mismatch: expected %d tokens, got %d\n", len(expected), len(got))
		os.Exit(1)
	}
	for i := 0; i < t; i++ {
		if got[i] != expected[i] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expected[i], got[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
