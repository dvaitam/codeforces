package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded testcases (one number per line).
const embeddedTestcases = `58716
11955
10544
41950
66576
64131
14294
39511
72255
38153
92610
16359
71754
43614
70816
26634
79060
71726
77020
37703
58325
12010
78156
50449
41555
75451
31733
38054
24100
24823
24475
4321
80317
86069
34086
62459
9055
11773
88961
99300
17068
19601
5064
10518
91661
70857
89587
51287
92442
68756
36127
68392
30867
28206
89060
77306
54974
75981
36072
59056
64573
86539
84042
91779
46840
10796
42509
80318
15119
63759
76949
82594
43944
24953
31855
2124
95877
35525
15353
92449
28896
48766
22345
43586
55853
8151
13186
19183
91445
28675
5928
75217
83126
70018
78927
89206
9698
3499
16311
83230`

// Embedded reference solution (from 971B.go).
func solveString(s string) string {
	sum := 0
	for _, ch := range s {
		sum += int(ch - '0')
	}
	return fmt.Sprint(sum)
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(embeddedTestcases))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		s := line
		input := s + "\n"
		expected := solveString(s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
