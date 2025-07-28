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

func runProg(bin string, input string) (string, error) {
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

func expected(n, k int) string {
	if k%4 == 0 {
		return "NO"
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	if k%4 == 2 {
		flip := true
		for i := 1; i <= n; i += 2 {
			if flip {
				fmt.Fprintf(&sb, "%d %d\n", i+1, i)
			} else {
				fmt.Fprintf(&sb, "%d %d\n", i, i+1)
			}
			flip = !flip
		}
	} else {
		for i := 1; i <= n; i += 2 {
			fmt.Fprintf(&sb, "%d %d\n", i, i+1)
		}
	}
	return strings.TrimSpace(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "bad line %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		k, _ := strconv.Atoi(parts[1])
		input := fmt.Sprintf("1\n%d %d\n", n, k)
		out, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		expect := expected(n, k)
		if strings.Fields(out) == nil {
			// unreachable but for linter
			fmt.Printf("%s", "")
		}
		eTokens := strings.Fields(expect)
		gTokens := strings.Fields(out)
		if len(eTokens) != len(gTokens) {
			fmt.Printf("case %d failed: token length mismatch\nexpected: %v\ngot: %v\n", idx, eTokens, gTokens)
			os.Exit(1)
		}
		for i := range eTokens {
			if eTokens[i] != gTokens[i] {
				fmt.Printf("case %d failed: expected %s got %s\n", idx, expect, out)
				os.Exit(1)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
