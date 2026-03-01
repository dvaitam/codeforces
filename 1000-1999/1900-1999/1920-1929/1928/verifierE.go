package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func genTests() []byte {
	rand.Seed(46)
	t := 110
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rand.Intn(30) + 1
		x := rand.Intn(60)
		y := rand.Intn(60) + 1
		s := rand.Intn(100)
		fmt.Fprintf(&sb, "%d %d %d %d\n", n, x, y, s)
	}
	return []byte(sb.String())
}

func runCmd(cmd *exec.Cmd, input []byte) (string, error) {
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func parseInts(input []byte) ([]int64, error) {
	fields := strings.Fields(string(input))
	vals := make([]int64, 0, len(fields))
	for _, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, err
		}
		vals = append(vals, v)
	}
	return vals, nil
}

func readToken(tokens []string, pos *int) (string, error) {
	if *pos >= len(tokens) {
		return "", io.EOF
	}
	tok := tokens[*pos]
	*pos++
	return tok, nil
}

func verifyAnswer(input []byte, out string, expectYes []bool) error {
	in, err := parseInts(input)
	if err != nil {
		return fmt.Errorf("invalid generated input: %w", err)
	}
	if len(in) == 0 {
		return fmt.Errorf("empty input")
	}
	t := int(in[0])
	idx := 1
	if len(expectYes) != t {
		return fmt.Errorf("internal error: expectYes length mismatch")
	}
	tokens := strings.Fields(out)
	pos := 0
	for tc := 0; tc < t; tc++ {
		if idx+3 >= len(in) {
			return fmt.Errorf("broken input on case %d", tc+1)
		}
		n := int(in[idx])
		x := in[idx+1]
		y := in[idx+2]
		s := in[idx+3]
		idx += 4

		tok, err := readToken(tokens, &pos)
		if err != nil {
			return fmt.Errorf("missing answer for case %d", tc+1)
		}
		up := strings.ToUpper(tok)
		switch up {
		case "NO":
			if expectYes[tc] {
				return fmt.Errorf("case %d: contestant answered NO but solution exists", tc+1)
			}
		case "YES":
			seq := make([]int64, n)
			for i := 0; i < n; i++ {
				vTok, err := readToken(tokens, &pos)
				if err != nil {
					return fmt.Errorf("case %d: expected %d sequence values", tc+1, n)
				}
				v, convErr := strconv.ParseInt(vTok, 10, 64)
				if convErr != nil {
					return fmt.Errorf("case %d: invalid integer %q", tc+1, vTok)
				}
				seq[i] = v
			}
			if seq[0] != x {
				return fmt.Errorf("case %d: a1=%d but expected x=%d", tc+1, seq[0], x)
			}
			sum := seq[0]
			for i := 1; i < n; i++ {
				if seq[i] != seq[i-1]+y && seq[i] != seq[i-1]%y {
					return fmt.Errorf("case %d: invalid transition at i=%d (%d -> %d)", tc+1, i+1, seq[i-1], seq[i])
				}
				sum += seq[i]
			}
			if sum != s {
				return fmt.Errorf("case %d: sum=%d but expected %d", tc+1, sum, s)
			}
			if !expectYes[tc] {
				return fmt.Errorf("case %d: contestant found sequence but reference says impossible", tc+1)
			}
		default:
			return fmt.Errorf("case %d: expected YES/NO, got %q", tc+1, tok)
		}
	}
	if pos != len(tokens) {
		return fmt.Errorf("extra output tokens after case %d", t)
	}
	return nil
}

func parseExpected(refOut string) ([]bool, error) {
	tokens := strings.Fields(refOut)
	expectYes := make([]bool, 0)
	for i := 0; i < len(tokens); i++ {
		up := strings.ToUpper(tokens[i])
		if up != "YES" && up != "NO" {
			continue
		}
		yes := up == "YES"
		expectYes = append(expectYes, yes)
		if yes {
			// Skip sequence values for this testcase in reference output.
			for i+1 < len(tokens) {
				next := strings.ToUpper(tokens[i+1])
				if next == "YES" || next == "NO" {
					break
				}
				i++
			}
		}
	}
	if len(expectYes) == 0 {
		return nil, fmt.Errorf("could not parse reference yes/no answers")
	}
	return expectYes, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	input := genTests()
	candCmd := exec.Command(os.Args[1])
	candOut, err := runCmd(candCmd, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate error: %v\n", err)
		os.Exit(1)
	}
	refCmd := exec.Command("go", "run", "1928E.go")
	refOut, err := runCmd(refCmd, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference error: %v\n", err)
		os.Exit(1)
	}
	expectYes, err := parseExpected(refOut)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference parse error: %v\n", err)
		os.Exit(1)
	}
	if err := verifyAnswer(input, candOut, expectYes); err != nil {
		fmt.Println("WA")
		fmt.Printf("reason: %v\n", err)
		fmt.Println("input:\n" + string(input))
		fmt.Println("expected:\n" + refOut)
		fmt.Println("got:\n" + candOut)
		os.Exit(1)
	}
	fmt.Println("OK")
}
