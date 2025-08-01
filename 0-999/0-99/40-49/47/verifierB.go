package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(comps []string) string {
	perms := []string{"ABC", "ACB", "BAC", "BCA", "CAB", "CBA"}
	for _, perm := range perms {
		ok := true
		for _, comp := range comps {
			if len(comp) < 3 {
				ok = false
				break
			}
			lhs := strings.IndexRune(perm, rune(comp[0]))
			rhs := strings.IndexRune(perm, rune(comp[2]))
			if comp[1] == '<' {
				if lhs >= rhs {
					ok = false
					break
				}
			} else if comp[1] == '>' {
				if lhs <= rhs {
					ok = false
					break
				}
			} else {
				ok = false
				break
			}
		}
		if ok {
			return perm
		}
	}
	return "Impossible"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		comps := strings.Fields(line)
		expect := expected(comps)
		input := strings.Join(comps, "\n") + "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
