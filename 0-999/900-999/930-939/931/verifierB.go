package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(n, a, b int) string {
	rounds := 0
	for a != b {
		a = (a + 1) / 2
		b = (b + 1) / 2
		rounds++
	}
	total := 0
	for tmp := n; tmp > 1; tmp /= 2 {
		total++
	}
	if rounds == total {
		return "Final!"
	}
	return fmt.Sprintf("%d", rounds)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("could not open testcasesB.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("empty test file")
		os.Exit(1)
	}
	var t int
	fmt.Sscan(scan.Text(), &t)
	for i := 0; i < t; i++ {
		var n, a, b int
		if !scan.Scan() {
			fmt.Printf("case %d missing n\n", i+1)
			os.Exit(1)
		}
		fmt.Sscan(scan.Text(), &n)
		if !scan.Scan() {
			fmt.Printf("case %d missing a\n", i+1)
			os.Exit(1)
		}
		fmt.Sscan(scan.Text(), &a)
		if !scan.Scan() {
			fmt.Printf("case %d missing b\n", i+1)
			os.Exit(1)
		}
		fmt.Sscan(scan.Text(), &b)
		exp := expected(n, a, b)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d %d\n", n, a, b)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != exp {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
