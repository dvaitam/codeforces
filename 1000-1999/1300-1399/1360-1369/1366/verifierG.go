package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strings"
)

func apply(s string) (string, bool) {
	stack := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '.' {
			if len(stack) == 0 {
				return "", false
			}
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, c)
		}
	}
	return string(stack), true
}

func minDel(s, t string) int {
	n := len(s)
	best := n + 1
	for mask := 0; mask < 1<<n; mask++ {
		var sb strings.Builder
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				sb.WriteByte(s[i])
			}
		}
		res, ok := apply(sb.String())
		if ok && res == t {
			del := n - bits.OnesCount(uint(mask))
			if del < best {
				best = del
			}
		}
	}
	return best
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesG.txt")
	if err != nil {
		panic(err)
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
		fields := strings.Fields(line)
		if len(fields) != 2 {
			fmt.Printf("invalid test %d\n", idx)
			os.Exit(1)
		}
		s := fields[0]
		t := fields[1]
		expected := minDel(s, t)

		input := fmt.Sprintf("%s\n%s\n", s, t)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\n%s", idx, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(out.String())
		if outStr != fmt.Sprintf("%d", expected) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, expected, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
