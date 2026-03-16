package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesBRaw = `100
4 3
YES YES
2 1
YES YES
9 4
YES NO NO YES YES NO
5 5
NO
7 3
YES NO YES YES NO
9 4
YES YES NO YES YES NO
5 5
NO
7 4
YES NO YES YES
4 2
NO YES NO
6 6
NO
6 4
YES NO NO
9 4
NO NO YES NO YES NO
3 2
YES NO
4 3
YES YES
7 7
NO
3 1
YES NO NO
6 1
YES NO NO YES NO NO
2 1
YES NO
6 5
NO YES
2 2
YES
3 1
YES YES YES
7 2
YES YES YES YES YES NO
7 4
YES NO YES YES
7 5
NO YES YES
4 1
YES NO YES NO
3 2
YES YES
8 8
YES
4 3
NO NO
5 3
YES NO NO
6 4
NO NO YES
9 9
NO
6 3
YES NO NO YES
6 5
NO NO
3 2
NO YES
5 4
NO YES
7 6
YES YES
9 9
YES
8 2
NO NO NO NO NO NO NO
4 1
YES NO YES NO
6 3
NO NO NO NO
6 6
YES
2 2
YES
3 2
YES NO
3 2
NO YES
2 2
NO
3 2
NO YES
7 6
NO YES
3 3
YES
9 5
YES NO YES NO YES
4 3
YES YES
4 2
YES NO NO
7 4
YES YES NO NO
9 2
YES YES YES YES YES YES NO NO
3 1
NO YES YES
3 1
YES YES YES
2 2
NO
6 3
NO NO YES NO
6 3
NO YES YES NO
6 1
NO YES NO YES NO YES
6 5
YES NO
7 6
NO NO
9 2
NO NO NO YES NO YES YES YES
7 7
YES
2 2
NO
7 4
YES NO YES YES
4 3
YES YES
2 1
NO NO
9 1
YES NO YES YES YES NO YES NO YES
7 3
NO YES YES NO NO
9 9
YES
3 2
NO NO
9 8
YES NO
6 2
NO NO YES YES YES
6 3
NO YES NO YES
7 5
YES YES YES
2 1
YES NO
7 3
YES YES YES YES NO
8 6
NO NO YES
7 2
NO YES YES NO YES YES
9 1
YES NO NO NO YES YES NO YES NO
5 2
NO YES NO NO
3 3
YES
7 3
YES NO YES NO YES
9 7
NO YES YES
5 1
NO YES NO YES YES
3 1
NO YES NO
3 3
YES
8 3
NO NO YES YES YES NO
8 6
YES NO NO
8 7
YES YES
8 2
YES YES NO YES YES NO NO
5 1
NO YES NO YES YES
2 1
NO NO
3 2
NO YES
8 6
YES NO NO
7 7
YES
5 3
YES NO NO
5 1
YES NO NO YES YES
2 2
YES
2 2
NO
`

func generateNames() []string {
	names := make([]string, 0, 52)
	for i := 0; i < 52; i++ {
		first := 'A' + rune(i/26)
		second := 'a' + rune(i%26)
		names = append(names, fmt.Sprintf("%c%c", first, second))
	}
	return names
}

func solveCase(n, k int, tokens []string) []string {
	m := n - k + 1
	pool := generateNames()
	res := make([]string, n)
	idx := 0
	for i := 0; i < k-1; i++ {
		res[i] = pool[idx]
		idx++
	}
	for i := 0; i < m; i++ {
		if tokens[i] == "YES" {
			res[i+k-1] = pool[idx]
			idx++
		} else {
			res[i+k-1] = res[i]
		}
	}
	return res
}

func runCandidate(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("execution failed: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	scan := bufio.NewScanner(strings.NewReader(testcasesBRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	pool := generateNames()
	_ = pool
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("missing n")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		k, _ := strconv.Atoi(scan.Text())
		m := n - k + 1
		tokens := make([]string, m)
		for i := 0; i < m; i++ {
			scan.Scan()
			tokens[i] = scan.Text()
		}
		expected := solveCase(n, k, tokens)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < m; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(tokens[i])
		}
		sb.WriteByte('\n')
		out, err := runCandidate(os.Args[1], []byte(sb.String()))
		if err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
		gotTokens := strings.Fields(out)
		if len(gotTokens) != n {
			fmt.Printf("case %d failed: expected %d tokens got %d\n", caseIdx+1, n, len(gotTokens))
			os.Exit(1)
		}
		for i := 0; i < n; i++ {
			if gotTokens[i] != expected[i] {
				fmt.Printf("case %d failed at position %d: expected %s got %s\n", caseIdx+1, i+1, expected[i], gotTokens[i])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed!")
}
