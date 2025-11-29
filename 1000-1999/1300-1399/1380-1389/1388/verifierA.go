package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesA.txt.
const testcasesRaw = `100
50
98
54
6
34
66
63
52
39
62
46
75
28
65
18
37
18
97
13
80
33
69
91
78
19
40
13
94
10
88
43
61
72
13
46
56
41
79
82
27
71
62
57
67
34
8
71
2
12
93
52
91
86
81
1
79
64
43
32
94
42
91
9
25
73
29
31
19
70
58
12
11
41
66
63
14
39
71
38
91
16
71
43
70
27
78
71
76
37
57
12
77
50
41
74
31
38
24
25
24`

func expected(n int) (string, string) {
	if n <= 30 {
		return "NO", ""
	}
	switch n {
	case 36:
		return "YES", "5 6 10 15"
	case 40:
		return "YES", "5 6 14 15"
	case 44:
		return "YES", "6 7 10 21"
	default:
		return "YES", fmt.Sprintf("6 10 14 %d", n-30)
	}
}

func runCase(bin string, n int) error {
	input := fmt.Sprintf("1\n%d\n", n)
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	outLines := strings.Fields(strings.TrimSpace(out.String()))
	exp1, exp2 := expected(n)
	if len(outLines) == 0 {
		return fmt.Errorf("no output")
	}
	if outLines[0] != exp1 {
		return fmt.Errorf("expected %s got %s", exp1, outLines[0])
	}
	if exp1 == "YES" {
		got := strings.Join(outLines[1:], " ")
		if got != exp2 {
			return fmt.Errorf("expected numbers %q got %q", exp2, got)
		}
	} else {
		if len(outLines) > 1 {
			return fmt.Errorf("extra output after NO")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		fmt.Println("no embedded testcases")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(fields[0])
	if len(fields) != t+1 {
		fmt.Println("embedded testcases count mismatch")
		os.Exit(1)
	}
	for i := 0; i < t; i++ {
		n, _ := strconv.Atoi(fields[i+1])
		if err := runCase(bin, n); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
