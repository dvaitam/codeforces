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

func randomID(rng *rand.Rand) string {
	l := rng.Intn(3) + 1
	b := make([]byte, l)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func randomSpaces(rng *rand.Rand) string {
	return strings.Repeat(" ", rng.Intn(3))
}

func generateCase(rng *rand.Rand) (string, string) {
	sections := rng.Intn(3) + 1
	var lines []string
	// some global key-values
	gCount := rng.Intn(3)
	for i := 0; i < gCount; i++ {
		key := randomID(rng)
		val := randomID(rng)
		line := randomSpaces(rng) + key + randomSpaces(rng) + "=" + randomSpaces(rng) + val + randomSpaces(rng)
		lines = append(lines, line)
	}
	for s := 0; s < sections; s++ {
		name := fmt.Sprintf("sec%d", s+1)
		secLine := randomSpaces(rng) + "[" + name + "]" + randomSpaces(rng)
		lines = append(lines, secLine)
		kvCount := rng.Intn(5) + 1
		for i := 0; i < kvCount; i++ {
			key := randomID(rng)
			val := randomID(rng)
			line := randomSpaces(rng) + key + randomSpaces(rng) + "=" + randomSpaces(rng) + val + randomSpaces(rng)
			lines = append(lines, line)
		}
	}
	// add some comments
	for i := 0; i < rng.Intn(3); i++ {
		lines = append(lines, randomSpaces(rng)+";comment")
	}
	n := len(lines)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, l := range lines {
		sb.WriteString(l)
		sb.WriteByte('\n')
	}
	input := sb.String()
	expected := computeExpected(lines)
	return input, expected
}

func computeExpected(lines []string) string {
	global := make(map[string]string)
	sections := make(map[string]map[string]string)
	var sectionNames []string
	current := ""
	for _, line := range lines {
		trimmedLeft := strings.TrimLeft(line, " \t")
		if len(trimmedLeft) > 0 && trimmedLeft[0] == ';' {
			continue
		}
		trimmed := strings.TrimSpace(line)
		if len(trimmed) >= 2 && trimmed[0] == '[' && trimmed[len(trimmed)-1] == ']' {
			name := strings.TrimSpace(trimmed[1 : len(trimmed)-1])
			current = name
			if _, ok := sections[current]; !ok {
				sections[current] = make(map[string]string)
				sectionNames = append(sectionNames, current)
			}
		} else {
			if idx := strings.Index(line, "="); idx != -1 {
				key := strings.TrimSpace(line[:idx])
				value := strings.TrimSpace(line[idx+1:])
				if current == "" {
					global[key] = value
				} else {
					sections[current][key] = value
				}
			}
		}
	}
	var gKeys []string
	for k := range global {
		gKeys = append(gKeys, k)
	}
	sort.Strings(gKeys)
	var out strings.Builder
	for _, k := range gKeys {
		out.WriteString(fmt.Sprintf("%s=%s\n", k, global[k]))
	}
	sort.Strings(sectionNames)
	for _, sec := range sectionNames {
		out.WriteString(fmt.Sprintf("[%s]\n", sec))
		var keys []string
		for k := range sections[sec] {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			out.WriteString(fmt.Sprintf("%s=%s\n", k, sections[sec][k]))
		}
	}
	res := strings.TrimRight(out.String(), "\n")
	return res
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected '%s' got '%s'", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
