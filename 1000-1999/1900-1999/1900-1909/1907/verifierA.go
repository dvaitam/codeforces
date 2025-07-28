package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildOfficial() (string, error) {
	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("1907A_official_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", tmp, "1907A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build official solution: %v %s", err, string(out))
	}
	return tmp, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	official, err := buildOfficial()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(official)

	// generate 100 positions
	cols := "abcdefgh"
	var positions []string
	for _, c := range cols {
		for r := 1; r <= 8; r++ {
			positions = append(positions, fmt.Sprintf("%c %d", c, r))
		}
	}
	for len(positions) < 100 {
		positions = append(positions, positions[len(positions)%64])
	}

	groups := [][]string{positions[:64], positions[64:]}
	total := 0
	for i, g := range groups {
		input := fmt.Sprintf("%d\n", len(g))
		for _, p := range g {
			input += p + "\n"
		}
		exp, err := runBinary(official, input)
		if err != nil {
			fmt.Printf("official solution failed: %v\n", err)
			os.Exit(1)
		}
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("binary failed on group %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(out) {
			fmt.Printf("mismatch on group %d\nexpected:\n%s\nactual:\n%s\n", i+1, exp, out)
			os.Exit(1)
		}
		total += len(g)
	}
	fmt.Printf("All %d test cases passed.\n", total)
}
