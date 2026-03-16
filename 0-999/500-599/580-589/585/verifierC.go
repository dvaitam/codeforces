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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func expectedC(x, y int64) string {
	if gcd(x, y) != 1 {
		return "Impossible\n"
	}
	type step struct {
		k int64
		c byte
	}
	var res []step
	for x > 1 || y > 1 {
		if x > y {
			k := (x - 1) / y
			res = append(res, step{k, 'A'})
			x -= k * y
		} else {
			k := (y - 1) / x
			res = append(res, step{k, 'B'})
			y -= k * x
		}
	}
	var sb strings.Builder
	for _, st := range res {
		sb.WriteString(fmt.Sprintf("%d%c", st.k, st.c))
	}
	sb.WriteByte('\n')
	return sb.String()
}

const testcasesCRaw = `19 44
38 31
40 39
43 12
40 29
21 38
37 50
41 10
33 28
24 25
19 16
14 17
16 35
39 22
40 41
10 21
49 17
35 46
46 39
8 29
31 14
39 23
39 16
35 32
14 19
25 18
33 13
23 27
21 47
2 3
39 11
6 17
33 19
45 32
25 4
49 18
40 47
36 29
27 31
47 6
41 6
33 2
43 38
47 39
7 37
3 7
47 20
9 50
37 40
34 13
18 19
11 12
38 15
15 17
9 43
13 38
2 23
35 43
31 44
9 49
6 47
15 1
48 37
44 1
35 3
15 43
21 43
29 15
42 41
47 18
13 7
23 32
18 23
35 38
19 14
19 36
26 5
12 25
34 9
17 47
36 72
50 100
49 98
18 54
22 66
34 68
24 72
49 196
26 104
10 40
47 188
38 76
43 86
21 84
30 120
50 200
28 56
31 62
27 54
8 32
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesCRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		var xStr, yStr string
		fmt.Sscan(line, &xStr, &yStr)
		x, _ := strconv.ParseInt(xStr, 10, 64)
		y, _ := strconv.ParseInt(yStr, 10, 64)
		expect := expectedC(x, y)
		input := fmt.Sprintf("%s %s\n", xStr, yStr)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		expectTrim := strings.TrimSpace(expect)
		if got != expectTrim {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, expectTrim, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
