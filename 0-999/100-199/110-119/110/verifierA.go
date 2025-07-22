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

func expectedNearlyLucky(n uint64) string {
	count := 0
	for tmp := n; tmp > 0; tmp /= 10 {
		d := tmp % 10
		if d == 4 || d == 7 {
			count++
		}
	}
	if count == 0 {
		return "NO"
	}
	for tmp := count; tmp > 0; tmp /= 10 {
		d := tmp % 10
		if d != 4 && d != 7 {
			return "NO"
		}
	}
	return "YES"
}

func runCase(bin string, input string, expected string) error {
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
		log.Fatalf("Usage: go run verifierA.go /path/to/binary")
	}
	bin := os.Args[1]
	var cases []uint64
	for i := 1; i <= 100; i++ {
		cases = append(cases, uint64(i))
	}
	extras := []uint64{444444444444444444, 777777777777777777, 474747474747474747, 123456789876543210}
	cases = append(cases, extras...)

	for idx, n := range cases {
		expect := expectedNearlyLucky(n)
		if err := runCase(bin, strconv.FormatUint(n, 10), expect); err != nil {
			log.Fatalf("case %d (%d) failed: %v", idx+1, n, err)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(cases))
}
