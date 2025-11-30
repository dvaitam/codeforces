package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n     int
	cards []string
}

// Embedded testcases (previously stored in testcasesB.txt) to keep the verifier self contained.
const rawTestcasesB = `
2 JS 6S
4 AC 9C AD 3C
1 8C
5 AS KC 6D JS 7S
1 2S
4 QD 8S TD AC
4 TD 7D QD AC
3 2C TS 4H
1 KH
6 KC TD 6H JC TC JS
4 5C 8D 7H 3C
6 TS AD TC 7C KS 9S
3 KC QD 4D
1 AD
5 TD 8H JH 9H QS
4 AD TD 8S 9H
5 TD TC 9H 8H 2H
4 JS AD QD TD
7 3H 2S 3S 9S AH 5H 3D
3 6S 4D 6D
6 6H 9H 9C 3S 6C 7C
7 5H 3H KD JC 2D 2C 4S
6 4C KC TD QC 5S 8H
6 QC 2H 4D 2H 3S 6H
6 4C JH 4S TS JD JC
2 AS 8D
3 3D JC JD
4 3C 6C 2H JC
3 2D 5H AD
3 8D 6S 8H
7 QC AD 3S 3D 4D TD 6H
5 TH 7H 7S 6D JC
2 JS 7S
4 3C AD 4H 3C
1 JD
5 3H 7H JS 9H 3S
7 6S JS 3C 3S 5D AC 4S
4 4D 4S 8C AH
5 6C 7S 5H 2S 2H
6 JH 9C 7C 3S 7C 3H
2 AC QH
3 4D 6D 5H
1 6S
7 9S QH 5C 6S 7D 7H 5H
1 TS
2 5S AD
4 3H TS KS 2S
3 AH 9C 4S
5 AH 3D 4D 4H 6S
6 TH 4D 4S AH JD 4H
4 TD 2D 6S QC
7 8H TC TC 2C 7D 6C 2C
5 2S KH JD JD 4H
7 6C JC 4S 5C 2D TH TC
6 QD 5H 9C 5C 7H QD
1 3H
2 TD 6H
6 6H 4C JS 3C 4D 6C
2 JS 9C
6 QH 8D TS TS AH QS
3 KS 4S 9D
7 8C 8D 7C 4C 5S 8C 3H
3 5C KS 5C
5 2S QD 6D 4H 4D
3 6H QC AD
5 7C 8S AD JC 5H
7 3S 3S TH QD 3H JH 8H
7 TH 2S 9C 7H TC 7C 3C
4 5S 6D 9C KH
6 4C JD 7S QC JC 8H
7 JS 9D QH QS 8D QC AH
7 4S AS 7H AC QH 4C 6C
2 9S 6S
6 JC 3H 3C 2D TD KS
4 QH JH 5D 5H
3 3S KH 9S
2 6H 7D
4 TC 4C AH JH
6 5H JD QS JC 7C AD
7 6D 3D JC JD JH 9D 4D
6 9H 6C 5S KD KH 3S
2 8H 9S
2 2S AS
7 5S 9C 7H 3D 3D 8D 9C
4 AD 5D 6C TC
2 9H 7C
5 3D 3S 2S 9H 8H
2 8D AD
7 2S 8D QS JC 6D 3C QH
1 2S
5 4S 6S 8S 5S 9D
6 6D QC 8H QH 6D 5S
5 AD 7C JS 7C TD
6 TC QS KH KS 6D 3D
1 5C
7 2S QS TC TH 3H 2D TS
4 QD 8C 2H 3H
7 7S 6S 8S KH 7D 6C AS
7 QH 3C 5D 7H TC JC 3D
6 9S 6D 5H 8H 3C 7D
`

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(rawTestcasesB, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 1 {
			return nil, fmt.Errorf("line %d: empty", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d: expected %d cards got %d", idx+1, n, len(fields)-1)
		}
		cases = append(cases, testCase{n: n, cards: fields[1:]})
	}
	return cases, nil
}

// match and solve208BCase are lifted from 208B.go so the verifier needs no external oracle.
func match(a, b string) bool {
	return a[0] == b[0] || a[1] == b[1]
}

func solve208BCase(cards []string) string {
	piles := append([]string(nil), cards...)
	for {
		moved := false
		for i := range piles {
			if i >= 3 && match(piles[i], piles[i-3]) {
				piles[i-3] = piles[i]
				piles = append(piles[:i], piles[i+1:]...)
				moved = true
				break
			}
			if i >= 1 && match(piles[i], piles[i-1]) {
				piles[i-1] = piles[i]
				piles = append(piles[:i], piles[i+1:]...)
				moved = true
				break
			}
		}
		if !moved {
			break
		}
	}
	if len(piles) == 1 {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		expect := solve208BCase(tc.cards)
		var sb strings.Builder
		sb.WriteString(fmt.Sprint(tc.n))
		sb.WriteByte('\n')
		for i, c := range tc.cards {
			sb.WriteString(c)
			if i+1 < len(tc.cards) {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(sb.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		got := strings.ToUpper(strings.TrimSpace(out.String()))
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
