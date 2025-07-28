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

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v\nstderr: %s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		a := rng.Intn(1000) - 500
		b := rng.Intn(1000) - 500
		input := fmt.Sprintf("%d %d\n", a, b)
		expect := fmt.Sprintf("%d", a*b)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i, err)
			os.Exit(1)
		}
		if out != expect {
			fmt.Printf("case %d failed: expected %s got %s\ninput: %s", i, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
