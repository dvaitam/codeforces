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

var rawTestcases = []string{
	"7\n1 0 1\n1 1 1\n1 1 0\n0 1 0\n0 1 0\n1 0 0\n1 1 0",
	"6\n1 1 0\n1 1 1\n0 0 0\n1 0 1\n1 0 1\n0 0 0",
	"4\n0 1 0\n0 1 1\n0 1 1\n0 1 0",
	"10\n1 1 0\n1 1 0\n1 0 0\n0 0 1\n1 0 0\n0 0 0\n0 1 1\n0 0 1\n1 1 1\n1 0 1",
	"10\n0 1 1\n0 0 0\n1 0 0\n1 0 1\n1 0 0\n0 0 0\n0 0 0\n0 0 1\n0 1 0\n0 0 0",
	"3\n0 1 0\n0 0 1\n0 1 0",
	"4\n0 1 1\n1 0 0\n1 0 0\n1 0 1",
	"6\n1 0 0\n0 0 0\n1 1 0\n1 0 0\n1 1 1\n1 1 1",
	"3\n0 1 0\n1 0 1\n0 0 1",
	"6\n1 1 0\n1 1 1\n0 0 0\n1 0 0\n0 1 1\n1 0 1",
	"10\n1 0 0\n1 0 1\n0 1 1\n0 0 1\n1 1 1\n0 1 0\n0 0 0\n1 1 1\n0 0 0\n0 0 0",
	"2\n0 1 1\n0 0 1",
	"7\n0 0 1\n1 0 1\n0 1 0\n0 1 0\n0 0 0\n1 1 1\n0 1 1",
	"7\n1 0 0\n1 1 0\n0 0 1\n1 1 1\n0 1 0\n0 1 0\n0 1 1",
	"7\n0 1 0\n1 0 0\n1 0 0\n1 1 1\n1 1 0\n0 1 1\n1 1 0",
	"8\n0 1 1\n0 1 1\n0 0 1\n0 0 0\n0 0 0\n1 0 0\n1 1 1\n1 0 1",
	"2\n1 0 1\n0 1 0",
	"6\n0 0 0\n1 0 0\n1 1 1\n0 0 0\n0 0 0\n1 1 1",
	"9\n0 0 1\n0 1 1\n1 1 1\n1 0 0\n0 1 1\n1 0 0\n1 0 1\n1 1 0\n0 0 0",
	"4\n0 0 1\n0 1 1\n0 0 1\n1 1 0",
	"10\n0 0 1\n0 1 1\n1 1 0\n0 0 0\n1 1 1\n0 1 1\n1 1 1\n1 0 0\n0 0 1\n1 0 0",
	"8\n1 0 1\n0 0 1\n1 0 1\n1 1 0\n0 1 0\n0 0 0\n1 0 1\n0 1 1",
	"2\n0 1 1\n0 1 0",
	"7\n1 0 0\n0 1 1\n1 0 1\n1 0 0\n0 1 1\n0 0 0\n0 1 1",
	"2\n0 1 1\n0 0 1",
	"2\n0 1 0\n0 1 1",
	"7\n0 1 1\n1 0 1\n0 0 1\n0 1 1\n0 0 0\n1 0 0\n0 0 0",
	"2\n0 1 1\n0 0 0",
	"3\n1 0 1\n1 1 0\n0 0 1",
	"6\n1 1 1\n1 0 0\n0 1 0\n1 0 0\n0 1 0\n0 0 1",
	"6\n1 0 0\n0 1 0\n0 0 1\n1 1 0\n1 1 1\n0 0 1",
	"6\n0 1 1\n0 1 1\n1 0 1\n1 1 0\n0 0 1\n0 0 1",
	"7\n0 1 1\n0 0 1\n1 0 1\n1 0 1\n0 1 0\n1 1 0\n0 1 0",
	"5\n0 1 1\n0 1 0\n1 1 0\n0 1 0\n0 1 1",
	"9\n1 1 0\n1 1 1\n0 0 0\n0 1 1\n1 0 0\n0 0 1\n1 1 0\n0 0 1\n1 1 1",
	"8\n1 0 0\n0 0 1\n1 1 0\n1 1 1\n1 0 1\n1 1 1\n0 1 1\n0 1 0",
	"10\n0 0 1\n1 1 0\n0 0 1\n0 1 0\n0 0 1\n1 0 0\n1 0 0\n1 1 1\n1 0 1\n1 1 0",
	"8\n0 1 0\n0 0 1\n1 1 1\n1 1 0\n0 1 0\n0 1 0\n0 0 0\n1 0 0",
	"4\n1 0 1\n1 1 1\n0 1 1\n1 1 1",
	"4\n1 1 0\n0 0 0\n1 0 1\n0 0 0",
	"7\n0 1 0\n1 0 1\n0 0 1\n1 1 1\n1 0 0\n0 1 0\n0 1 0",
	"4\n1 1 1\n1 1 1\n1 0 0\n1 0 0",
	"8\n1 1 0\n1 1 1\n1 0 0\n0 0 0\n0 1 0\n1 0 1\n1 0 0\n0 1 0",
	"3\n1 1 1\n1 0 1\n0 1 1",
	"9\n1 1 1\n0 0 0\n0 0 1\n0 0 1\n0 1 0\n0 1 1\n0 1 0\n1 0 0\n0 1 1",
	"4\n1 1 1\n0 0 0\n1 1 0\n0 1 0",
	"3\n0 0 0\n1 0 0\n0 0 1",
	"6\n0 1 0\n0 1 1\n1 1 0\n0 0 1\n1 1 0\n1 0 1",
	"4\n1 1 1\n0 0 0\n1 0 1\n1 0 1",
	"6\n0 1 0\n1 0 0\n1 0 1\n0 0 0\n0 0 0\n1 0 0",
	"6\n1 0 0\n1 1 0\n1 1 0\n1 1 1\n0 1 1\n1 0 0",
	"7\n1 0 0\n0 0 0\n1 0 0\n0 0 0\n0 0 1\n0 1 0\n0 1 1",
	"7\n0 0 1\n0 1 1\n0 1 0\n1 0 1\n1 1 1\n1 0 1\n1 1 0",
	"9\n1 0 1\n1 0 1\n1 1 0\n0 1 0\n1 1 0\n1 1 1\n0 0 1\n0 1 0\n0 1 0",
	"10\n0 0 1\n1 1 1\n1 1 1\n1 0 0\n0 0 1\n0 1 0\n1 1 0\n0 1 1\n0 0 0\n0 0 1",
	"7\n1 1 1\n0 1 1\n1 1 1\n0 1 0\n0 0 1\n1 1 0\n1 0 0",
	"6\n0 0 1\n1 0 0\n0 1 0\n1 0 1\n0 1 0\n0 0 1",
	"2\n1 1 1\n1 1 1",
	"10\n1 0 1\n0 1 0\n1 0 1\n1 1 1\n0 1 1\n0 1 1\n0 1 1\n0 0 0\n1 1 0\n0 0 0",
	"9\n1 0 1\n0 1 1\n0 1 1\n0 0 1\n1 0 0\n0 0 1\n0 0 1\n0 1 0\n0 0 1",
	"4\n1 1 1\n1 1 0\n1 0 0\n0 0 1",
	"2\n0 1 0\n1 0 1",
	"3\n1 0 0\n1 0 0\n0 1 1",
	"10\n1 1 1\n1 0 0\n1 0 0\n0 0 1\n0 0 1\n0 1 1\n0 0 1\n0 1 0\n1 1 1\n1 1 1",
	"3\n0 1 1\n0 1 0\n1 0 0",
	"1\n1 0 0",
	"7\n0 1 1\n1 1 0\n0 1 1\n0 0 0\n1 1 0\n1 0 0\n1 0 0",
	"7\n1 0 0\n0 0 0\n0 1 0\n1 1 0\n0 0 1\n1 0 0\n1 0 0",
	"7\n1 1 0\n1 0 0\n0 1 0\n0 0 0\n1 0 0\n0 0 0\n0 1 0",
	"7\n1 1 0\n0 1 0\n1 0 0\n1 1 0\n1 1 0\n1 1 0\n1 0 1",
	"7\n0 1 1\n1 0 1\n0 1 0\n1 1 0\n0 0 1\n0 1 1\n1 0 1",
	"10\n1 0 0\n1 1 1\n0 1 0\n1 1 0\n1 1 0\n0 1 0\n1 1 1\n0 1 1\n0 1 1\n1 0 1",
	"3\n0 0 1\n0 0 1\n1 1 0",
	"2\n1 0 0\n0 0 0",
	"4\n0 0 1\n1 1 1\n0 1 1\n0 1 1",
	"9\n0 1 1\n0 1 1\n1 1 1\n1 1 0\n0 0 0\n0 0 0\n1 0 1\n0 0 0\n0 1 1",
	"6\n0 1 1\n1 1 0\n0 0 0\n0 0 1\n0 1 1\n1 0 0",
	"7\n0 1 0\n1 0 1\n1 0 1\n0 1 0\n1 1 1\n1 1 0\n0 1 0",
	"5\n0 1 0\n0 1 1\n1 0 1\n1 0 0\n1 0 1",
	"2\n0 0 0\n0 0 0",
	"7\n1 0 1\n1 1 0\n0 0 0\n1 1 1\n0 0 1\n1 1 0\n0 1 1",
	"4\n1 1 0\n0 1 0\n1 1 1\n1 0 0",
	"7\n1 1 0\n1 1 1\n1 1 1\n1 0 1\n1 0 1\n0 1 0\n0 1 1",
	"4\n0 1 1\n1 1 0\n1 1 1\n0 1 0",
	"2\n1 1 1\n0 1 0",
	"10\n0 1 0\n1 0 1\n0 0 0\n1 0 0\n0 0 1\n0 0 0\n0 0 1\n0 0 1\n0 0 0\n0 1 1",
	"8\n0 1 0\n1 0 1\n0 1 1\n1 1 0\n0 1 0\n0 1 0\n0 0 1\n0 1 0",
	"7\n0 0 1\n0 1 1\n0 1 1\n0 1 1\n1 0 1\n0 1 0\n1 1 1",
	"1\n1 0 1",
	"8\n0 1 0\n1 1 1\n0 0 1\n1 1 1\n0 0 0\n1 0 0\n0 1 0\n1 0 1",
	"7\n0 1 1\n0 0 1\n0 1 0\n1 0 1\n0 1 0\n1 0 0\n1 1 0",
	"1\n0 0 0",
	"4\n1 0 0\n1 1 1\n1 0 0\n1 1 1",
	"7\n1 1 1\n0 0 0\n0 0 1\n1 0 1\n0 0 1\n0 0 1\n1 0 0",
	"3\n1 1 1\n0 0 1\n0 0 1",
	"10\n1 1 0\n0 0 0\n1 0 1\n1 0 1\n1 0 1\n1 1 0\n1 0 0\n0 1 1\n0 1 1\n1 1 0",
	"2\n1 0 0\n0 0 1",
	"4\n0 1 1\n1 0 1\n1 1 0\n0 1 0",
	"3\n1 0 1\n1 0 1\n0 1 0",
	"5\n0 0 0\n1 0 0\n1 0 0\n0 0 1\n0 1 1",
}

// solve231A mirrors the reference solution from 231A.go using an in-memory reader.
func solve231A(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return "", err
	}
	ans := 0
	for i := 0; i < n; i++ {
		var a, b, c int
		if _, err := fmt.Fscan(reader, &a, &b, &c); err != nil {
			return "", err
		}
		if a+b+c >= 2 {
			ans++
		}
	}
	return strconv.Itoa(ans), nil
}

func runCase(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, tc := range rawTestcases {
		expected, err := solve231A(tc)
		if err != nil {
			fmt.Printf("case %d parse error: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runCase(bin, tc+"\n")
		if err != nil {
			fmt.Printf("case %d execution failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
