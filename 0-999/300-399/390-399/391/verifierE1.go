package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := ioutil.ReadFile("problemE1.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open problemE1.txt:", err)
		os.Exit(1)
	}
	cases := strings.Split(string(data), "\n---\n")
	total := len(cases)
	passed := 0
	for i, c := range cases {
		lines := strings.Split(strings.TrimSpace(c), "\n")
		if len(lines) == 0 {
			continue
		}
		input := strings.Join(lines[:len(lines)-1], "\n") + "\n"
		expected := strings.TrimSpace(lines[len(lines)-1])
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		result := strings.TrimSpace(string(out))
		if err != nil {
			fmt.Printf("Case %d: runtime error: %v\n", i+1, err)
			fmt.Printf("Output: %s\n", result)
			continue
		}
		if result == expected {
			passed++
		} else {
			fmt.Printf("Case %d failed: expected %s got %s\n", i+1, expected, result)
		}
	}
	fmt.Printf("%d/%d cases passed\n", passed, total)
}
