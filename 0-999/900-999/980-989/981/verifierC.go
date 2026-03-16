package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", oracle, "981C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func normalize(out string) (string, error) {
	out = strings.TrimSpace(out)
	lines := strings.Split(out, "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	if len(lines) == 0 {
		return "", fmt.Errorf("empty output")
	}
	if lines[0] == "No" {
		return "No", nil
	}
	if lines[0] != "Yes" {
		return "", fmt.Errorf("first line must be Yes or No")
	}
	if len(lines) < 2 {
		return "", fmt.Errorf("missing count line")
	}
	k, err := strconv.Atoi(lines[1])
	if err != nil {
		return "", err
	}
	pairs := lines[2:]
	if len(pairs) != k {
		return "", fmt.Errorf("expected %d pairs, got %d", k, len(pairs))
	}
	if k == 1 {
		f := strings.Fields(pairs[0])
		if len(f) != 2 {
			return "", fmt.Errorf("invalid pair: %s", pairs[0])
		}
		a, err1 := strconv.Atoi(f[0])
		b, err2 := strconv.Atoi(f[1])
		if err1 != nil || err2 != nil {
			return "", fmt.Errorf("invalid pair: %s", pairs[0])
		}
		if a > b {
			a, b = b, a
		}
		pairs[0] = fmt.Sprintf("%d %d", a, b)
	}
	sort.Strings(pairs)
	return "Yes\n" + fmt.Sprintln(k) + strings.Join(pairs, "\n"), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	const testcasesRaw = `8
1 2
1 3
3 4
2 5
1 6
1 7
3 8
2
1 2
6
1 2
2 3
3 4
4 5
1 6
7
1 2
1 3
2 4
1 5
3 6
3 7
4
1 2
2 3
1 4
4
1 2
1 3
1 4
6
1 2
1 3
3 4
2 5
4 6
4
1 2
2 3
2 4
6
1 2
2 3
3 4
4 5
2 6
4
1 2
2 3
3 4
7
1 2
1 3
3 4
3 5
3 6
1 7
4
1 2
2 3
1 4
5
1 2
1 3
3 4
3 5
6
1 2
2 3
3 4
2 5
2 6
8
1 2
2 3
1 4
4 5
3 6
2 7
4 8
5
1 2
2 3
1 4
4 5
4
1 2
1 3
1 4
4
1 2
1 3
3 4
2
1 2
4
1 2
1 3
2 4
5
1 2
1 3
1 4
3 5
8
1 2
2 3
1 4
3 5
4 6
4 7
4 8
4
1 2
1 3
3 4
8
1 2
1 3
1 4
1 5
5 6
6 7
4 8
7
1 2
1 3
2 4
1 5
4 6
1 7
7
1 2
1 3
2 4
3 5
5 6
1 7
6
1 2
1 3
1 4
1 5
1 6
4
1 2
2 3
1 4
3
1 2
2 3
4
1 2
2 3
2 4
2
1 2
2
1 2
4
1 2
2 3
3 4
2
1 2
2
1 2
6
1 2
1 3
1 4
1 5
2 6
5
1 2
1 3
2 4
3 5
8
1 2
2 3
3 4
1 5
5 6
4 7
1 8
2
1 2
4
1 2
2 3
1 4
4
1 2
1 3
2 4
7
1 2
1 3
2 4
1 5
5 6
4 7
6
1 2
2 3
2 4
3 5
4 6
7
1 2
2 3
2 4
3 5
2 6
2 7
3
1 2
2 3
7
1 2
2 3
2 4
3 5
2 6
1 7
2
1 2
7
1 2
1 3
2 4
3 5
5 6
2 7
2
1 2
5
1 2
2 3
2 4
1 5
3
1 2
1 3
5
1 2
1 3
3 4
4 5
6
1 2
1 3
1 4
4 5
1 6
6
1 2
2 3
3 4
1 5
2 6
6
1 2
2 3
2 4
3 5
1 6
2
1 2
4
1 2
2 3
2 4
8
1 2
2 3
2 4
4 5
5 6
1 7
6 8
5
1 2
2 3
1 4
4 5
7
1 2
2 3
3 4
3 5
4 6
5 7
4
1 2
2 3
3 4
7
1 2
2 3
1 4
4 5
4 6
4 7
4
1 2
1 3
2 4
3
1 2
1 3
6
1 2
1 3
1 4
4 5
1 6
8
1 2
1 3
3 4
1 5
2 6
6 7
2 8
2
1 2
7
1 2
1 3
2 4
4 5
4 6
3 7
6
1 2
1 3
2 4
2 5
1 6
7
1 2
2 3
3 4
3 5
3 6
3 7
5
1 2
2 3
2 4
4 5
2
1 2
8
1 2
1 3
2 4
3 5
5 6
2 7
4 8
2
1 2
2
1 2
8
1 2
2 3
1 4
3 5
3 6
1 7
1 8
2
1 2
3
1 2
2 3
7
1 2
1 3
1 4
2 5
1 6
1 7
7
1 2
1 3
1 4
1 5
4 6
4 7
2
1 2
5
1 2
1 3
3 4
1 5
3
1 2
2 3
6
1 2
1 3
3 4
1 5
5 6
3
1 2
2 3
4
1 2
2 3
2 4
3
1 2
1 3
6
1 2
2 3
3 4
3 5
3 6
6
1 2
1 3
1 4
1 5
2 6
4
1 2
1 3
1 4
2
1 2
7
1 2
1 3
2 4
3 5
2 6
1 7
3
1 2
2 3
3
1 2
2 3
3
1 2
2 3
3
1 2
1 3
4
1 2
1 3
2 4
8
1 2
2 3
1 4
3 5
2 6
4 7
1 8
8
1 2
2 3
1 4
3 5
1 6
6 7
3 8
7
1 2
2 3
3 4
4 5
1 6
6 7`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for {
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		n, _ := strconv.Atoi(line)
		var sb strings.Builder
		sb.WriteString(line + "\n")
		for i := 0; i < n-1; i++ {
			scanner.Scan()
			l := strings.TrimSpace(scanner.Text())
			sb.WriteString(l + "\n")
		}
		idx++
		input := sb.String()
		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		ne, err := normalize(exp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "normalize oracle output on case %d: %v\n", idx, err)
			os.Exit(1)
		}
		ng, err := normalize(got)
		if err != nil {
			fmt.Fprintf(os.Stderr, "normalize candidate output on case %d: %v\n", idx, err)
			os.Exit(1)
		}
		if ne != ng {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
