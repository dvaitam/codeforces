package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func expected(tests []string) string {
	var results []string
	for _, s := range tests {
		var pairs []string
		for a := 1; a <= 12; a++ {
			if 12%a != 0 {
				continue
			}
			b := 12 / a
			found := false
			for j := 0; j < b; j++ {
				allX := true
				for i := 0; i < a; i++ {
					idx := i*b + j
					if s[idx] != 'X' {
						allX = false
						break
					}
				}
				if allX {
					found = true
					break
				}
			}
			if found {
				pairs = append(pairs, fmt.Sprintf("%dx%d", a, b))
			}
		}
		line := fmt.Sprintf("%d", len(pairs))
		if len(pairs) > 0 {
			line += " " + strings.Join(pairs, " ")
		}
		results = append(results, line)
	}
	return strings.Join(results, "\n")
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tcase := 0; tcase < 100; tcase++ {
		t := rng.Intn(4) + 1
		tests := make([]string, t)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t))
		for i := 0; i < t; i++ {
			var str strings.Builder
			for j := 0; j < 12; j++ {
				if rng.Intn(2) == 0 {
					str.WriteByte('X')
				} else {
					str.WriteByte('O')
				}
			}
			s := str.String()
			tests[i] = s
			sb.WriteString(s)
			sb.WriteByte('\n')
		}
		input := sb.String()
		expectedOut := expected(tests)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tcase+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expectedOut) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", tcase+1, expectedOut, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
