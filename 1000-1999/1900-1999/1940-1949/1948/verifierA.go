package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// countSpecial counts the number of special characters in the string as defined
// in problem A.
func countSpecial(s string) int {
	b := []byte(s)
	n := len(b)
	cnt := 0
	for i := 0; i < n; i++ {
		left := i > 0 && b[i] == b[i-1]
		right := i+1 < n && b[i] == b[i+1]
		if (left && !right) || (!left && right) {
			cnt++
		}
	}
	return cnt
}

func validString(s string, n int) bool {
	if len(s) > 200 {
		return false
	}
	for _, ch := range s {
		if ch < 'A' || ch > 'Z' {
			return false
		}
	}
	return countSpecial(s) == n
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	// generate 100 test cases: n=1..50 repeated twice
	cases := make([]int, 100)
	for i := 0; i < 50; i++ {
		cases[i] = i + 1
		cases[i+50] = i + 1
	}

	var input bytes.Buffer
	fmt.Fprintf(&input, "%d\n", len(cases))
	for _, n := range cases {
		fmt.Fprintf(&input, "%d\n", n)
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to run binary:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	for idx, n := range cases {
		if !scanner.Scan() {
			fmt.Printf("missing output for case %d\n", idx+1)
			os.Exit(1)
		}
		line := scanner.Text()
		if n%2 == 1 {
			if line != "NO" {
				fmt.Printf("case %d: expected NO, got %s\n", idx+1, line)
				os.Exit(1)
			}
			continue
		}
		if line != "YES" {
			fmt.Printf("case %d: expected YES, got %s\n", idx+1, line)
			os.Exit(1)
		}
		if !scanner.Scan() {
			fmt.Printf("case %d: expected string line\n", idx+1)
			os.Exit(1)
		}
		ans := scanner.Text()
		if !validString(ans, n) {
			fmt.Printf("case %d: invalid string %q for n=%d\n", idx+1, ans, n)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Println("warning: extra output detected")
	}
	fmt.Println("All tests passed!")
}
