package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesF.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open problemF.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	passed := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 4 {
			fmt.Println("invalid test case:", line)
			continue
		}
		n := parts[0]
		k := parts[1]
		q := parts[2]
		// number of initial values = n
		nVal := 0
		fmt.Sscan(n, &nVal)
		qVal := 0
		fmt.Sscan(q, &qVal)
		pos := 3
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%s %s %s\n", n, k, q))
		for i := 0; i < nVal && pos < len(parts); i++ {
			input.WriteString(parts[pos] + "\n")
			pos++
		}
		for i := 0; i < qVal && pos < len(parts); i++ {
			input.WriteString(parts[pos] + "\n")
			pos++
		}
		if pos >= len(parts) {
			fmt.Println("invalid test case:", line)
			continue
		}
		expectedTokens := strings.Split(parts[pos], ";")
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		resultLines := strings.Split(strings.TrimSpace(out.String()), "\n")
		if err != nil {
			fmt.Printf("Case %d: runtime error: %v\n", idx, err)
			fmt.Printf("Output: %s\n", out.String())
			continue
		}
		ok := true
		if len(resultLines) != len(expectedTokens) {
			ok = false
		} else {
			for i := range resultLines {
				if strings.TrimSpace(resultLines[i]) != expectedTokens[i] {
					ok = false
					break
				}
			}
		}
		if ok {
			passed++
		} else {
			fmt.Printf("Case %d failed: expected %v got %v\n", idx, expectedTokens, resultLines)
		}
	}
	fmt.Printf("%d/%d cases passed\n", passed, idx)
}
