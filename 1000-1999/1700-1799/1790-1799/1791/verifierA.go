package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func caseCount(data []byte) int {
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if scan.Scan() {
		n, _ := strconv.Atoi(scan.Text())
		return n
	}
	return 0
}

func run(bin string, data []byte) ([]byte, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(data)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.Bytes(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	dir, _ := os.Getwd()
	oracle := filepath.Join(dir, "oracleA")
	if out, err := exec.Command("go", "build", "-o", oracle, "1791A.go").CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build oracle failed: %v\n%s", err, out)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read testcases: %v\n", err)
		os.Exit(1)
	}

	candOut, err := run(candidate, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n%s", err, candOut)
		os.Exit(1)
	}
	refOut, err := run("./"+filepath.Base(oracle), data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle runtime error: %v\n%s", err, refOut)
		os.Exit(1)
	}
	cLines := bytes.Split(bytes.TrimSpace(candOut), []byte{'\n'})
	rLines := bytes.Split(bytes.TrimSpace(refOut), []byte{'\n'})
	if len(cLines) != len(rLines) {
		fmt.Printf("output length mismatch\n")
		os.Exit(1)
	}
	for i := range cLines {
		if !bytes.Equal(bytes.TrimSpace(cLines[i]), bytes.TrimSpace(rLines[i])) {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, rLines[i], cLines[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", caseCount(data))
}
