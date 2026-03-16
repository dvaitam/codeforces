package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)
const testcasesARaw = `
865 25 14 11
266 33 16 104
941 20 16 92
598 14 17 36
289 9 4 159
819 17 18 181
830 39 5 80
102 47 3 176
339 31 18 26
363 28 11 157
656 14 18 123
454 34 9 16
825 36 1 24
737 26 1 157
506 22 8 187
334 46 3 49
940 37 8 62
823 10 18 115
94 6 11 131
956 32 4 78
565 19 4 141
341 35 7 155
561 38 10 114
94 39 13 82
590 16 10 48
194 12 2 157
673 17 16 18
92 44 5 39
946 3 3 180
946 35 13 181
538 18 17 61
870 14 19 108
594 18 15 127
677 42 12 22
333 40 4 125
602 41 11 49
249 2 9 30
723 15 12 44
341 28 2 26
802 10 8 12
837 37 18 155
697 5 1 32
651 13 20 148
123 26 3 95
854 8 2 156
23 13 6 184
127 31 7 187
820 4 1 140
436 40 4 67
72 15 3 166
309 23 14 47
63 33 15 11
611 7 13 52
267 23 16 146
174 45 7 197
60 44 6 42
351 34 9 31
612 29 6 4
483 44 14 146
896 33 10 167
366 25 9 40
575 45 1 118
760 6 11 190
47 35 9 35
246 49 16 91
625 19 12 152
970 41 20 34
733 20 13 192
425 42 3 1
609 13 11 41
246 15 15 97
728 44 19 107
33 26 19 108
791 43 2 43
457 5 9 180
162 29 17 125
930 36 20 194
1 3 16 84
320 30 2 107
193 36 3 186
134 1 13 174
428 21 1 55
15 46 1 173
542 40 4 49
122 39 7 78
287 45 6 26
488 26 3 6
282 29 4 66
137 42 17 167
661 23 4 40
286 2 2 11
211 44 9 143
323 24 19 11
867 48 20 168
507 46 15 164
446 24 18 46
213 25 19 75
10 9 5 70
342 22 12 184
96 22 20 10
`


func expected(n, m, a, b int) int {
	cost1 := n * a
	cost2 := (n/m)*b + (n%m)*a
	cost3 := ((n + m - 1) / m) * b
	if cost2 < cost1 {
		cost1 = cost2
	}
	if cost3 < cost1 {
		cost1 = cost3
	}
	return cost1
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesARaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		var n, m, a, b int
		fmt.Sscan(line, &n, &m, &a, &b)
		exp := expected(n, m, a, b)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d %d %d\n", n, m, a, b)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != fmt.Sprint(exp) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
