package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded testcases from testcasesB.txt.
const testcasesRaw = `A>B B<C A>C
B>C A>B A>C
B<C A>C A<B
A>C A<B B<C
A>B B<C A>C
A>B B>C A<C
A<B A<C B<C
B>C A>C A<B
A<C B<C A<B
A>B B<C A<C
A>C A<B B>C
B>C A>B A<C
B<C A>B A>C
A<B A<C B<C
A<C A>B B<C
A<B A<C B<C
A>B A<C B<C
A>B A>C B>C
A<B B<C A>C
A<B A<C B<C
B>C A<C A<B
A<C B<C A>B
A<C B<C A<B
B<C A<B A>C
A<C B<C A<B
A<C A>B B<C
A<C B>C A>B
B>C A<B A<C
A>C A<B B<C
A>B B>C A<C
A<C A<B B<C
B>C A>C A<B
A<B A<C B>C
A>B A>C B>C
A<C A<B B>C
A<C B>C A>B
A>B A>C B>C
A<B A>C B>C
A<B A<C B<C
A<B A<C B<C
A>B B<C A>C
B>C A<C A<B
A>C A<B B>C
A<B B>C A>C
A>C A<B B<C
A<B A<C B>C
A<C B<C A>B
A<C A<B B<C
B<C A>C A>B
A>B B<C A<C
A<B A>C B<C
A<C B>C A<B
A<B A<C B>C
A>B A<C B>C
B<C A>C A>B
B<C A<C A>B
A>C A>B B>C
B<C A>B A<C
B>C A<B A>C
A<C A>B B>C
A>C A<B B<C
A>B B>C A>C
A<B B>C A<C
B>C A<C A>B
B>C A>B A<C
A>C B<C A<B
B<C A<B A<C
A<B B>C A>C
B<C A<B A>C
B>C A<C A<B
A<B A<C B<C
A<B B>C A<C
A<C A<B B<C
A<B B>C A<C
B>C A<C A<B
A>B B>C A>C
A<B B<C A<C
B<C A>B A<C
A>C B>C A>B
A<C B<C A<B
B>C A>B A<C
A<B A>C B>C
A<B A<C B<C
A<B B>C A<C
A<C A<B B<C
B>C A>C A>B
A>B B>C A>C
B<C A<C A<B
A>B B<C A<C
A>C A<B B<C
A>B A>C B<C
B<C A>C A>B
A<C A<B B<C
B>C A<C A<B
A>B A<C B<C
B>C A<C A<B
B<C A<C A<B
A>B B<C A>C
A<B B<C A<C
B<C A<C A<B`

func parseTestcases() [][]string {
	lines := strings.Split(testcasesRaw, "\n")
	var res [][]string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		res = append(res, strings.Fields(line))
	}
	return res
}

// Embedded solver logic from 47B.go.
func solve(comps []string) string {
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

	cases := parseTestcases()

	for idx, comps := range cases {
		expect := solve(comps)
		input := strings.Join(comps, "\n") + "\n"

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
