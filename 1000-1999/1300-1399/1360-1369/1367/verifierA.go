package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesRaw = `100
baabbccbbbbbbb
caaccaabbaaa
cccccaabba
cbbb
abbbbbbccccaaccbbb
baaccaaaaccbbccccc
cb
accbbccaaaac
aaaccbba
bccb
bccb
cbbc
cccccbbb
cbbb
baaaaaaa
baaaaccaaa
ac
cbbccccbbccaaaaccc
cbbbbbbccccccb
bcca
cccbbaaaaaaccbba
baabbbba
acca
cc
cccaaaaaaccaacccca
abbaaaaccaaaaa
baac
ca
bccaabbaaaaaaccbbb
aaaccbbaaccaac
abbbbccbbccaac
accaaaab
baaccbbccaaaabbccb
bccbbbbccbbaacccca
caabbccaaccbbaaa
bccbbccbbcccccca
bccbbccaaa
cbbaaaaa
bccccccbbaabbccc
cccaaaabbaabbc
bccbbc
ab
bbbaabbaaccc
caaa
cbbbbaaaaaacca
caaaaaaccccaabbbbc
abbbbc
abbb
baac
cccbbaaaabbaaaaaaa
cbbbbccaac
cccbbccbbbbccaaa
cbbaaaaaabbbbb
caabbccaaaab
accbbb
caabbaabbccaaa
accccaabbb
bbbaaaaccbbb
baabbaaccbbb
bb
accccb
aaaa
abbaaaab
cbbbbbbccccccaabba
abbccaabbaab
acca
bccaaaabbbbbbaacca
caaaaaab
bccaaaaccbbaac
cccbbccbbccbbb
cccaaccccaaaabbc
baaccaabbcca
cbbccccbbaaaac
aa
abbaab
aaaccbbbbccb
caaccaaccccccaaccb
bccccbbb
cccaaaaccaacca
cccbbbbaab
bbbbbbbaaa
abbccb
ab
abbccccaacccca
baaccbbbbbba
baaa
cc
cbbaac
cbbaabbbbaaaaccccb
caaccccccbba
bccccaaccaaaab
baabbbbccaaaaaab
acccccccca
accccccc
cbbaaaabbccbba
baaa
accaabbbbccbba
bccbbaabbaaaabbc
bb
aaabba`

// Embedded reference logic from 1367A.go.
func solve(b string) string {
	if len(b) == 0 {
		return ""
	}
	a := make([]byte, 0, len(b)/2+1)
	a = append(a, b[0])
	for i := 1; i < len(b); i += 2 {
		a = append(a, b[i])
	}
	return string(a)
}

func parseTestcases(raw string) ([]string, error) {
	sc := bufio.NewScanner(strings.NewReader(raw))
	sc.Split(bufio.ScanLines)
	if !sc.Scan() {
		return nil, fmt.Errorf("empty test data")
	}
	var t int
	if _, err := fmt.Sscan(sc.Text(), &t); err != nil {
		return nil, fmt.Errorf("invalid test count")
	}
	tests := make([]string, 0, t)
	for sc.Scan() {
		tests = append(tests, strings.TrimSpace(sc.Text()))
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if len(tests) != t {
		return nil, fmt.Errorf("expected %d tests, got %d", t, len(tests))
	}
	return tests, nil
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse tests:", err)
		os.Exit(1)
	}

	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", len(tests)))
	expected := make([]string, 0, len(tests))
	for _, b := range tests {
		input.WriteString(b)
		input.WriteByte('\n')
		expected = append(expected, solve(b))
	}

	got, err := run(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "candidate failed:", err)
		os.Exit(1)
	}
	fields := strings.Fields(got)
	if len(fields) != len(expected) {
		fmt.Fprintf(os.Stderr, "wrong number of outputs: expected %d got %d\n", len(expected), len(fields))
		os.Exit(1)
	}
	for i := range expected {
		if fields[i] != expected[i] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expected[i], fields[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
