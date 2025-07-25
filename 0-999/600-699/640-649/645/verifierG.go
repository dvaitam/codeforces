package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesG.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open problemG.txt:", err)
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
		if len(parts) < 5 {
			fmt.Println("invalid test case:", line)
			continue
		}
		nVal := 0
		fmt.Sscan(parts[0], &nVal)
		a := parts[1]
		var input strings.Builder
		input.WriteString(parts[0] + " " + a + "\n")
		pos := 2
		for i := 0; i < nVal && pos+1 < len(parts); i++ {
			input.WriteString(fmt.Sprintf("%s %s\n", parts[pos], parts[pos+1]))
			pos += 2
		}
		if pos >= len(parts) {
			fmt.Println("invalid test case:", line)
			continue
		}
		expectedStr := parts[pos]
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		resultStr := strings.TrimSpace(out.String())
		if err != nil {
			fmt.Printf("Case %d: runtime error: %v\n", idx, err)
			fmt.Printf("Output: %s\n", resultStr)
			continue
		}
		expected, _ := strconv.ParseFloat(expectedStr, 64)
		got, errp := strconv.ParseFloat(resultStr, 64)
		if errp != nil {
			fmt.Printf("Case %d failed: invalid output %s\n", idx, resultStr)
			continue
		}
		if math.Abs(expected-got) <= 1e-6*math.Max(1, math.Abs(expected)) {
			passed++
		} else {
			fmt.Printf("Case %d failed: expected %s got %s\n", idx, expectedStr, resultStr)
		}
	}
	fmt.Printf("%d/%d cases passed\n", passed, idx)
}
