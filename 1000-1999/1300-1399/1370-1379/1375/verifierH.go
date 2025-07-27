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

func buildRef() (string, error) {
	dir := "."
	src := "1375H.go"
	bin := filepath.Join(os.TempDir(), "refH.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build ref: %v, %s", err, string(out))
	}
	return bin, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	refBin, err := buildRef()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	file, err := os.Open("testcasesH.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		q, _ := strconv.Atoi(parts[1])
		need := 2 + n + 2*q
		if len(parts) != need {
			fmt.Printf("test %d wrong count\n", idx)
			os.Exit(1)
		}
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, q))
		for i := 0; i < n; i++ {
			input.WriteString(fmt.Sprintf("%s ", parts[2+i]))
		}
		input.WriteByte('\n')
		for i := 0; i < q; i++ {
			l := parts[2+n+2*i]
			r := parts[2+n+2*i+1]
			input.WriteString(fmt.Sprintf("%s %s\n", l, r))
		}

		// run reference
		cmdRef := exec.Command(refBin)
		cmdRef.Stdin = strings.NewReader(input.String())
		var refOut bytes.Buffer
		cmdRef.Stdout = &refOut
		if err := cmdRef.Run(); err != nil {
			fmt.Printf("ref run error on test %d: %v\n", idx, err)
			os.Exit(1)
		}
		expected := strings.TrimSpace(refOut.String())

		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input.String())
		var outBuf bytes.Buffer
		cmd.Stdout = &outBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		out := strings.TrimSpace(outBuf.String())
		if out != expected {
			fmt.Printf("Test %d failed\nexpected:\n%s\ngot:\n%s\n", idx, expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
