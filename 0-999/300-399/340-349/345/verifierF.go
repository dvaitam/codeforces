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

func run(bin string, input string) (string, error) {
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

func normalize(s string) string {
	fields := strings.Fields(s)
	for i, w := range fields {
		fields[i] = strings.ToLower(w)
	}
	return strings.Join(fields, " ")
}

func solveCase(countries [][]string) []string {
	counts := make(map[string]int)
	for _, list := range countries {
		set := make(map[string]bool)
		for _, sup := range list {
			set[normalize(sup)] = true
		}
		for s := range set {
			counts[s]++
		}
	}
	maxCount := 0
	for _, c := range counts {
		if c > maxCount {
			maxCount = c
		}
	}
	var best []string
	for s, c := range counts {
		if c == maxCount {
			best = append(best, s)
		}
	}
	sort.Strings(best)
	return best
}

func runCase(bin string, countries [][]string) error {
	var sb strings.Builder
	for idx, list := range countries {
		if idx > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprintf("Country%d\n", idx))
		for _, s := range list {
			sb.WriteString("* ")
			sb.WriteString(s)
			sb.WriteByte('\n')
		}
	}
	out, err := run(bin, sb.String())
	if err != nil {
		return err
	}
	gotLines := strings.Split(strings.TrimSpace(out), "\n")
	expected := solveCase(countries)
	if len(gotLines) != len(expected) {
		return fmt.Errorf("expected %d lines got %d", len(expected), len(gotLines))
	}
	for i := range expected {
		if strings.TrimSpace(gotLines[i]) != expected[i] {
			return fmt.Errorf("line %d expected %q got %q", i+1, expected[i], strings.TrimSpace(gotLines[i]))
		}
	}
	return nil
}

func randomSup(rng *rand.Rand) string {
	words := rng.Intn(3) + 1
	var parts []string
	for i := 0; i < words; i++ {
		l := rng.Intn(5) + 1
		b := make([]byte, l)
		for j := 0; j < l; j++ {
			b[j] = byte('a' + rng.Intn(26))
		}
		parts = append(parts, string(b))
	}
	return strings.Join(parts, " ")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		c := rng.Intn(3) + 1
		countries := make([][]string, c)
		for j := 0; j < c; j++ {
			n := rng.Intn(3) + 1
			countries[j] = make([]string, n)
			for k := 0; k < n; k++ {
				countries[j][k] = randomSup(rng)
			}
		}
		if err := runCase(bin, countries); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
