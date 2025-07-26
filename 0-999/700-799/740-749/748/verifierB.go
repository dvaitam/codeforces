package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const alpha = 26

func expected(s, t string) string {
	if len(s) != len(t) {
		return "-1"
	}
	f := make([]int, alpha)
	for i := range f {
		f[i] = -1
	}
	pairs := make([][2]byte, 0)
	for i := 0; i < len(s); i++ {
		a := int(s[i] - 'a')
		b := int(t[i] - 'a')
		if f[a] == -1 && f[b] == -1 {
			if a != b {
				f[a] = b
				f[b] = a
				if a < b {
					pairs = append(pairs, [2]byte{byte(a), byte(b)})
				}
			} else {
				f[a] = a
			}
		} else if f[a] != -1 {
			if f[a] != b {
				return "-1"
			}
			if f[b] == -1 {
				f[b] = a
				if a < b {
					pairs = append(pairs, [2]byte{byte(a), byte(b)})
				}
			} else if f[b] != a {
				return "-1"
			}
		} else { // f[a]==-1 && f[b]!=-1
			if f[b] != a {
				return "-1"
			}
			f[a] = b
			if a < b {
				pairs = append(pairs, [2]byte{byte(a), byte(b)})
			}
		}
	}
	// deduplicate pairs (since added multiple times maybe)
	seen := make(map[[2]byte]bool)
	uniq := make([][2]byte, 0)
	for _, p := range pairs {
		if !seen[p] {
			seen[p] = true
			uniq = append(uniq, p)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(uniq)))
	for _, p := range uniq {
		sb.WriteString(fmt.Sprintf("%c %c\n", 'a'+p[0], 'a'+p[1]))
	}
	return strings.TrimRight(sb.String(), "\n")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Printf("invalid testcase %d\n", idx)
			os.Exit(1)
		}
		s := parts[0]
		t := parts[1]
		expect := expected(s, t)
		input := fmt.Sprintf("%s\n%s\n", s, t)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expect {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", idx, input, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
