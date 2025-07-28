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

func solveB(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(in, &t)
	var out strings.Builder
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)
		substr1 := make(map[string]struct{})
		substr2 := make(map[string]struct{})
		substr3 := make(map[string]struct{})
		for i := 0; i < n; i++ {
			substr1[s[i:i+1]] = struct{}{}
			if i+1 < n {
				substr2[s[i:i+2]] = struct{}{}
			}
			if i+2 < n {
				substr3[s[i:i+3]] = struct{}{}
			}
		}
		found := ""
		for c := byte('a'); c <= 'z' && found == ""; c++ {
			str := string([]byte{c})
			if _, ok := substr1[str]; !ok {
				found = str
			}
		}
		if found == "" {
			for c1 := byte('a'); c1 <= 'z' && found == ""; c1++ {
				for c2 := byte('a'); c2 <= 'z'; c2++ {
					str := string([]byte{c1, c2})
					if _, ok := substr2[str]; !ok {
						found = str
						break
					}
				}
			}
		}
		if found == "" {
			for c1 := byte('a'); c1 <= 'z' && found == ""; c1++ {
				for c2 := byte('a'); c2 <= 'z' && found == ""; c2++ {
					for c3 := byte('a'); c3 <= 'z'; c3++ {
						str := string([]byte{c1, c2, c3})
						if _, ok := substr3[str]; !ok {
							found = str
							break
						}
					}
				}
			}
		}
		out.WriteString(found)
		out.WriteByte('\n')
	}
	return strings.TrimSpace(out.String())
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []string {
	r := rand.New(rand.NewSource(2))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := r.Intn(20) + 1
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			sb.WriteByte(byte('a' + r.Intn(26)))
		}
		sb.WriteByte('\n')
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		expected := solveB(t)
		got, err := runBinary(bin, t)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("test %d failed. input: %sexpected %s got %s\n", i+1, t, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
