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
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
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
