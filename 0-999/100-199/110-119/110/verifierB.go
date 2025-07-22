package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func expectedString(n int) string {
	letters := []byte{'a', 'b', 'c', 'd'}
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		res[i] = letters[i%4]
	}
	return string(res)
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input + "\n")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execution error: %v (output: %s)", err, out.String())
	}
	result := strings.TrimSpace(out.String())
	if result != expected {
		return fmt.Errorf("expected %q got %q", expected, result)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run verifierB.go /path/to/binary")
	}
	bin := os.Args[1]
	var cases []int
	for i := 1; i <= 100; i++ {
		cases = append(cases, i)
	}
	extras := []int{101, 256, 512, 1000, 100000}
	cases = append(cases, extras...)

	for idx, n := range cases {
		expect := expectedString(n)
		if err := runCase(bin, strconv.Itoa(n), expect); err != nil {
			log.Fatalf("case %d (%d) failed: %v", idx+1, n, err)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(cases))
}
