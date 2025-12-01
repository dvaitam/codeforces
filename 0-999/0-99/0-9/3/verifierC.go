package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded copy of testcasesC.txt.
const testcasesCData = `
OOX O.O OOO
O.X .XO XX.
O.. .XO X.X
.OO .XO OO.
.X. OO. OX.
XX. O.. .X.
OOX .O. XX.
XXX .OX XO.
OXO .O. X.O
.X. ..O OX.
OO. XOX XXX
..O OXX .XX
XX. ..O ..O
.XX ..O .OO
O.. .OX O.X
O.. OXX X.O
X.X OXO OXX
X.X X.. ...
XXX .X. .XO
XOX X.X XX.
XOX .X. X.O
.XO XXX .OO
OXX .OX .X.
OXO O.O .X.
.XX .XX O.O
X.O .XX O.O
..O .OO .OX
..X O.X O.X
.OX XOO .O.
O.. .X. OO.
O.X X.X .OX
XX. OO. ..O
XO. .O. .XX
OXO .XO .O.
.XX OOO OXO
X.. X.X XO.
OOX XX. X..
.XX X.. XOO
.XX OO. XXO
OXO X.. ..O
XXO XXX X.O
.OO .X. ...
O.. O.O O.X
XO. OXX XOO
OO. XO. XXO
XX. OOO .XO
XO. XXO X..
XOO OOO XX.
OOO OXO X.O
OXO O.. XX.
.O. XXX X.X
XOX XO. .OO
O.. .XO XOX
O.X OXO XX.
X.O .XX OOO
X.X X.X XXO
OO. XX. OX.
O.. O.O .OO
O.. X.. XXO
..O OX. XO.
XO. O.. OXX
.XX XXX OXO
OXX .OO .O.
.X. X.. .X.
OXO ..O OO.
.XX .X. XO.
.OO XOO OOO
OXX .XX O.O
XXO OXO ..X
..X OOX .OO
OXX OXX X..
X.O X.. .OX
OOX X.. O..
O.X ... OXO
O.. X.X XOO
OXO O.X XXO
OX. ... XXX
... .O. OXX
O.O XXO XXO
X.X OO. OXO
O.O XOX XO.
XOO XXX OX.
XXX X.. .X.
XOO .X. XXX
OXO O.O XXX
OOO OO. .OX
.XX ..O XOX
XXO .XX XOO
O.X X.X OX.
XXO O.O .X.
OO. OXX O.O
..X OOX O.O
OXO OOX XXO
X.. .X. OOX
OOX XOO .XO
O.. XOX OX.
O.O XXO X.O
..X O.O X.O
X.O O.X XO.
XXO O.. .OO
`

type testCase struct {
	board [3]string
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesCData), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 3 {
			return nil, fmt.Errorf("line %d: expected 3 rows", idx+1)
		}
		cases = append(cases, testCase{board: [3]string{parts[0], parts[1], parts[2]}})
	}
	return cases, nil
}

func win(board [3]string, ch byte) bool {
	for i := 0; i < 3; i++ {
		if board[i][0] == ch && board[i][1] == ch && board[i][2] == ch {
			return true
		}
		if board[0][i] == ch && board[1][i] == ch && board[2][i] == ch {
			return true
		}
	}
	if board[0][0] == ch && board[1][1] == ch && board[2][2] == ch {
		return true
	}
	if board[0][2] == ch && board[1][1] == ch && board[2][0] == ch {
		return true
	}
	return false
}

// solve mirrors 3C.go verdict logic.
func solve(board [3]string) string {
	xCnt, oCnt := 0, 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			switch board[i][j] {
			case 'X':
				xCnt++
			case '0':
				oCnt++
			default:
				// '.' or other handled by main validation in original; assume testcases valid
			}
		}
	}
	if !(xCnt == oCnt || xCnt == oCnt+1) {
		return "illegal"
	}
	xWin := win(board, 'X')
	oWin := win(board, '0')
	if xWin && oWin {
		return "illegal"
	}
	if xWin {
		if xCnt == oCnt+1 {
			return "the first player won"
		}
		return "illegal"
	}
	if oWin {
		if xCnt == oCnt {
			return "the second player won"
		}
		return "illegal"
	}
	if xCnt+oCnt == 9 {
		return "draw"
	}
	if xCnt == oCnt {
		return "first"
	}
	return "second"
}

func runCandidate(bin string, tc testCase) (string, error) {
	input := fmt.Sprintf("%s\n%s\n%s\n", tc.board[0], tc.board[1], tc.board[2])
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		expect := solve(tc.board)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
