package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD4.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, input := range tests {
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(refOut) != strings.TrimSpace(candOut) {
			fmt.Fprintf(os.Stderr, "Mismatch on test %d\nInput:\n%sExpected: %s\nGot: %s\n", idx+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_207D4.bin"
	cmd := exec.Command("go", "build", "-o", refName, "207D4.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
}

func runProgram(target, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func buildTests() []string {
	return []string{
		"1\nGlobal Markets\nTrade discussions dominate the global markets today with new tariffs looming.\n",
		"2\nLaboratory Notes\nWe analyze algorithms and data structures for academic research.\n",
		"3\nPoetry Collection\nA river flows silently while hearts dream.\n",
		"4\nEconomic Outlook\nCurrency exchange and trade balances show improvement.\n",
		"5\nUniversity Schedule\nThe semester timetable lists lectures and exams.\n",
		"6\nLove Letter\nDearest friend, the stars shine brightly tonight.\n",
		"7\nCorporate Report\nQuarterly profits and trade agreements strengthen the company.\n",
		"8\nCampus Life\nStudents discuss assignments and research topics.\n",
		"9\nNature Journal\nThe forest whispers through winds and waterfalls.\n",
		"10\nTech Review\nProcessors, compilers, and algorithms are compared in depth.\n",
	}
}
