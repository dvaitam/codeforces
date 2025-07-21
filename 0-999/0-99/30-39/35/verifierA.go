package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(start int, shuffles [][2]int) int {
	pos := start
	for _, s := range shuffles {
		a, b := s[0], s[1]
		if pos == a {
			pos = b
		} else if pos == b {
			pos = a
		}
	}
	return pos
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	pairs := [][2]int{{1, 2}, {2, 1}, {1, 3}, {3, 1}, {2, 3}, {3, 2}}
	caseNum := 1
	for start := 1; start <= 3; start++ {
		for _, p1 := range pairs {
			for _, p2 := range pairs {
				for _, p3 := range pairs {
					sh := [][2]int{p1, p2, p3}
					input := fmt.Sprintf("%d\n%d %d\n%d %d\n%d %d\n", start, p1[0], p1[1], p2[0], p2[1], p3[0], p3[1])
					exp := fmt.Sprintf("%d", expected(start, sh))
					out, err := run(bin, input)
					if err != nil {
						fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum, err, input)
						os.Exit(1)
					}
					fields := strings.Fields(out)
					if len(fields) == 0 {
						fmt.Fprintf(os.Stderr, "case %d: empty output\ninput:\n%s", caseNum, input)
						os.Exit(1)
					}
					if fields[0] != exp {
						fmt.Fprintf(os.Stderr, "case %d: expected %s got %s\ninput:\n%s", caseNum, exp, fields[0], input)
						os.Exit(1)
					}
					caseNum++
				}
			}
		}
	}
	fmt.Printf("All %d tests passed\n", caseNum-1)
}
