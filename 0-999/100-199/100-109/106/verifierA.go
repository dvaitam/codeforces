package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcases = `D 6C QS
H 9D 9H
C QD 7S
D 8D 8D
S 6D TD
D 9H 6S
H KC 6S
C 6D 8C
C 6C 6C
D 6H QC
S 8H JH
C 8H 9H
S KS QH
C TH 6S
S TD 6D
H 9D KS
D QS KH
C 9S 9H
S AS AH
D QD 6D
C TH 6S
C JS TS
C JD 9C
D JH AS
H 6S AH
S 9H KH
D KH 9C
D TS AH
C AH 8D
D 7H AH
C 9S 9D
C 6C TS
S 6S JC
C KS JD
H AH 6D
H QD JC
C 7H JH
D TC 7S
C QH 8S
S 7C TC
H QC JC
D KH QH
S 6C AC
S QH 9S
C 9H 6D
S TH KH
C KC AH
C 9H KD
C AC QD
H 9D AD
C 7H AD
D AC KS
D AS 6C
D 6S KS
C JS QC
C 9C JC
H 6D TH
C TH KD
S JH TS
S QS QS
D AS 7D
D 9C 7H
H 9H 6S
C 8S 8H
H AD KD
D AH 8S
S AC TH
D 6S 9H
D AD TS
C JD JC
S 6D 8D
D JC 8H
S 6C QH
C TS JC
D JH 8H
S QS TD
D 6H TD
S AD 6S
C QC TC
C JH TS
D JS JC
H 9H KD
D KH QS
S 9D QD
D KH 8C
H KC 8C
H KD 7H
C TS TH
S 7H KS
C 9S 6S
H TS KS
C 8S AH
H AS 7C
C 6H 6H
S 7C 9S
H JS KD
S AC 9C
H QD AS
H 8S 8H
H TC QH`

func expected(trump, a, b string) string {
	order := "6789TJQKA"
	m := trump[0]
	if a[1] == m && b[1] != m {
		return "YES"
	}
	if a[1] == b[1] && strings.IndexByte(order, a[0]) > strings.IndexByte(order, b[0]) {
		return "YES"
	}
	return "NO"
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

	cases, err := loadTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		input := fmt.Sprintf("%s\n%s %s\n", tc.trump, tc.a, tc.b)
		want := strings.ToUpper(expected(tc.trump, tc.a, tc.b))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got = strings.ToUpper(strings.TrimSpace(got))
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

type testCase struct {
	trump string
	a     string
	b     string
}

func loadTestcases() ([]testCase, error) {
	scanner := bufio.NewScanner(strings.NewReader(testcases))
	var cases []testCase
	lineNum := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lineNum++
		parts := strings.Fields(line)
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid test line %d", lineNum)
		}
		cases = append(cases, testCase{
			trump: parts[0],
			a:     parts[1],
			b:     parts[2],
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("no testcases loaded")
	}
	return cases, nil
}
