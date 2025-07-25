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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
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
		if len(parts) < 3 {
			fmt.Fprintf(os.Stderr, "invalid line %d\n", idx)
			os.Exit(1)
		}
		mIdx := 1
		mVal := 0
		fmt.Sscan(parts[mIdx], &mVal)
		var input strings.Builder
		input.WriteString(parts[0] + " " + parts[1] + "\n")
		pos := 2
		for i := 0; i < mVal && pos+1 < len(parts); i++ {
			input.WriteString(fmt.Sprintf("%s %s\n", parts[pos], parts[pos+1]))
			pos += 2
		}
		expected := parts[len(parts)-1]
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		result := strings.TrimSpace(out.String())
		if err != nil {
			fmt.Printf("Case %d: runtime error: %v\n", idx, err)
			fmt.Printf("Output: %s\n", result)
			continue
		}
		if result == expected {
			passed++
		} else {
			fmt.Printf("Case %d failed: expected %s got %s\n", idx, expected, result)
		}
	}
	fmt.Printf("%d/%d cases passed\n", passed, idx)
}
