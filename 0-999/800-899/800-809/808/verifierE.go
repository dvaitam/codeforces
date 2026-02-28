package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const randomSeed int64 = 808
const randomCaseCount = 300

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "808E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, string(out))
	}
	return oracle, nil
}

func buildCase(n, m int, items [][2]int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, it := range items {
		sb.WriteString(fmt.Sprintf("%d %d\n", it[0], it[1]))
	}
	return sb.String()
}

func fixedCases() []string {
	return []string{
		buildCase(1, 1, [][2]int{{1, 1}}),
		buildCase(3, 3, [][2]int{{1, 10}, {2, 20}, {3, 100}}),
		buildCase(6, 7, [][2]int{{1, 3}, {1, 4}, {2, 5}, {2, 6}, {3, 10}, {3, 11}}),
		buildCase(8, 10, [][2]int{{1, 1000}, {1, 1}, {1, 1}, {2, 100}, {2, 99}, {3, 1000}, {3, 998}, {3, 997}}),
		buildCase(10, 1, [][2]int{{2, 100}, {3, 200}, {2, 50}, {3, 1}, {2, 2}, {3, 3}, {2, 4}, {3, 5}, {2, 6}, {3, 7}}),
		// Keep a high-value mix while staying within signed 32-bit answer range.
		// This avoids platform-specific printf("%I64d") behaviour in legacy C++ submissions.
		buildCase(12, 30, [][2]int{{1, 1000000}, {1, 999999}, {2, 1000000}, {2, 999998}, {3, 1000000}, {3, 999997}, {1, 7}, {2, 8}, {3, 9}, {1, 10}, {2, 11}, {3, 12}}),
	}
}

func firstIntToken(s string) (int64, error) {
	for _, tok := range strings.Fields(s) {
		if v, err := strconv.ParseInt(tok, 10, 64); err == nil {
			return v, nil
		}
	}
	return 0, fmt.Errorf("no integer token in output")
}

func randomCases() []string {
	rng := rand.New(rand.NewSource(randomSeed))
	cases := make([]string, 0, randomCaseCount)

	for i := 0; i < randomCaseCount; i++ {
		n := rng.Intn(70) + 1
		m := rng.Intn(150) + 1
		items := make([][2]int, n)
		for j := 0; j < n; j++ {
			w := rng.Intn(3) + 1
			v := rng.Intn(2000) + 1
			items[j] = [2]int{w, v}
		}
		cases = append(cases, buildCase(n, m, items))
	}

	return cases
}

func runCase(bin, input string) (int64, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return 0, stderr.String(), err
	}

	rawOut := strings.TrimSpace(out.String())
	if rawOut == "" {
		// Some binaries write the answer to stderr (or mix streams); accept first integer from either stream.
		rawOut = strings.TrimSpace(stderr.String())
	}
	v, parseErr := firstIntToken(rawOut)
	if parseErr != nil {
		merged := strings.TrimSpace(out.String())
		if errTxt := strings.TrimSpace(stderr.String()); errTxt != "" {
			if merged != "" {
				merged += "\n"
			}
			merged += errTxt
		}
		if merged != "" {
			if v2, err := firstIntToken(merged); err == nil {
				return v2, strings.TrimSpace(stderr.String()), nil
			}
		}
		return 0, strings.TrimSpace(stderr.String()), fmt.Errorf("%w: stdout=%q stderr=%q", parseErr, out.String(), stderr.String())
	}

	return v, strings.TrimSpace(stderr.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	cases := append(fixedCases(), randomCases()...)
	for i, c := range cases {
		idx := i + 1

		expected, oracleErrOut, err := runCase(oracle, c)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle run error on case %d: %v\nstderr: %s\n", idx, err, oracleErrOut)
			os.Exit(1)
		}

		got, binErrOut, err := runCase(bin, c)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", idx, err, binErrOut)
			os.Exit(1)
		}

		if got != expected {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %d\n     got: %d\n", idx, c, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed (fixed + random, seed=%d)\n", len(cases), randomSeed)
}
