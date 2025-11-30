package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `100
4 4 0111 1110
5 2 10010 10
6 1 110111 0
5 4 11000 1011
2 6 10 000010
1 3 1 011
6 1 101101 1
5 2 10000 11
1 1 0 0
1 1 1 1
5 2 01111 10
3 5 011 00010
6 2 101100 00
1 5 0 00001
1 3 0 000
2 6 01 000101
1 2 0 11
4 2 0100 10
3 3 100 000
2 4 01 1111
4 3 1000 001
4 1 0010 0
1 5 0 00011
3 5 000 00010
6 3 001111 101
1 1 0 0
4 5 0111 01110
6 5 001100 11010
6 4 100111 1111
5 2 11101 10
4 4 0101 0110
2 6 11 101010
2 4 00 0011
1 2 1 01
4 2 1111 00
1 5 0 10010
2 5 10 10101
4 3 0110 010
2 1 10 0
4 5 1010 00000
2 5 01 10110
5 3 00000 010
4 6 1010 000001
6 2 110110 10
4 5 0100 01000
4 4 1010 0010
1 3 0 010
6 1 010010 0
6 5 001010 10000
1 1 1 1
2 6 10 000000
4 2 0101 01
3 1 001 1
4 6 0100 000000
1 2 1 11
4 4 0000 1001
3 3 000 101
4 1 0000 1
6 1 000010 1
3 5 101 11000
6 5 110111 11000
3 1 001 0
1 6 0 111111
5 6 00101 000000
3 4 101 1110
6 2 010010 01
4 1 1110 1
5 5 10110 01110
5 1 10101 1
5 2 11011 01
5 2 00010 01
3 3 100 101
1 2 0 10
5 5 11000 00000
5 4 01110 1011
5 1 10001 0
3 4 101 0011
4 2 1000 10
5 5 00110 10100
2 6 10 111000
2 2 00 01
2 5 11 10001
5 4 01001 0101
4 2 0111 11
3 2 001 10
4 2 1011 10
6 5 000000 10010
6 1 101101 0
5 5 00010 00100
5 5 10101 01000
2 5 01 10101
6 3 100001 100
3 2 010 11
3 2 001 01
5 3 11010 101
4 2 1000 10
5 6 01011 101010
6 6 111111 111111
6 3 000000 101
4 5 0101 00000
5 2 10111 00
2 5 10 01101
6 4 110000 0101
4 1 0011 0
5 4 11111 0010
2 5 01 01100
6 1 111111 1
6 1 100111 0
2 2 01 10
5 4 10111 1111
3 4 111 0000
6 5 101010 11110
5 4 00000 1000
5 3 10110 010
3 5 101 00000`

func isGood(s string) bool {
	for i := 0; i+1 < len(s); i++ {
		if s[i] == s[i+1] {
			return false
		}
	}
	return true
}

func expected(n, m int, s, t string) string {
	if isGood(s) {
		return "YES"
	}
	if !isGood(t) {
		return "NO"
	}
	has00, has11 := false, false
	for i := 0; i+1 < len(s); i++ {
		if s[i] == '0' && s[i+1] == '0' {
			has00 = true
		}
		if s[i] == '1' && s[i+1] == '1' {
			has11 = true
		}
	}
	if has00 && has11 {
		return "NO"
	}
	first, last := t[0], t[len(t)-1]
	if has00 {
		if first == '1' && last == '1' {
			return "YES"
		}
		return "NO"
	}
	if has11 {
		if first == '0' && last == '0' {
			return "YES"
		}
		return "NO"
	}
	return "NO"
}

func loadCases() ([]string, []string) {
	tokens := strings.Fields(testcasesRaw)
	if len(tokens) == 0 {
		fmt.Println("no embedded testcases")
		os.Exit(1)
	}
	t, err := strconv.Atoi(tokens[0])
	if err != nil {
		fmt.Println("invalid testcase count")
		os.Exit(1)
	}
	pos := 1
	var inputs []string
	var expects []string
	for i := 0; i < t; i++ {
		if pos+3 >= len(tokens) {
			fmt.Printf("case %d incomplete\n", i+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(tokens[pos])
		m, _ := strconv.Atoi(tokens[pos+1])
		s := tokens[pos+2]
		tstr := tokens[pos+3]
		pos += 4
		inputs = append(inputs, fmt.Sprintf("1\n%d %d\n%s\n%s\n", n, m, s, tstr))
		expects = append(expects, expected(n, m, s, tstr))
	}
	return inputs, expects
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	inputs, expects := loadCases()
	for idx, input := range inputs {
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expects[idx] {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expects[idx], got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}
