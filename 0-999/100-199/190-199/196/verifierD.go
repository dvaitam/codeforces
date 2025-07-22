package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func nextString(d int, s string) string {
	n := len(s)
	if d == 1 {
		return "Impossible"
	}
	t := []byte(s)
	ok := func(pos int) bool {
		if pos-d+1 >= 0 && t[pos] == t[pos-d+1] {
			return false
		}
		if pos-d >= 0 && t[pos] == t[pos-d] {
			return false
		}
		return true
	}
	orig := []byte(s)
	for i := n - 1; i >= 0; i-- {
		for c := orig[i] + 1; c <= 'z'; c++ {
			t[i] = c
			if !ok(i) {
				continue
			}
			viable := true
			for j := i + 1; j < n && viable; j++ {
				placed := false
				for c2 := byte('a'); c2 <= 'z'; c2++ {
					t[j] = c2
					if ok(j) {
						placed = true
						break
					}
				}
				if !placed {
					viable = false
				}
			}
			if viable {
				return string(t)
			}
		}
		t[i] = orig[i]
	}
	return "Impossible"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanLines)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	var t int
	fmt.Sscanf(scan.Text(), "%d", &t)
	dvals := make([]int, t)
	strs := make([]string, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		fmt.Sscanf(scan.Text(), "%d", &dvals[i])
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		strs[i] = scan.Text()
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got := outScan.Text()
		expected := nextString(dvals[i], strs[i])
		if got != expected {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
