package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", oracle, "677D.go")
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	const testcasesDRaw = `2 2 3 2 3 1 3
4 3 9 2 6 9 9 9 8 7 3 4 3 5 1
2 1 2 2 1
4 4 13 3 6 2 1 3 8 4 5 11 7 13 12 10 7 13 9
4 2 6 6 1 5 4 6 2 6 3
3 3 2 2 2 2 1 2 1 2 1 1
3 4 7 1 7 5 5 7 1 4 6 2 3 5 3
2 1 2 2 1
1 1 1 1
3 3 3 3 2 2 2 1 2 2 2 3
4 1 3 3 1 2 3
4 3 9 9 6 1 7 9 5 1 7 3 1 4 8
3 4 1 1 1 1 1 1 1 1 1 1 1 1 1
3 3 5 1 1 5 4 3 5 5 2 3
3 2 6 6 5 3 6 1 4
4 2 2 1 2 2 1 1 2 1 1
3 2 3 3 1 3 3 2 3
2 4 5 1 3 4 3 5 4 2 4
2 2 1 1 1 1 1
4 3 9 3 2 5 2 4 9 1 6 9 4 7 8
1 2 2 2 1
4 1 1 1 1 1 1
1 3 1 1 1 1
2 2 3 1 3 2 3
4 1 3 2 3 3 1
1 4 3 2 1 1 1
1 1 1 1
1 4 3 2 3 1 2
4 4 10 5 4 6 7 2 3 9 1 7 2 10 10 3 1 6 8
4 1 4 4 3 4 2
4 1 2 2 2 1 2
2 4 3 3 2 2 3 2 1 2 3
1 4 1 1 1 1 1
2 2 4 3 2 1 4
3 1 1 1 1 1
3 1 1 1 1 1
1 1 1 1
1 2 1 1 1
4 3 6 4 3 5 4 1 4 5 2 4 6 6 2
4 4 5 2 2 1 4 4 5 4 5 2 2 3 2 5 2 5 5
3 2 6 3 6 6 4 2 5
3 2 3 3 2 2 2 1 1
3 2 3 1 2 3 3 2 3
2 4 1 1 1 1 1 1 1 1 1
2 3 2 2 2 1 1 1 2
2 1 2 1 2
1 1 1 1
3 1 2 2 2 1
1 3 2 1 2 2
2 4 7 5 2 7 1 3 1 6 4
4 3 2 2 2 2 2 2 1 2 1 2 2 2 1
3 2 6 5 6 2 4 3 1
1 2 2 1 2
4 4 10 2 2 3 9 9 10 10 7 8 2 5 5 7 1 3 4
2 3 3 2 3 1 3 2 2
1 2 2 2 1
1 2 1 1 1
2 2 4 2 4 3 2
3 1 2 2 2 1
3 4 9 5 8 1 7 1 2 4 8 3 9 9 8
2 1 2 1 2
2 4 2 2 1 2 1 1 2 1 2
3 3 6 1 4 6 5 6 1 6 3 2
4 2 8 2 1 7 2 4 8 3 6
1 4 2 2 1 2 2
4 4 15 8 15 1 12 13 15 2 14 3 8 9 11 10 7 4 1
4 2 3 1 1 3 3 2 3 1 3
1 2 1 1 1
4 2 1 1 1 1 1 1 1 1 1
4 1 4 3 2 2 4
4 4 7 3 5 1 3 1 7 5 1 7 6 3 3 7 4 7 2
2 3 1 1 1 1 1 1 1
1 1 1 1
1 1 1 1
4 1 2 1 2 2 2
4 3 11 11 7 5 10 11 6 1 5 8 7 2 9
4 2 7 2 3 7 4 6 5 2 1
2 4 8 7 3 3 4 8 2 2 6
4 4 11 1 2 4 7 9 6 10 5 2 8 2 11 11 4 4 1
1 1 1 1
2 2 4 3 4 4 2
1 1 1 1
3 1 1 1 1 1
1 4 3 1 3 2 3
3 4 7 1 5 1 7 4 3 6 6 1 3 3 2
4 1 3 3 2 3 3
2 3 1 1 1 1 1 1 1
2 3 2 1 1 2 2 2 2
2 4 8 6 1 4 1 5 2 8 7
1 2 1 1 1
4 1 3 3 1 2 3
2 1 2 2 1
1 2 1 1 1
4 3 3 2 2 2 3 3 1 3 2 2 2 1 2
3 3 5 4 1 3 5 2 5 2 5 1
1 2 2 1 2
4 2 4 4 3 4 3 2 1 4 3
2 4 7 3 1 2 5 7 1 6 6
2 4 6 5 1 6 2 1 1 4 6
4 4 14 9 11 3 5 3 2 11 3 6 12 2 14 7 13 10 8`

	scanner := bufio.NewScanner(strings.NewReader(testcasesDRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 3 {
			fmt.Printf("case %d malformed\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		p, _ := strconv.Atoi(parts[2])
		if len(parts) != 3+n*m {
			fmt.Printf("case %d expected %d grid values got %d\n", idx, n*m, len(parts)-3)
			os.Exit(1)
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, m, p)
		k := 3
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(parts[k])
				k++
			}
			sb.WriteByte('\n')
		}
		input := sb.String()

		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed: expected %s got %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
