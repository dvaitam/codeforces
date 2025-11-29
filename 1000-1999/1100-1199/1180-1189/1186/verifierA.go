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

const testcasesRaw = `50 98 54
6 34 66
63 52 39
62 46 75
28 65 18
37 18 97
13 80 33
69 91 78
19 40 13
94 10 88
43 61 72
13 46 56
41 79 82
27 71 62
57 67 34
8 71 2
12 93 52
91 86 81
1 79 64
43 32 94
42 91 9
25 73 29
31 19 70
58 12 11
41 66 63
14 39 71
38 91 16
71 43 70
27 78 71
76 37 57
12 77 50
41 74 31
38 24 25
24 5 79
85 34 61
9 12 87
97 17 20
5 11 90
70 88 51
91 68 36
67 31 28
87 76 54
75 36 58
64 85 83
90 46 11
42 79 15
63 76 81
43 25 32
3 94 35
15 91 29
48 22 43
55 8 13
19 90 29
6 74 82
69 78 88
10 4 16
82 25 78
74 16 51
12 48 15
5 78 3
25 24 92
16 62 27
94 8 87
3 70 55
80 13 34
9 29 10
83 39 45
56 24 8
65 60 6
77 13 90
51 26 34
46 94 61
73 22 90
87 27 99
8 87 21
21 44 68
33 16 77
57 86 23
2 61 88
53 73 66
40 84 46
50 85 33
20 72 89
2 59 95
11 43 95
6 70 36
18 31 98
62 46 79
37 87 46
76 82 80
17 92 40
50 96 54
84 11 1
77 25 90
43 21 31
29 82 58
49 91 87
73 54 5
52 90 73
54 99 85`

// referenceSolve mirrors 1186A.go.
func referenceSolve(n, m, k int) string {
	if m >= n && k >= n {
		return "Yes"
	}
	return "No"
}

func runCase(bin string, n, m, k int) error {
	expect := referenceSolve(n, m, k)
	input := fmt.Sprintf("%d %d %d\n", n, m, k)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scanner.Split(bufio.ScanWords)
	idx := 0
	for {
		if !scanner.Scan() {
			break
		}
		n, _ := strconv.Atoi(scanner.Text())
		if !scanner.Scan() {
			fmt.Printf("case %d: missing m\n", idx+1)
			os.Exit(1)
		}
		m, _ := strconv.Atoi(scanner.Text())
		if !scanner.Scan() {
			fmt.Printf("case %d: missing k\n", idx+1)
			os.Exit(1)
		}
		k, _ := strconv.Atoi(scanner.Text())
		idx++
		if err := runCase(bin, n, m, k); err != nil {
			fmt.Printf("case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
