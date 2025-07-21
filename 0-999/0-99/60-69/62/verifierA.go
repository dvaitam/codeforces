package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func canHold(girl, boy int) bool {
	return boy >= girl-1 && boy <= 2*(girl+1)
}

func solveCase(al, ar, bl, br int) string {
	if canHold(al, br) || canHold(ar, bl) {
		return "YES"
	}
	return "NO"
}

func runCase(bin string, al, ar, bl, br int) (string, error) {
	input := fmt.Sprintf("%d %d\n%d %d\n", al, ar, bl, br)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		al := rng.Intn(6)
		ar := rng.Intn(6)
		bl := rng.Intn(6)
		br := rng.Intn(6)
		expected := solveCase(al, ar, bl, br)
		got, err := runCase(bin, al, ar, bl, br)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s (input: %d %d %d %d)\n", i+1, expected, got, al, ar, bl, br)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
