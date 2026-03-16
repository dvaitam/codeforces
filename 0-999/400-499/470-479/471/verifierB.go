package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)
const testcasesBRaw = `
8 4 2 4 4 5 3 4 1
8 3 2 4 1 5 2 1 3
8 4 1 5 1 1 4 1 3
1 1
10 1 3 3 2 2 5 3 2 1 4
8 3 4 2 3 4 4 2 4
3 5 3 2
4 2 4 3 4
7 4 4 2 2 4 2 5
1 4
1 2
2 2 3
1 2
4 5 3 5 1
9 3 3 4 4 1 5 5 4 5
8 4 3 4 2 3 3 1 1
1 2
6 1 3 1 2 1 4
4 5 4 5 2
8 2 3 5 1 5 1 3 3
9 4 3 3 1 5 1 2 3 1
4 5 3 2 2
5 3 3 5 4 3
8 3 2 1 3 5 1 1 4
8 4 1 4 4 4 4 1 1
2 2 1
3 4 2 4
10 1 4 5 4 1 2 2 4 2 2
5 3 3 4 1 5
5 5 5 2 3 4
9 5 4 5 3 3 2 1 1 5
2 2 4
4 2 3 1 5
9 4 1 1 4 3 1 5 3 2
9 3 2 2 1 5 3 3 2 3
8 3 5 2 2 1 5 5 3
6 5 1 2 4 2 2
9 1 2 2 4 5 2 2 2 2
7 3 5 5 2 4 1 5
1 5
10 3 4 4 3 1 2 5 2 4 4
9 3 1 3 2 5 4 2 3 3
7 1 2 1 3 2 3 4
4 3 3 2 3
1 3
10 5 1 2 3 1 4 1 1 2 1
1 2
6 1 1 3 4 2 2
8 4 2 3 3 2 3 4 4
1 4
5 5 5 4 1 5
2 4 4
3 1 5 2
10 5 2 1 3 2 2 2 1 2 5
3 1 4 5
2 5 4
3 5 5 1
5 3 4 1 1 4
2 3 3
3 4 2 5
6 2 4 3 3 4 4
1 3
9 3 5 4 1 5 5 5 3 1
8 4 1 4 3 4 1 1 3
1 3
10 3 2 5 5 3 4 3 2 1 5
6 2 5 5 3 2 2
6 1 5 1 5 2 3
6 3 3 3 4 4 5
7 2 1 2 5 1 4 2
6 1 4 3 5 2 1
9 4 3 2 5 2 1 2 5 1
9 5 5 4 4 3 3 3 1 4
5 3 5 5 5 3
6 2 4 2 1 5 2
10 4 3 4 1 5 4 5 2 1 3
9 2 2 3 4 1 2 5 2 3
3 4 1 5
8 2 4 5 1 2 2 4 1
7 4 3 3 4 3 5 3
2 3 1
8 1 3 2 4 4 4 4 1
10 4 3 5 2 5 3 3 1 4 4
9 2 1 1 5 3 3 1 1 2
2 5 4
1 3
1 3
7 2 3 4 2 5 2 4
5 5 1 2 2 2
8 1 5 1 5 4 2 3 5
2 1 5
3 3 2 3
6 3 3 5 4 2 4
9 2 1 5 2 5 1 3 1 2
8 5 2 4 5 2 3 2 4
9 1 5 1 5 3 1 1 1 4
10 3 5 4 3 4 5 3 1 2 3
1 2
3 1 3 5
4 1 4 1 3
7 1 5 2 3 1 4 1
8 3 5 5 3 3 5 4 4
3 1 4 3
4 4 1 3 1
2 1 3
1 2
7 5 1 3 4 5 4 4
2 1 5
7 3 1 5 1 1 3 3
2 4 1
3 5 3 1
1 4
6 2 5 2 2 2 2
4 5 3 1 4
7 1 2 2 3 3 2 2
6 2 2 4 4 3 5
3 4 5 1
3 5 1 4
3 2 1 1
6 5 1 1 1 1 5
`


type pair struct {
	val int
	pos int
}

func expectedB(n int, arr []int) string {
	pairs := make([]pair, n)
	for i := 0; i < n; i++ {
		pairs[i] = pair{arr[i], i + 1}
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].val < pairs[j].val })
	l1, l2 := -1, -1
	for i := 0; i < n-1; i++ {
		if pairs[i].val == pairs[i+1].val {
			if l1 == -1 {
				l1 = i
			} else if l2 == -1 {
				l2 = i
			}
		}
	}
	if l2 == -1 {
		return "NO"
	}
	base := make([]int, n)
	for i := 0; i < n; i++ {
		base[i] = pairs[i].pos
	}
	join := func(seq []int) string {
		sb := make([]string, len(seq))
		for i, v := range seq {
			sb[i] = strconv.Itoa(v)
		}
		return strings.Join(sb, " ")
	}
	seq1 := join(base)
	seq2 := append([]int(nil), base...)
	seq2[l1], seq2[l1+1] = seq2[l1+1], seq2[l1]
	seq2s := join(seq2)
	seq3 := append([]int(nil), base...)
	seq3[l2], seq3[l2+1] = seq3[l2+1], seq3[l2]
	seq3s := join(seq3)
	return "YES\n" + seq1 + "\n" + seq2s + "\n" + seq3s
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesBRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 1 {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		if len(parts) != n+1 {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(parts[1+i])
			arr[i] = v
		}
		expect := expectedB(n, arr)
		var sb strings.Builder
		sb.WriteString(line)
		sb.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(sb.String())
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
		if got != expect {
			fmt.Printf("test %d failed. Expected:\n%s\nGot:\n%s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
