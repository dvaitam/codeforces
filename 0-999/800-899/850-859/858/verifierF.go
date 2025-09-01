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

	f, err := os.Open("testcasesF.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	cases, err := readCases(f)
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
