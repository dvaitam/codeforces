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

const testcasesRaw = `2 2 1 11 11 16 12 34 18
3 1 2 0 0 1 18 20 22 29 4 3 14 18
3 1 2 1 1 0 18 16 32 32 9 2 26 2
1 3 2 011 8 11 10 17 19 8 26 12
3 2 1 01 10 11 4 18 14 35
1 3 3 110 20 13 30 31 8 10 13 16 6 2 25 10
2 1 1 0 0 2 3 19 15
3 3 2 001 111 101 20 4 35 22 11 7 18 7
3 2 1 01 01 10 4 5 11 6
3 3 3 000 001 010 2 20 2 26 6 4 21 10 2 1 19 14
3 1 2 0 0 0 10 12 23 17 2 17 16 18
3 1 3 1 0 1 12 16 30 21 7 2 12 7 11 17 19 20
3 2 3 00 11 11 13 9 17 26 1 15 3 25 2 18 10 22
1 2 2 11 19 20 23 29 13 14 33 16
1 3 1 100 8 15 20 33
2 1 2 1 0 6 15 8 23 6 15 22 30
3 3 1 011 110 100 5 1 17 14
2 1 1 0 0 17 20 20 26
1 3 3 011 6 4 21 16 3 1 11 15 4 9 8 29
3 3 3 100 100 001 18 11 29 29 2 20 22 35 15 14 26 31
1 1 2 1 1 5 5 13 11 11 22 13
2 3 1 010 011 13 18 17 27
1 2 3 00 10 6 26 8 10 13 20 22 14 4 17 21
2 2 2 10 10 16 14 17 23 11 5 16 25
3 2 3 00 00 00 13 1 16 13 18 17 27 31 16 19 22 32
1 2 1 10 14 7 25 10
1 3 1 100 16 13 24 19
3 1 1 0 0 0 15 13 26 30
1 1 3 1 5 19 17 39 14 17 29 27 16 16 36 22
3 3 1 011 100 101 19 10 34 12
1 3 1 000 2 10 2 24
2 1 1 1 1 17 13 33 29
1 3 1 010 10 18 29 31
2 2 3 00 00 10 17 28 25 11 3 26 11 10 14 22 26
1 1 3 0 8 10 18 11 2 16 15 20 16 20 18 24
2 2 1 11 10 4 16 8 16
1 3 3 010 18 12 24 24 16 4 17 23 15 20 35 30
3 1 3 1 0 1 10 4 26 10 2 13 16 24 7 15 18 35
1 1 1 1 9 1 25 19
3 1 1 0 1 1 4 5 17 23
2 1 1 1 0 4 14 8 14
2 2 3 10 11 9 3 20 5 4 12 4 23 12 6 12 13
2 1 3 0 0 1 7 4 7 10 12 10 31 8 5 13 19
1 2 2 10 1 7 12 17 16 10 25 27
3 2 1 00 10 10 5 8 15 24
1 1 1 1 12 14 13 18
3 1 2 0 0 0 14 10 31 23 5 19 18 28
3 2 1 01 10 11 1 14 11 28
1 2 2 10 6 4 14 7 18 20 22 34
2 1 2 1 0 8 15 18 31 5 12 19 32
3 1 2 0 1 0 15 20 29 20 7 10 10 30
2 3 3 011 010 18 13 26 33 1 4 9 5 1 9 13 25
3 3 2 101 110 000 9 10 26 20 4 17 11 22
1 2 2 10 19 17 39 23 18 4 31 24
3 2 3 11 11 00 4 4 16 16 19 15 23 32 10 12 30 27
3 2 1 11 11 01 10 5 25 6
3 1 1 1 1 1 1 17 3 19
3 3 3 101 000 110 5 19 14 25 4 14 18 24 13 6 23 19
3 3 2 010 100 011 13 14 28 26 8 7 22 13
3 3 1 100 001 000 20 10 39 12
3 3 2 111 011 111 7 11 15 12 2 2 7 13
1 2 3 00 3 14 10 33 13 18 20 32 7 11 26 14
3 1 2 1 1 1 9 1 25 2 7 12 9 18
3 2 1 01 11 11 16 12 23 13
2 3 1 011 101 16 15 30 18
1 1 1 0 5 14 11 28
3 1 2 1 0 0 8 16 15 20 9 12 19 25
1 3 2 011 17 20 31 37 9 9 16 9
1 3 3 001 8 7 17 7 18 17 31 18 4 13 24 21
1 3 3 101 8 8 10 24 10 11 17 22 16 10 34 15
1 1 3 1 12 19 32 19 5 13 9 18 17 3 21 9
2 3 3 000 011 20 19 24 39 16 4 35 4 17 20 28 35
2 2 1 00 11 18 11 20 19
1 3 2 011 13 2 19 3 11 8 21 22
3 3 3 011 010 100 12 1 27 21 2 1 9 2 1 8 21 18
1 1 2 1 5 7 19 20 5 12 14 17
3 2 3 11 01 11 2 19 5 32 13 6 13 22 5 20 21 24
1 2 1 00 1 6 18 11
3 1 2 0 1 0 20 20 21 28 11 13 11 33
1 2 1 11 5 15 12 31
2 1 3 1 1 9 16 21 16 10 17 19 34 16 2 33 20
3 2 3 01 10 11 16 2 16 10 2 9 20 18 7 17 23 27
2 2 1 01 01 6 5 16 5
3 1 3 0 1 1 10 10 20 25 13 20 26 25 1 5 19 6
2 1 2 0 1 9 20 15 22 18 14 26 19
3 1 1 0 0 1 14 9 23 18
1 2 3 11 18 17 35 27 11 7 24 11 1 17 5 35
2 2 2 01 00 12 17 17 23 12 16 12 23
3 1 2 0 1 0 19 15 26 29 17 4 23 9
2 1 2 1 1 9 14 20 33 11 3 20 3
2 1 2 0 1 13 14 33 34 13 2 31 16
2 3 1 110 110 2 3 20 14
2 1 1 0 0 18 16 19 26
1 2 2 01 14 5 33 9 13 10 29 11
1 1 1 1 2 17 3 34
3 3 3 101 000 101 11 9 19 25 15 5 29 22 5 2 25 20
1 3 3 010 7 15 26 22 15 17 20 27 5 16 22 17
3 1 3 1 0 0 4 14 23 25 19 15 29 27 17 12 37 15
1 2 1 00 1 11 20 17
1 2 3 01 13 2 32 7 12 3 25 4 15 12 34 31`

// Embedded reference logic from 1202E.go.
type node struct {
	next [26]int
	fail int
	out  int
}

func newNode() node {
	n := node{fail: 0, out: 0}
	for i := 0; i < 26; i++ {
		n.next[i] = -1
	}
	return n
}

func buildAC(patterns []string) []node {
	nodes := make([]node, 1)
	nodes[0] = newNode()
	for _, s := range patterns {
		v := 0
		for i := 0; i < len(s); i++ {
			c := int(s[i] - 'a')
			if nodes[v].next[c] == -1 {
				nodes = append(nodes, newNode())
				nodes[v].next[c] = len(nodes) - 1
			}
			v = nodes[v].next[c]
		}
		nodes[v].out++
	}
	q := make([]int, 0)
	for c := 0; c < 26; c++ {
		v := nodes[0].next[c]
		if v != -1 {
			nodes[v].fail = 0
			q = append(q, v)
		} else {
			nodes[0].next[c] = 0
		}
	}
	for idx := 0; idx < len(q); idx++ {
		v := q[idx]
		f := nodes[v].fail
		nodes[v].out += nodes[f].out
		for c := 0; c < 26; c++ {
			u := nodes[v].next[c]
			if u != -1 {
				nodes[u].fail = nodes[f].next[c]
				q = append(q, u)
			} else {
				nodes[v].next[c] = nodes[f].next[c]
			}
		}
	}
	return nodes
}

func countEnds(nodes []node, text string) []int {
	n := len(text)
	res := make([]int, n)
	state := 0
	for i := 0; i < n; i++ {
		c := int(text[i] - 'a')
		state = nodes[state].next[c]
		res[i] = nodes[state].out
	}
	return res
}

func countStarts(nodes []node, text string) []int {
	n := len(text)
	res := make([]int, n)
	state := 0
	for i := n - 1; i >= 0; i-- {
		c := int(text[i] - 'a')
		state = nodes[state].next[c]
		res[i] = nodes[state].out
	}
	return res
}

func referenceSolve(text string, patterns []string) int64 {
	acForward := buildAC(patterns)
	rev := make([]string, len(patterns))
	for i := 0; i < len(patterns); i++ {
		b := []byte(patterns[i])
		for l, r := 0, len(b)-1; l < r; l, r = l+1, r-1 {
			b[l], b[r] = b[r], b[l]
		}
		rev[i] = string(b)
	}
	acBackward := buildAC(rev)

	pref := countEnds(acForward, text)
	suff := countStarts(acBackward, text)

	var ans int64
	for i := 0; i < len(text)-1; i++ {
		ans += int64(pref[i]) * int64(suff[i+1])
	}
	return ans
}

// encodeToken turns a numeric token into a lowercase string so that the solver receives valid input.
func encodeToken(tok string) string {
	val, err := strconv.ParseInt(tok, 10, 64)
	if err != nil {
		// fallback: map characters directly to letters
		var b strings.Builder
		for i := 0; i < len(tok); i++ {
			ch := tok[i]
			b.WriteByte('a' + ch%26)
		}
		return b.String()
	}
	if val == 0 {
		return "a"
	}
	var rev []byte
	for val > 0 {
		rev = append(rev, byte('a'+(val%26)))
		val /= 26
	}
	for i, j := 0, len(rev)-1; i < j; i, j = i+1, j-1 {
		rev[i], rev[j] = rev[j], rev[i]
	}
	return string(rev)
}

func runCase(bin, line string) error {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return fmt.Errorf("invalid test line")
	}
	text := encodeToken(fields[0])
	patterns := make([]string, len(fields)-1)
	for i := 1; i < len(fields); i++ {
		patterns[i-1] = encodeToken(fields[i])
	}
	var input strings.Builder
	fmt.Fprintf(&input, "%s\n%d\n", text, len(patterns))
	for _, p := range patterns {
		input.WriteString(p)
		input.WriteByte('\n')
	}

	expect := referenceSolve(text, patterns)

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(gotStr, 10, 64)
	if err != nil {
		return fmt.Errorf("cannot parse output %q", gotStr)
	}
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		if err := runCase(bin, line); err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
