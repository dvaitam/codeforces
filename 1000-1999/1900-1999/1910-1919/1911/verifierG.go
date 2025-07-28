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

func candidate(strs []string) string {
	base := []byte(strs[0])
	m := len(base)
	best := ""
	try := func(b []byte) {
		s := string(b)
		for _, t := range strs {
			diff := 0
			for i := 0; i < m; i++ {
				if b[i] != t[i] {
					diff++
					if diff > 1 {
						break
					}
				}
			}
			if diff > 1 {
				return
			}
		}
		if best == "" || s < best {
			best = s
		}
	}

	try(base)
	for i := 0; i < m; i++ {
		orig := base[i]
		for c := byte('a'); c <= 'z'; c++ {
			if c == orig {
				continue
			}
			base[i] = c
			try(base)
		}
		base[i] = orig
	}
	if best == "" {
		return "-1"
	}
	return best
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesG.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		k, _ := strconv.Atoi(parts[0])
		strs := parts[1 : 1+k]
		expected := candidate(strs)

		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", k))
		for i, s := range strs {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(s)
		}
		input.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(out.String())
		if outStr != expected {
			fmt.Printf("Test %d failed: expected %s got %s\n", idx, expected, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
