package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleF")
	cmd := exec.Command("go", "build", "-o", oracle, "858F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, string(out))
	}
	return oracle, nil
}

func readCases(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	cases := []string{}
	var sb strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			if sb.Len() > 0 {
				cases = append(cases, sb.String())
				sb.Reset()
			}
			continue
		}
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	if sb.Len() > 0 {
		cases = append(cases, sb.String())
	}
	return cases, scanner.Err()
}

func normalizeOutput(out string) (int, []string) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	norm := make([]string, 0, len(lines))
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		norm = append(norm, strings.Join(strings.Fields(ln), " "))
	}
	if len(norm) == 0 {
		return 0, nil
	}
	var cnt int
	fmt.Sscan(norm[0], &cnt)
	detail := append([]string(nil), norm[1:]...)
	sort.Strings(detail)
	return cnt, detail
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	const testcasesRaw = `4 4 1 2 2 4 4 1 3 2
10 1 7 5
5 8 5 3 3 4 1 5 1 4 1 2 2 3 5 2 4 2
7 1 2 3
10 8 2 10 8 3 10 5 2 6 1 3 3 7 1 2 3 6
6 13 1 3 6 1 3 6 5 3 2 4 6 5 5 4 2 6 6 4 2 5 2 1 4 1 1 5
2 1 1 2
5 5 4 2 3 4 1 5 1 4 2 1
10 5 6 1 5 4 3 2 9 6 8 10
2 1 2 1
2 1 2 1
9 0
9 15 1 9 2 3 6 8 8 9 2 7 3 9 8 3 8 1 8 7 2 1 8 2 1 7 2 6 4 2 7 9
3 0
3 3 1 3 3 2 1 2
2 1 1 2
4 4 1 4 1 2 3 4 2 3
6 10 2 3 5 4 3 6 3 5 1 6 6 4 2 1 4 3 2 5 6 5
4 3 2 4 1 2 3 1
5 2 5 2 5 4
7 12 4 2 7 5 7 1 6 1 4 3 4 7 5 1 5 2 5 6 1 4 2 1 3 2
5 1 3 1
1 0
5 8 3 2 2 4 3 4 1 2 1 3 5 1 5 3 4 1
2 1 1 2
1 0
3 3 2 3 1 2 3 1
7 15 6 5 1 7 1 2 4 6 4 7 7 3 2 6 4 3 1 4 7 6 1 3 7 2 4 2 3 5 5 1
1 0
1 0
5 4 3 5 1 2 3 1 2 5
7 2 5 6 1 5
2 1 1 2
1 0
7 3 3 7 7 2 6 7
8 5 1 2 4 3 5 8 7 8 2 5
5 10 4 1 4 3 5 1 1 2 3 2 4 5 4 2 5 2 3 1 5 3
2 0
8 11 4 1 2 7 4 6 7 8 6 5 5 4 4 3 3 6 8 6 8 5 5 1
8 10 6 4 6 1 3 4 2 6 7 3 4 2 7 6 5 6 1 3 2 1
8 0
4 4 2 4 4 3 3 2 1 3
5 10 2 1 2 3 3 4 1 5 4 1 4 5 2 5 5 3 1 3 4 2
3 0
4 6 4 2 3 4 2 3 2 1 3 1 4 1
7 0
6 2 4 6 3 2
6 14 1 3 4 3 5 3 5 2 5 1 2 3 2 4 6 2 1 6 5 6 6 3 1 4 2 1 6 4
10 7 3 7 7 5 10 7 7 8 2 3 5 2 2 8
8 9 8 6 4 5 3 8 1 7 2 3 4 1 1 6 3 7 8 2
3 3 2 1 3 2 1 3
5 9 3 5 1 5 2 3 5 4 3 1 4 1 2 4 1 2 2 5
7 4 6 1 5 4 5 2 7 3
10 0
2 1 2 1
1 0
9 3 9 1 3 6 6 4
1 0
5 2 1 2 3 1
1 0
10 2 7 10 8 6
3 2 1 3 3 2
3 1 2 3
10 1 3 7
5 4 5 2 1 2 4 1 5 4
2 1 1 2
3 1 3 2
8 8 1 6 7 6 4 6 8 4 8 6 4 3 7 4 5 1
1 0
9 14 2 6 1 5 4 2 5 9 1 2 1 3 5 7 2 7 4 8 3 5 4 1 9 6 8 1 8 7
7 8 3 7 4 2 5 1 6 4 6 7 5 7 6 5 1 6
6 15 3 2 1 4 5 2 6 5 6 3 4 2 1 6 5 1 1 2 5 4 3 1 5 3 2 6 6 4 3 4
4 4 1 3 1 2 3 2 3 4
9 5 6 9 5 3 7 8 2 3 6 8
5 3 1 4 3 2 5 1
10 7 6 4 3 5 2 6 3 10 10 2 9 2 10 4
7 7 5 3 6 1 1 5 7 5 6 3 5 4 2 4
1 0
2 0
8 11 2 5 3 4 4 1 8 2 8 4 2 4 3 8 8 6 3 5 5 4 2 6
6 13 2 1 3 2 2 5 4 1 4 5 4 6 3 1 4 3 2 4 3 5 6 5 6 2 1 6
10 10 3 2 7 2 8 3 7 3 1 9 4 9 10 7 6 5 6 7 4 5
3 3 1 3 3 2 1 2
6 15 2 5 6 5 6 4 1 3 1 4 2 4 3 2 4 3 2 6 1 6 3 5 5 1 4 5 3 6 1 2
7 2 5 3 5 4
10 10 9 3 3 1 8 3 4 1 6 4 8 9 10 9 5 10 4 2 4 5
4 1 1 3
2 1 1 2
7 0
9 12 2 3 3 6 9 7 4 7 7 6 6 4 1 2 7 8 6 2 8 1 7 1 6 9
1 0
2 1 2 1
1 0
7 12 4 7 6 1 7 2 2 5 6 2 7 5 3 1 4 2 1 2 6 4 3 2 7 1
10 3 10 9 5 2 8 7
6 6 5 4 3 1 4 1 4 6 1 2 3 6
4 0
5 1 2 4
10 7 5 9 6 8 7 10 6 1 1 2 10 5 7 9
3 3 2 3 2 1 1 3`

	cases, err := readCases(strings.NewReader(testcasesRaw))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read cases: %v\n", err)
		os.Exit(1)
	}

	for i, c := range cases {
		idx := i + 1
		cmdO := exec.Command(oracle)
		cmdO.Stdin = strings.NewReader(c)
		var outO bytes.Buffer
		cmdO.Stdout = &outO
		if err := cmdO.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "oracle run error on case %d: %v\n", idx, err)
			os.Exit(1)
		}
		expectedRaw := outO.String()
		expCnt, expList := normalizeOutput(expectedRaw)

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(c)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		gotRaw := out.String()
		gotCnt, gotList := normalizeOutput(gotRaw)

		if gotCnt != expCnt {
			fmt.Printf("case %d failed\nexpected: %s\n got: %s\n", idx, strings.TrimSpace(expectedRaw), strings.TrimSpace(gotRaw))
			os.Exit(1)
		}
		if len(gotList) != len(expList) {
			fmt.Printf("case %d failed (count mismatch)\nexpected: %s\n got: %s\n", idx, strings.TrimSpace(expectedRaw), strings.TrimSpace(gotRaw))
			os.Exit(1)
		}
		for j := range expList {
			if expList[j] != gotList[j] {
				fmt.Printf("case %d failed\nexpected: %s\n got: %s\n", idx, strings.TrimSpace(expectedRaw), strings.TrimSpace(gotRaw))
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
