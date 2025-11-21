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

const refSourceD1 = "1000-1999/1900-1999/1970-1979/1970/1970D1.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceD1)
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	refWordList, refPowerMap, err := collectData(refBin)
	if err != nil {
		fail("reference run failed: %v", err)
	}

	if err := validateSolution(refWordList, refPowerMap); err != nil {
		fail("reference validation failed: %v", err)
	}

	candWordList, candPowerMap, err := collectDataWithInput(candidate, refWordList, refPowerMap)
	if err != nil {
		fail("candidate run failed: %v", err)
	}

	if err := compareOutputs(refWordList, refPowerMap, candWordList, candPowerMap); err != nil {
		fail("candidate validation failed: %v", err)
	}

	fmt.Println("OK")
}

func buildReference(src string) (string, error) {
	tmp, err := os.CreateTemp("", "1970D1-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(src))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func collectData(bin string) ([]string, map[int][2]int, error) {
	input := "3 3\n1\n2\n3\n"
	output, err := runProgram(bin, input)
	if err != nil {
		return nil, nil, err
	}
	return parseOutput(output)
}

func collectDataWithInput(target string, words []string, powerMap map[int][2]int) ([]string, map[int][2]int, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", len(words), len(powerMap)))
	for power := range powerMap {
		sb.WriteString(fmt.Sprintf("%d\n", power))
	}
	output, err := runCandidate(target, sb.String())
	if err != nil {
		return nil, nil, err
	}
	return parseOutput(output)
}

func parseOutput(out string) ([]string, map[int][2]int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 3 {
		return nil, nil, fmt.Errorf("insufficient output lines")
	}
	n, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, nil, fmt.Errorf("invalid n: %v", err)
	}
	words := lines[1 : 1+n]
	powerLines := lines[1+n:]
	powerMap := make(map[int][2]int)
	for _, line := range powerLines {
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("invalid power map line: %s", line)
		}
		a, _ := strconv.Atoi(parts[0])
		b, _ := strconv.Atoi(parts[1])
		power := len(words[a-1]) + len(words[b-1])
		if _, exists := powerMap[power]; !exists {
			powerMap[power] = [2]int{a, b}
		}
	}
	return words, powerMap, nil
}

func validateSolution(words []string, powerMap map[int][2]int) error {
	if len(words) < 3 {
		return fmt.Errorf("insufficient number of words")
	}
	wordSet := make(map[string]bool)
	for _, word := range words {
		if wordSet[word] {
			return fmt.Errorf("duplicate word: %s", word)
		}
		wordSet[word] = true
	}
	if len(powerMap) == 0 {
		return fmt.Errorf("empty power map")
	}
	return nil
}

func compareOutputs(refWords []string, refPowerMap map[int][2]int, candWords []string, candPowerMap map[int][2]int) error {
	if len(refWords) != len(candWords) {
		return fmt.Errorf("word list length mismatch")
	}
	for i := range refWords {
		if refWords[i] != candWords[i] {
			return fmt.Errorf("word mismatch at index %d", i)
		}
	}
	if len(refPowerMap) != len(candPowerMap) {
		return fmt.Errorf("power map size mismatch")
	}
	for power, pair := range refPowerMap {
		if candPair, ok := candPowerMap[power]; !ok || candPair != pair {
			return fmt.Errorf("mismatch for power %d", power)
		}
	}
	return nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func runCandidate(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
