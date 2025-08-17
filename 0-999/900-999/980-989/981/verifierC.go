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

	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
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
