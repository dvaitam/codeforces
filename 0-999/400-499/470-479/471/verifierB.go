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
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
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
