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

func expectedAnswerB(s string) string {
	digits := make([]int, len(s))
	maxd := 0
	for i, c := range []byte(s) {
		d := int(c - '0')
		digits[i] = d
		if d > maxd {
			maxd = d
		}
	}
	results := make([]int, 0, maxd)
	for k := 0; k < maxd; k++ {
		num := 0
		for i := range digits {
			num *= 10
			if digits[i] > 0 {
				num++
				digits[i]--
			}
		}
		results = append(results, num)
	}
	nums := make([]string, len(results))
	for i, v := range results {
		nums[i] = fmt.Sprint(v)
	}
	return fmt.Sprintf("%d\n%s", maxd, strings.Join(nums, " "))
}

func generateCaseB(rng *rand.Rand) string {
	n := rng.Intn(12) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('0' + rng.Intn(10))
	}
	if b[0] == '0' {
		b[0] = '1'
	}
	return string(b)
}

func runCaseB(bin, num string) error {
	input := num + "\n"
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := expectedAnswerB(num)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		num := generateCaseB(rng)
		if err := runCaseB(bin, num); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n", i+1, err, num)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
