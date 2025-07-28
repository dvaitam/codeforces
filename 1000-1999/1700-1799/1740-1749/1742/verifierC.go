package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	board [8]string
	ans   string
}

func randomSet(n, k int) []int {
	perm := rand.Perm(n)
	return perm[:k]
}

func generateCases() []testCase {
	rand.Seed(3)
	cases := make([]testCase, 100)
	for idx := range cases {
		lastRed := rand.Intn(2) == 0
		numRows := rand.Intn(8) + 1
		numCols := rand.Intn(8) + 1
		rows := randomSet(8, numRows)
		cols := randomSet(8, numCols)
		var grid [8][8]byte
		if lastRed {
			for _, c := range cols {
				for r := 0; r < 8; r++ {
					grid[r][c] = 'B'
				}
			}
			for _, r := range rows {
				for c := 0; c < 8; c++ {
					grid[r][c] = 'R'
				}
			}
		} else {
			for _, r := range rows {
				for c := 0; c < 8; c++ {
					grid[r][c] = 'R'
				}
			}
			for _, c := range cols {
				for r := 0; r < 8; r++ {
					grid[r][c] = 'B'
				}
			}
		}
		var board [8]string
		for r := 0; r < 8; r++ {
			line := make([]byte, 8)
			for c := 0; c < 8; c++ {
				ch := grid[r][c]
				if ch == 0 {
					ch = '.'
				}
				line[c] = ch
			}
			board[r] = string(line)
		}
		ans := "B"
		if lastRed {
			ans = "R"
		}
		cases[idx] = testCase{board: board, ans: ans}
	}
	return cases
}

func buildIO(cases []testCase) (string, string) {
	var inBuilder strings.Builder
	var outBuilder strings.Builder
	fmt.Fprintf(&inBuilder, "%d\n", len(cases))
	for _, tc := range cases {
		for i := 0; i < 8; i++ {
			fmt.Fprintln(&inBuilder, tc.board[i])
		}
		fmt.Fprintln(&outBuilder, tc.ans)
	}
	return inBuilder.String(), outBuilder.String()
}

func run(binary, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(binary, ".go") {
		cmd = exec.Command("go", "run", binary)
	} else {
		cmd = exec.Command(binary)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func normalizeTokens(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.TrimSpace(s)
	return strings.Fields(s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	cases := generateCases()
	input, expectedOutput := buildIO(cases)
	actualOutput, err := run(binary, input)
	if err != nil {
		fmt.Printf("Runtime error: %v\n", err)
		os.Exit(1)
	}
	if strings.Join(normalizeTokens(actualOutput), " ") != strings.Join(normalizeTokens(expectedOutput), " ") {
		fmt.Println("Wrong answer")
		fmt.Println("Expected:")
		fmt.Println(expectedOutput)
		fmt.Println("Got:")
		fmt.Println(actualOutput)
		os.Exit(1)
	}
	fmt.Println("All test cases passed!")
}
