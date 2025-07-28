package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func runCmd(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = io.Discard
	err := cmd.Run()
	return out.String(), err
}

func solve(a []int, s string) string {
	blues := make([]int, 0, len(a))
	reds := make([]int, 0, len(a))
	for i, v := range a {
		if s[i] == 'B' {
			blues = append(blues, v)
		} else {
			reds = append(reds, v)
		}
	}
	sort.Ints(blues)
	sort.Sort(sort.Reverse(sort.IntSlice(reds)))
	pos := 1
	for _, v := range blues {
		if v < pos {
			return "NO"
		}
		pos++
	}
	pos = len(a)
	for _, v := range reds {
		if v > pos {
			return "NO"
		}
		pos--
	}
	return "YES"
}

func randCase() ([]int, string) {
	n := rand.Intn(20) + 1
	a := make([]int, n)
	var sb strings.Builder
	for i := 0; i < n; i++ {
		a[i] = rand.Intn(10) + 1
		if rand.Intn(2) == 0 {
			sb.WriteByte('B')
		} else {
			sb.WriteByte('R')
		}
	}
	return a, sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	rand.Seed(4)
	casesA := [][]int{{1}, {2}, {3, 1, 2}}
	casesS := []string{"B", "R", "BRB"}
	for len(casesA) < 100 {
		a, s := randCase()
		casesA = append(casesA, a)
		casesS = append(casesS, s)
	}

	for i := range casesA {
		a := casesA[i]
		s := casesS[i]
		input := fmt.Sprintf("1\n%d\n", len(a))
		for j, v := range a {
			if j+1 == len(a) {
				input += fmt.Sprintf("%d\n", v)
			} else {
				input += fmt.Sprintf("%d ", v)
			}
		}
		input += fmt.Sprintf("%s\n", s)
		expected := solve(append([]int(nil), a...), s)
		got, err := runCmd(candidate, input)
		if err != nil {
			fmt.Printf("test %d: candidate runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
