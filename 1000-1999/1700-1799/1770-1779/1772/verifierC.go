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

func solveC(n, k int) []int {
	s := []int{}
	present := make(map[int]bool)
	if n >= 1 {
		s = append(s, 1)
		present[1] = true
	}
	if n >= 2 {
		s = append(s, 2)
		present[2] = true
	}
	cnt := 2
	for i := 4; i <= k && len(s) < n; {
		s = append(s, i)
		present[i] = true
		cnt++
		i += cnt
	}
	if len(s) < n {
		for x := k; x >= 1 && len(s) < n; x-- {
			if !present[x] {
				s = append(s, x)
			}
		}
	}
	sort.Ints(s)
	return s
}

func expectedC(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var t int
	fmt.Sscanf(lines[0], "%d", &t)
	var out strings.Builder
	idx := 1
	for i := 0; i < t; i++ {
		var n, k int
		fmt.Sscanf(lines[idx], "%d %d", &n, &k)
		idx++
		ans := solveC(n, k)
		for j, v := range ans {
			if j > 0 {
				out.WriteByte(' ')
			}
			out.WriteString(fmt.Sprintf("%d", v))
		}
		if i+1 < t {
			out.WriteByte('\n')
		}
	}
	return strings.TrimSpace(out.String())
}

func genTestsC() []string {
	rand.Seed(3)
	tests := make([]string, 0, 100)
	for len(tests) < 100 {
		t := rand.Intn(20) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t))
		for i := 0; i < t; i++ {
			n := rand.Intn(39) + 2
			k := rand.Intn(40-n+1) + n
			sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		}
		tests = append(tests, sb.String())
	}
	return tests
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run verifierC.go <binary>\n")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsC()
	for i, tcase := range tests {
		want := expectedC(tcase)
		got, err := runBinary(bin, tcase)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("Test %d failed.\nInput:\n%s\nExpected:\n%s\nGot:\n%s\n", i+1, tcase, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
