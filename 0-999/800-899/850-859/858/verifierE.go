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

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "858E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, string(out))
	}
	return oracle, nil
}

func generateCases() []string {
	cases := []string{
		"1\n1 1\n",
		"1\na 0\n",
		"2\n1 1\n2 0\n",
		"2\na 1\nb 0\n",
		"3\n1 1\n2 1\n3 0\n",
		"4\nuser1 1\n2 1\n3 0\n4 0\n",
	}

	rng := rand.New(rand.NewSource(858))
	for t := 0; t < 200; t++ {
		n := rng.Intn(40) + 1
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		used := map[string]bool{}
		for i := 0; i < n; i++ {
			op := rng.Intn(2)
			name := ""
			for {
				nameType := rng.Intn(4)
				switch nameType {
				case 0:
					name = strconv.Itoa(rng.Intn(n) + 1)
				case 1:
					name = strconv.Itoa(n + rng.Intn(30) + 1)
				case 2:
					name = fmt.Sprintf("u%d", rng.Intn(1000))
				default:
					name = fmt.Sprintf("0%d", rng.Intn(100)+1)
				}
				if !used[name] {
					used[name] = true
					break
				}
			}
			sb.WriteString(name)
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(op))
			sb.WriteByte('\n')
		}
		cases = append(cases, sb.String())
	}
	return cases
}

func firstInt(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil || k < 0 {
		return 0, fmt.Errorf("invalid first token %q", fields[0])
	}
	return k, nil
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

	cases := generateCases()

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
		expectedK, err := firstInt(outO.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle parse error on case %d: %v\n", idx, err)
			os.Exit(1)
		}

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
		gotK, err := firstInt(out.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\noutput:\n%s\ninput:\n%s\n", idx, err, out.String(), c)
			os.Exit(1)
		}
		if gotK != expectedK {
			fmt.Printf("case %d failed\nexpected moves: %d\n got moves: %d\ninput:\n%s\n", idx, expectedK, gotK, c)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
