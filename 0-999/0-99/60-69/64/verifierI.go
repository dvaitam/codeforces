package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type rule struct {
	idx int
	asc bool
}

type row struct {
	cells []string
	idx   int
}

func sortTable(header []string, rules []rule, rows []row) []row {
	sort.SliceStable(rows, func(i, j int) bool {
		a, b := rows[i], rows[j]
		for _, r := range rules {
			if a.cells[r.idx] == b.cells[r.idx] {
				continue
			}
			if r.asc {
				return a.cells[r.idx] < b.cells[r.idx]
			}
			return a.cells[r.idx] > b.cells[r.idx]
		}
		return a.idx < b.idx
	})
	return rows
}

func randWord(rng *rand.Rand) string {
	l := rng.Intn(5) + 1
	b := make([]byte, l)
	for i := range b {
		if rng.Intn(2) == 0 {
			b[i] = byte('A' + rng.Intn(26))
		} else {
			b[i] = byte('a' + rng.Intn(26))
		}
	}
	return string(b)
}

func generateCase(rng *rand.Rand) (string, []string) {
	cols := rng.Intn(3) + 2
	headers := make([]string, cols)
	for i := range headers {
		headers[i] = randWord(rng)
	}
	var sb strings.Builder
	sb.WriteString(strings.Join(headers, " ") + "\n")
	// rules
	ruleCnt := rng.Intn(cols) + 1
	perm := rand.Perm(cols)
	var rules []rule
	var ruleParts []string
	for i := 0; i < ruleCnt; i++ {
		idx := perm[i]
		asc := rng.Intn(2) == 0
		order := "ASC"
		if !asc {
			order = "DESC"
		}
		ruleParts = append(ruleParts, fmt.Sprintf("%s %s", headers[idx], order))
		rules = append(rules, rule{idx: idx, asc: asc})
	}
	sb.WriteString(strings.Join(ruleParts, ", ") + "\n")
	rowsN := rng.Intn(10) + 1
	rows := make([]row, rowsN)
	for i := 0; i < rowsN; i++ {
		cells := make([]string, cols)
		for j := 0; j < cols; j++ {
			cells[j] = randWord(rng)
		}
		rows[i] = row{cells: cells, idx: i}
		sb.WriteString(strings.Join(cells, " ") + "\n")
	}
	sorted := sortTable(headers, rules, append([]row(nil), rows...))
	var expected []string
	expected = append(expected, strings.Join(headers, " "))
	for _, r := range sorted {
		expected = append(expected, strings.Join(r.cells, " "))
	}
	return sb.String(), expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expectedLines := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		outLines := strings.Split(strings.TrimSpace(out), "\n")
		if len(outLines) != len(expectedLines) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got %d\ninput:\n%s", i+1, len(expectedLines), len(outLines), input)
			os.Exit(1)
		}
		for j := range expectedLines {
			if strings.TrimSpace(outLines[j]) != expectedLines[j] {
				fmt.Fprintf(os.Stderr, "case %d failed: line %d expected '%s' got '%s'\ninput:\n%s", i+1, j+1, expectedLines[j], outLines[j], input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
