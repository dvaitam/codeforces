package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func bullsCows(a, b string) (int, int) {
	bulls := 0
	var seen [10]bool
	for i := 0; i < 4; i++ {
		if a[i] == b[i] {
			bulls++
		}
		seen[a[i]-'0'] = true
	}
	match := 0
	for i := 0; i < 4; i++ {
		if seen[b[i]-'0'] {
			match++
		}
	}
	return bulls, match - bulls
}

func enumerate(guesses []string, bulls, cows []int) string {
	var candidates []string
	digits := "0123456789"
	for i := 0; i < len(digits); i++ {
		for j := 0; j < len(digits); j++ {
			if j == i {
				continue
			}
			for k := 0; k < len(digits); k++ {
				if k == i || k == j {
					continue
				}
				for l := 0; l < len(digits); l++ {
					if l == i || l == j || l == k {
						continue
					}
					cand := []byte{digits[i], digits[j], digits[k], digits[l]}
					ok := true
					for idx, g := range guesses {
						b, c := bullsCows(string(cand), g)
						if b != bulls[idx] || c != cows[idx] {
							ok = false
							break
						}
					}
					if ok {
						candidates = append(candidates, string(cand))
						if len(candidates) > 1 {
							return "Need more data"
						}
					}
				}
			}
		}
	}
	if len(candidates) == 0 {
		return "Incorrect data"
	}
	return candidates[0]
}

func generateCases() []testCase {
	rand.Seed(3)
	cases := make([]testCase, 100)
	digits := "0123456789"
	allNums := []string{}
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if j == i {
				continue
			}
			for k := 0; k < 10; k++ {
				if k == i || k == j {
					continue
				}
				for l := 0; l < 10; l++ {
					if l == i || l == j || l == k {
						continue
					}
					allNums = append(allNums, string([]byte{digits[i], digits[j], digits[k], digits[l]}))
				}
			}
		}
	}
	for t := 0; t < 100; t++ {
		secret := allNums[rand.Intn(len(allNums))]
		n := rand.Intn(5) + 1
		var guesses []string
		var bVals []int
		var cVals []int
		for i := 0; i < n; i++ {
			g := allNums[rand.Intn(len(allNums))]
			b, c := bullsCows(secret, g)
			if rand.Intn(10) < 3 { // introduce errors sometimes
				if rand.Intn(2) == 0 {
					b = (b + rand.Intn(4)) % 5
					if b+c > 4 {
						c = 4 - b
					}
				} else {
					c = (c + rand.Intn(4)) % 5
					if b+c > 4 {
						b = 4 - c
					}
				}
			}
			guesses = append(guesses, g)
			bVals = append(bVals, b)
			cVals = append(cVals, c)
		}
		buf := bytes.Buffer{}
		fmt.Fprintln(&buf, n)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&buf, "%s %d %d\n", guesses[i], bVals[i], cVals[i])
		}
		expected := enumerate(guesses, bVals, cVals)
		cases[t] = testCase{input: buf.String(), expected: expected}
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierC.go <binary>")
		os.Exit(1)
	}
	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:%s\nactual:%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
