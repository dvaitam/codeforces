package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func solveA(r *bufio.Reader) string {
	var s string
	if _, err := fmt.Fscan(r, &s); err != nil {
		return ""
	}
	parts := strings.SplitN(s, ".", 2)
	intPart := parts[0]
	fracPart := parts[1]
	if intPart[len(intPart)-1] == '9' {
		return "GOTO Vasilisa."
	}
	if fracPart[0] >= '5' {
		b := []byte(intPart)
		b[len(b)-1]++
		intPart = string(b)
	}
	i := 0
	for i < len(intPart)-1 && intPart[i] == '0' {
		i++
	}
	intPart = intPart[i:]
	return intPart
}

func generateCaseA(rng *rand.Rand) string {
	lenInt := rng.Intn(5) + 1
	intPart := make([]byte, lenInt)
	if lenInt == 1 {
		intPart[0] = byte('0' + rng.Intn(10))
	} else {
		intPart[0] = byte('1' + rng.Intn(9))
	}
	for i := 1; i < lenInt; i++ {
		intPart[i] = byte('0' + rng.Intn(10))
	}
	if rng.Intn(5) == 0 {
		intPart[lenInt-1] = '9'
	}
	lenFrac := rng.Intn(5) + 1
	fracPart := make([]byte, lenFrac)
	for i := 0; i < lenFrac; i++ {
		fracPart[i] = byte('0' + rng.Intn(10))
	}
	if rng.Intn(2) == 0 {
		fracPart[0] = byte('0' + rng.Intn(5))
	} else {
		fracPart[0] = byte('5' + rng.Intn(5))
	}
	return fmt.Sprintf("%s.%s\n", string(intPart), string(fracPart))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseA(rng)
		expect := solveA(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
