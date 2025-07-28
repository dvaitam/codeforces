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

func isPalindrome(s string) bool {
	i, j := 0, len(s)-1
	for i < j {
		if s[i] != s[j] {
			return false
		}
		i++
		j--
	}
	return true
}

func allSame(s string) bool {
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return false
		}
	}
	return true
}

func solveCase(s string) (bool, []string) {
	n := len(s)
	if n == 1 {
		return false, nil
	}
	if !isPalindrome(s) {
		return true, []string{s}
	}
	if allSame(s) {
		return false, nil
	}
	if n%2 == 0 {
		for _, j := range []int{2, 3} {
			if j <= n-2 {
				left, right := s[:j], s[j:]
				if !isPalindrome(left) && !isPalindrome(right) {
					return true, []string{left, right}
				}
			}
		}
		for j := 2; j <= n-2 && j <= 6; j++ {
			left, right := s[:j], s[j:]
			if !isPalindrome(left) && !isPalindrome(right) {
				return true, []string{left, right}
			}
		}
		return false, nil
	}
	if n >= 2 {
		x := s[0]
		y := s[1]
		if x != y {
			alt := true
			for i := 0; i < n; i++ {
				if i%2 == 0 && s[i] != x {
					alt = false
					break
				}
				if i%2 == 1 && s[i] != y {
					alt = false
					break
				}
			}
			if alt {
				return false, nil
			}
		}
	}
	outerChar := s[0]
	outerSame := true
	mid := n / 2
	for i := 0; i < n; i++ {
		if i == mid {
			continue
		}
		if s[i] != outerChar {
			outerSame = false
			break
		}
	}
	if outerSame && s[mid] != outerChar {
		return false, nil
	}
	for _, j := range []int{2, 3} {
		if j <= n-2 {
			left, right := s[:j], s[j:]
			if !isPalindrome(left) && !isPalindrome(right) {
				return true, []string{left, right}
			}
		}
	}
	for j := 2; j <= n-2 && j <= 6; j++ {
		left, right := s[:j], s[j:]
		if !isPalindrome(left) && !isPalindrome(right) {
			return true, []string{left, right}
		}
	}
	return false, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	file, err := os.Open("testcasesE.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if s == "" {
			continue
		}
		idx++
		ok, parts := solveCase(s)
		var exp string
		if !ok {
			exp = "NO"
		} else {
			exp = "YES\n" + strconv.Itoa(len(parts)) + "\n" + strings.Join(parts, " ")
		}
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(s + "\n")

		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input.String())
		var outBuf bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(outBuf.String())
		if outStr != exp {
			fmt.Printf("Test %d failed: expected %q got %q\n", idx, exp, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
