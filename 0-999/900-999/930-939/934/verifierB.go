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

func solveCase(k int) string {
	if k > 36 {
		return "-1"
	}
	half := k / 2
	if k%2 == 0 {
		return strings.Repeat("8", half)
	}
	return strings.Repeat("8", half) + "4"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	const testcasesRaw = `100
9
37
55
52
49
5
17
8
32
49
29
31
42
25
51
14
7
32
2
58
54
25
28
39
49
50
1
45
29
18
47
52
15
38
7
58
21
2
2
2
42
35
1
57
25
44
14
28
47
2
34
15
49
29
32
36
15
23
15
44
15
49
30
19
60
2
27
54
59
36
60
42
7
12
41
47
56
19
8
48
22
58
47
46
33
60
28
33
54
59
43
13
20
19
38
57
32
55
33
26`

	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		k, _ := strconv.Atoi(scan.Text())
		expected := solveCase(k)
		input := fmt.Sprintf("%d\n", k)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n%s", caseNum, err, out.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", caseNum, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
