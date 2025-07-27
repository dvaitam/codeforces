package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func solveCase(n int, lists [][]int) (bool, int, int) {
	used := make([]bool, n+1)
	unmatched := 0
	for i := 1; i <= n; i++ {
		matched := false
		for _, p := range lists[i-1] {
			if !used[p] {
				used[p] = true
				matched = true
				break
			}
		}
		if !matched {
			unmatched = i
		}
	}
	if unmatched == 0 {
		return false, 0, 0
	}
	prince := 0
	for j := 1; j <= n; j++ {
		if !used[j] {
			prince = j
			break
		}
	}
	return true, unmatched, prince
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(2)
	t := 100
	cases := make([][][]int, t)
	ns := make([]int, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(20) + 1
		ns[i] = n
		lists := make([][]int, n)
		for d := 0; d < n; d++ {
			k := rand.Intn(n + 1)
			perm := rand.Perm(n)
			lst := perm[:k]
			for idx := range lst {
				lst[idx]++
			}
			sort.Ints(lst)
			lists[d] = append([]int(nil), lst...)
		}
		cases[i] = lists
	}

	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", t))
	for idx := 0; idx < t; idx++ {
		n := ns[idx]
		input.WriteString(fmt.Sprintf("%d\n", n))
		lists := cases[idx]
		for _, lst := range lists {
			input.WriteString(fmt.Sprintf("%d", len(lst)))
			for _, v := range lst {
				input.WriteString(fmt.Sprintf(" %d", v))
			}
			input.WriteByte('\n')
		}
	}
	in := input.String()

	var expectedOut strings.Builder
	for idx := 0; idx < t; idx++ {
		improve, girl, prince := solveCase(ns[idx], cases[idx])
		if improve {
			expectedOut.WriteString("IMPROVE\n")
			expectedOut.WriteString(fmt.Sprintf("%d %d\n", girl, prince))
		} else {
			expectedOut.WriteString("OPTIMAL\n")
		}
	}
	want := strings.TrimSpace(expectedOut.String())

	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Printf("Runtime error: %v\n%s", err, out.String())
		os.Exit(1)
	}
	gotLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	wantLines := strings.Split(want, "\n")
	if len(gotLines) != len(wantLines) {
		fmt.Println("Wrong answer: line count mismatch")
		os.Exit(1)
	}
	for i := range wantLines {
		if strings.TrimSpace(gotLines[i]) != strings.TrimSpace(wantLines[i]) {
			fmt.Printf("Wrong answer on line %d: expected %s got %s\n", i+1, wantLines[i], gotLines[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
