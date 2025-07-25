package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveB(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var s, t string
	fmt.Fscan(reader, &s)
	fmt.Fscan(reader, &t)
	if len(s) != len(t) {
		return "impossible"
	}
	diff := 0
	for i := range s {
		if s[i] != t[i] {
			diff++
		}
	}
	if diff%2 == 1 {
		return "impossible"
	}
	res := make([]byte, len(s))
	pickFromT := false
	for i := range s {
		if s[i] == t[i] {
			res[i] = s[i]
		} else {
			if pickFromT {
				res[i] = t[i]
			} else {
				res[i] = s[i]
			}
			pickFromT = !pickFromT
		}
	}
	return string(res)
}

func genTests() []string {
	rand.Seed(2)
	tests := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(20) + 1
		var s, t strings.Builder
		for j := 0; j < n; j++ {
			s.WriteByte(byte('0' + rand.Intn(2)))
			t.WriteByte(byte('0' + rand.Intn(2)))
		}
		tests = append(tests, s.String()+"\n"+t.String())
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierB <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		reader := bufio.NewReader(strings.NewReader(tc))
		var s, t string
		fmt.Fscan(reader, &s)
		fmt.Fscan(reader, &t)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		diff := 0
		for k := range s {
			if s[k] != t[k] {
				diff++
			}
		}
		if diff%2 == 1 {
			if strings.TrimSpace(got) != "impossible" {
				fmt.Printf("test %d failed expected impossible, got %s\n", i+1, got)
				os.Exit(1)
			}
			continue
		}
		ans := strings.TrimSpace(got)
		if len(ans) != len(s) {
			fmt.Printf("test %d failed length mismatch\n", i+1)
			os.Exit(1)
		}
		d1 := 0
		d2 := 0
		for k := range ans {
			if ans[k] != '0' && ans[k] != '1' {
				fmt.Printf("test %d failed invalid character\n", i+1)
				os.Exit(1)
			}
			if ans[k] != s[k] {
				d1++
			}
			if ans[k] != t[k] {
				d2++
			}
		}
		if d1 != d2 || d1 != diff/2 {
			fmt.Printf("test %d failed distances mismatch\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
