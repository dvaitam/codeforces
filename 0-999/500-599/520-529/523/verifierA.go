package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runCandidate(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expectedOutput(w, h int, lines []string) string {
	img := make([][]byte, h)
	for i := 0; i < h; i++ {
		img[i] = []byte(lines[i])
	}
	outRows := 2 * w
	outCols := 2 * h
	var b strings.Builder
	for r2 := 0; r2 < outRows; r2++ {
		r := r2 / 2
		for c2 := 0; c2 < outCols; c2++ {
			c := c2 / 2
			b.WriteByte(img[c][r])
		}
		if r2+1 < outRows {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read test count: %v\n", err)
		os.Exit(1)
	}
	for caseNum := 1; caseNum <= t; caseNum++ {
		var w, h int
		fmt.Fscan(reader, &w, &h)
		lines := make([]string, h)
		for i := 0; i < h; i++ {
			fmt.Fscan(reader, &lines[i])
		}
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", w, h)
		for i := 0; i < h; i++ {
			input.WriteString(lines[i])
			input.WriteByte('\n')
		}
		want := expectedOutput(w, h, lines)
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", caseNum, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
