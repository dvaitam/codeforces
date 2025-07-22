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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
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

func parse(tokens []string) (string, bool) {
	type frame struct{ state int }
	frames := make([]frame, 0, len(tokens))
	var out strings.Builder
	for i, tok := range tokens {
		if tok == "int" {
			out.WriteString("int")
			for len(frames) > 0 && frames[len(frames)-1].state == 1 {
				out.WriteByte('>')
				frames = frames[:len(frames)-1]
			}
			if len(frames) > 0 && frames[len(frames)-1].state == 0 {
				out.WriteByte(',')
				frames[len(frames)-1].state = 1
			} else if len(frames) == 0 && i != len(tokens)-1 {
				return "Error occurred", false
			}
		} else if tok == "pair" {
			out.WriteString("pair<")
			frames = append(frames, frame{state: 0})
		} else {
			return "Error occurred", false
		}
	}
	if len(frames) != 0 {
		return "Error occurred", false
	}
	return out.String(), true
}

func genValid(depth int) []string {
	if depth > 3 || rand.Intn(2) == 0 {
		return []string{"int"}
	}
	left := genValid(depth + 1)
	right := genValid(depth + 1)
	return append(append([]string{"pair"}, left...), right...)
}

func generateCases() []testCase {
	rand.Seed(3)
	cases := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		valid := rand.Intn(2) == 0
		if valid {
			tokens := genValid(0)
			n := 0
			for _, t := range tokens {
				if t == "int" {
					n++
				}
			}
			buf := bytes.Buffer{}
			fmt.Fprintf(&buf, "%d\n", n)
			for j, t := range tokens {
				if j > 0 {
					fmt.Fprint(&buf, " ")
				}
				fmt.Fprint(&buf, t)
			}
			buf.WriteByte('\n')
			res, _ := parse(tokens)
			cases[i] = testCase{input: buf.String(), expected: res}
		} else {
			n := rand.Intn(4) + 1
			total := n + rand.Intn(3) // maybe extra pair tokens
			tokens := make([]string, total)
			intsPlaced := 0
			for j := 0; j < total; j++ {
				if intsPlaced < n && rand.Intn(2) == 0 {
					tokens[j] = "int"
					intsPlaced++
				} else {
					tokens[j] = "pair"
				}
			}
			for intsPlaced < n {
				tokens[rand.Intn(total)] = "int"
				intsPlaced++
			}
			buf := bytes.Buffer{}
			fmt.Fprintf(&buf, "%d\n", n)
			for j, t := range tokens {
				if j > 0 {
					fmt.Fprint(&buf, " ")
				}
				fmt.Fprint(&buf, t)
			}
			buf.WriteByte('\n')
			res, _ := parse(tokens)
			cases[i] = testCase{input: buf.String(), expected: res}
		}
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
