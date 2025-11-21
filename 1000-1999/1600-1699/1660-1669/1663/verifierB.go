package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runCandidate(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	if stderr.Len() > 0 {
		return "", fmt.Errorf("unexpected stderr output: %s", stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func solveRef(r int) int {
	thresholds := []int{-1000, 1200, 1400, 1600, 1900, 2100, 2300, 2400, 2600, 3000}
	for p := len(thresholds) - 1; p > 0; p-- {
		if r >= thresholds[p-1] {
			return thresholds[p]
		}
	}
	return thresholds[0]
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for r := -45; r <= 2999; r++ {
		input := fmt.Sprintf("%d\n", r)
		expect := solveRef(r)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("r=%d: runtime error: %v\n", r, err)
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Sscan(out, &got); err != nil {
			fmt.Printf("r=%d: failed to parse output %q: %v\n", r, out, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("r=%d: mismatch, expect %d got %d\n", r, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
