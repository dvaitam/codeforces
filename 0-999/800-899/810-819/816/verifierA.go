package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func isPalindrome(h, m int) bool {
	return h/10 == m%10 && h%10 == m/10
}

func expectedMinutes(h, m int) int {
	ans := 0
	for {
		if isPalindrome(h, m) {
			return ans
		}
		ans++
		m++
		if m == 60 {
			m = 0
			h = (h + 1) % 24
		}
	}
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for h := 0; h < 24; h++ {
		for m := 0; m < 60; m++ {
			input := fmt.Sprintf("%02d:%02d\n", h, m)
			exp := expectedMinutes(h, m)
			gotStr, err := runBinary(bin, input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "runtime error on input %q: %v\n", input, err)
				os.Exit(1)
			}
			got, err := strconv.Atoi(strings.TrimSpace(gotStr))
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid output on input %q: %q\n", input, gotStr)
				os.Exit(1)
			}
			if got != exp {
				fmt.Fprintf(os.Stderr, "wrong answer on input %q: expected %d got %d\n", input, exp, got)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed!")
}
