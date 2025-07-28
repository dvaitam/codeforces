package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runExe(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func prefix(s string) []int {
	n := len(s)
	pi := make([]int, n)
	for i := 1; i < n; i++ {
		j := pi[i-1]
		for j > 0 && s[i] != s[j] {
			j = pi[j-1]
		}
		if s[i] == s[j] {
			j++
		}
		pi[i] = j
	}
	return pi
}

func solveAll(s string, queries []string) []string {
	sb := []byte(s)
	piS := prefix(s)
	n := len(sb)
	results := make([]string, len(queries))
	for idx, t := range queries {
		tb := []byte(t)
		piT := make([]int, len(tb))
		j := piS[n-1]
		get := func(pos int) byte {
			if pos < n {
				return sb[pos]
			}
			return tb[pos-n]
		}
		for i := 0; i < len(tb); i++ {
			c := tb[i]
			for j > 0 && get(j) != c {
				if j <= n {
					j = piS[j-1]
				} else {
					j = piT[j-n-1]
				}
			}
			if get(j) == c {
				j++
			}
			piT[i] = j
		}
		strs := make([]string, len(piT))
		for i, v := range piT {
			strs[i] = strconv.Itoa(v)
		}
		results[idx] = strings.Join(strs, " ")
	}
	return results
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) < 2 {
		fmt.Println("bad test file")
		os.Exit(1)
	}
	s := lines[0]
	q, _ := strconv.Atoi(strings.TrimSpace(lines[1]))
	if len(lines) != q+2 {
		fmt.Println("bad test file")
		os.Exit(1)
	}
	queries := lines[2:]
	expected := solveAll(s, queries)

	var input strings.Builder
	input.WriteString(s)
	input.WriteByte('\n')
	fmt.Fprintf(&input, "%d\n", q)
	for _, t := range queries {
		input.WriteString(t)
		input.WriteByte('\n')
	}

	got, err := runExe(bin, input.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	outLines := strings.Split(strings.TrimSpace(got), "\n")
	if len(outLines) != q {
		fmt.Println("wrong number of output lines")
		os.Exit(1)
	}
	for i := 0; i < q; i++ {
		if strings.TrimSpace(outLines[i]) != expected[i] {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, expected[i], strings.TrimSpace(outLines[i]))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
