package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

var ones = map[int]string{
	0: "zero", 1: "one", 2: "two", 3: "three", 4: "four", 5: "five",
	6: "six", 7: "seven", 8: "eight", 9: "nine", 10: "ten",
	11: "eleven", 12: "twelve", 13: "thirteen", 14: "fourteen", 15: "fifteen",
	16: "sixteen", 17: "seventeen", 18: "eighteen", 19: "nineteen",
}

var tens = map[int]string{
	20: "twenty", 30: "thirty", 40: "forty", 50: "fifty",
	60: "sixty", 70: "seventy", 80: "eighty", 90: "ninety",
}

func expectedWord(s int) string {
	if s < 20 {
		return ones[s]
	}
	if s%10 == 0 {
		return tens[s]
	}
	return tens[s-s%10] + "-" + ones[s%10]
}

type testCase struct {
	name  string
	value int
}

func runCandidate(bin string, s int) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%d\n", s))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	var tests []testCase
	for s := 0; s <= 19; s++ {
		tests = append(tests, testCase{fmt.Sprintf("direct_%d", s), s})
	}
	for s := 20; s <= 90; s += 10 {
		tests = append(tests, testCase{fmt.Sprintf("tens_%d", s), s})
	}
	rng := rand.New(rand.NewSource(535))
	for i := 0; i < 30; i++ {
		tests = append(tests, testCase{fmt.Sprintf("rand_%d", i+1), rng.Intn(100)})
	}
	for idx, tc := range tests {
		expect := expectedWord(tc.value)
		got, err := runCandidate(bin, tc.value)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:%d\n", idx+1, tc.name, err, tc.value)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d (%s) failed: expect %q got %q (input %d)\n", idx+1, tc.name, expect, got, tc.value)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
