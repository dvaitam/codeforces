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

const testcasesBRaw = `100
25 rbrrrrrrbbrbbrbrbbrrbrrrb
36 rrrbbbrbrrbrbbbbbrbbrrbrrbrbrrbrrbrb
13 bbrrbbbbbbrrb
14 rrrrrbrbrrbbbr
8 brbrrbbb
45 bbbbbbbrbrbbbbbbrbbbrbrbbbrrrbbrbbrbrrrbbbbbr
34 rbrbbrrrrrrbbrbrbrbbrrrrbrrrbbbrbb
15 rrrbrrbbrbrbrrb
3 rrr
30 brbbbbrrrbbbbbbbbrrbbrrbbrrbrb
42 rbbrbbbbrrrbrrrrbbrrbbbrrrrbrbbrbbrrrbrbrb
4 rbbr
26 rrrbbrrrrbrbrrbrrbbrbbbbbb
25 bbrrrrbrbrbrbrbrbbbrbbrrr
14 bbbbbrrrbbrbrr
34 rrrrbbbrrrbbrbrrrbbbbbbbrbrrbbrrrb
37 bbrbrrrrbbbbrrrbrrrrrrbbbbrrbbrrbrbbr
27 brrrbbrbbbbrbrbrrbbrrbrbrrb
34 bbrrrbrrbbbrrbbbbrrbbrrbbrbbrbbrrr
2 rr
47 rbrbbrbrrbbbrbbbbbbbrrbbbbrbrrrbbbrrrrrrbbbrbrb
9 brbbbrrrb
9 brbbbrrrb
38 rrrbbrrbrrbrrrbrrrbbbrbbrrbrrbbrrbrrbr
49 brbrrbbrbrbrrbrbrrbbrbbrrrrbrrrbbbbrrrbbbbrrrbbbr
41 rrrrrbbbbrrrbrrrrbrrrrbrrbrbbbrrrbbbrbrbb
40 brrbbrbbrrrrbrrrbrbrbbbrrrrrrbbrbbrbbbbr
4 bbrb
46 rrrrbrrrrrbrrbbbbrbrbbbrbrbrbrbbrrrrrbbbrbbrbb
17 rrrrrrbbrbbrrrbrr
30 rbbbbbbrbrbrrbbbrbbrrrrbrbrrrr
18 bbbbbrbbrbrbbrrbrb
44 rbbbrrbrrrbbbrrbbrbbbbbrbbbbrrbrbbrrrrbbbrrr
46 brbrbrrrbbbrbrrbrrbrbrbbrbrbbbbbbrbbrrbbrrbrrb
42 rrrbrrrbbrrbbbbbrbbbbbbbrbrbbrrrbbrbrrbrbr
11 rrrrrbrrrbr
44 brrbrrrbbrbrrbrrrbbrbrbbrbbbrrrrrrrrbbbbrbrb
47 rrbbrrbbbbbrrrrrbrrrrrbrbbbrrrbrbbrbbrrbbbrbrbr
16 rbbbrbrrrrrrrbrb
32 brbrrrrbrrbrrbrrbbbrrbbbbrbrbrrb
49 rrbbrrbbbbrbbrbrbbbrbrrrrrbrbbbbrbbrbrbrbrbbrbbbr
39 rrrrrbbrbbbbrbbrbrrbbrbrbrrrrrrbbrrbrbr
7 bbrbbrb
21 rrrbbrrbbbrrbrbbrbbrr
11 bbbbbrbrrbb
16 rrbbrbbrrrbrbbbr
46 bbbbrbbbbbbrbrrrbbrbrbbrrbrrbrrbrbrrbrrrbrbrbr
26 bbbrbrrrbrrbbrrrbrbrrbrrbb
17 brrrbrrbrrrbrbbbr
11 brrrbbrbbbb
11 bbbrrrrbrrb
22 rbrrbrrrrrrrbbbbbbbrbr
6 bbbrrr
7 rrrrbbb
50 bbbrbrrrbbrbrbrbrrbrbrbrrrrrbbrbrbrbbrrrbrrbbrbrbb
43 bbbbbrrbrrrbbbbrrrbbrrrbbrrbrrbbrbrrrrbbrrr
15 rrrrrrrbrrbrbrb
48 brrbbrrrrbrrrbrbbrrrbrbbrbrbrbbbrbbbbrbbbbbrbbrb
35 bbbrrrbrbrbrbrrrrbbrbbrbbbrbrbrbbrb
33 rrbrrbrrrbrbrbrrrbrbrrbrbrrrbbrrr
42 rbbbrbbbrbrbrrbrrbbrbrbrbrbrbrbbrrbbbbbbrb
43 brrrrbbrrrrrrrbbbbbrrbrbbrbbrrbbbrrrbbrbbrr
20 bbbbrbrrbrrbrrrbrbbb
21 rbrrrrbbrbbbbrbbrrrbr
19 rbbrbbrbrrbrbrbrbbb
32 bbrbbbbrbrrbrrbbrrbrbbrrrbbbrbbr
13 rrbrrbbrrbrbr
44 rbbbrbrrrbbbrbrrrbrrrrbrrbbbrrbbbrbbbrbrrbbr
19 rrrrbbrbbbbrbbrbbbr
46 bbbbbrbrbbrbrrbrrbbbbbbbrrrrrrrbbbbbbbrbrbrbbr
43 rrrrbbrbrrrrrbrrrrrrrrbbrbrbbrrrrrrbbrrbbbb
17 bbbbrrbbbbrbrrrrb
38 rbbbbrrbrbrbbrbrbbbrbbbrrrbbrbrrrrrrbr
49 bbrrrrbrbbbbrrbrrrbrbbbbbbrrrbrbbrbbbrbbrbrrrbrbb
42 bbbrrrrbrbbbbrbbbbrbrbrbbrbrrrrrbbrrrrrrbb
10 rrbbbbbrbb
34 rrrrrrrrbbrrrbbbrbrrbbbbrbrbrrrrrb
42 brrbbrbrrrbbrbrrrrbrrbrrrbbrbrrbrbrbbrbbrr
14 brbbrbrrrbrrrb
3 bbb
36 rbrrrrrrbbbrbrrrrrrbbbbbrbbbbbrbrbbr
19 bbbbbbrrrbrbbbrbbbr
35 rbrbrbbbbbrbbrbrbbrrrbrbbrrrrrbrrbb
25 bbbrbrrbbbrrrrbrrrbbbrrrb
30 brbbrbrbbbbbbrrbbbbrbbbbbrrrrb
23 rbbbrbbbbbrrbrrrrrbbrrb
15 brbbrrbrbbrrbrb
13 rrrrbbrbrbbrr
10 brrbbbbrrb
3 rrb
44 rbbrbbbbrrrrbbbbrbbbbbbbbbrbrrrrbbbrbrrrrrbr
42 bbbbbbrrbrrbrrbrbrrbrrbrbrbrrbrrrbbbrbrrbr
31 rrbrrbbbrrbbrrrbrrbbbrrrbrrrrrb
6 rbbbrr
16 rbrbrbrrbrbbbbrr
29 bbrrrbrrbbbbrbrbrbbrbbbbrbbrr
45 brbrbbbrbrbrrbrrbrbrbrrrrrrrrbbbrrbrrrbbbrbbr
45 rbbrrrrrbbbbrbbrrbbrbrrbbrrrrbbbrrrrbrrrbrrrb
26 rrrrrrrbrbrrrrbbrbrbrbbrbb
12 bbbbrrbrrrrb
`

func expected(s string) int {
	n := len(s)
	var cost1_b, cost1_r int // pattern starting with 'r'
	var cost2_b, cost2_r int // pattern starting with 'b'
	for i := 0; i < n; i++ {
		c := s[i]
		if i%2 == 0 {
			if c == 'b' {
				cost1_b++
			} else {
				cost2_r++
			}
		} else {
			if c == 'r' {
				cost1_r++
			} else {
				cost2_b++
			}
		}
	}
	cost1 := cost1_b
	if cost1_r > cost1 {
		cost1 = cost1_r
	}
	cost2 := cost2_b
	if cost2_r > cost2 {
		cost2 = cost2_r
	}
	if cost1 < cost2 {
		return cost1
	}
	return cost2
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scan := bufio.NewScanner(strings.NewReader(testcasesBRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "empty test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for idx := 1; idx <= t; idx++ {
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "bad test format at case %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "missing string at case %d\n", idx)
			os.Exit(1)
		}
		s := scan.Text()
		if len(s) != n {
			fmt.Fprintf(os.Stderr, "case %d: length mismatch\n", idx)
			os.Exit(1)
		}
		input := fmt.Sprintf("%d\n%s\n", n, s)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid output\n", idx)
			os.Exit(1)
		}
		exp := expected(s)
		if got != exp {
			fmt.Printf("case %d failed\nexpected: %d\ngot: %d\n", idx, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
