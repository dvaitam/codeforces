package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseA struct {
	format int
	time   string
}

func genTestsA() []testCaseA {
	rand.Seed(42)
	tests := make([]testCaseA, 100)
	for i := range tests {
		f := 12
		if rand.Intn(2) == 0 {
			f = 24
		}
		hh1 := rand.Intn(10)
		hh2 := rand.Intn(10)
		mm1 := rand.Intn(10)
		mm2 := rand.Intn(10)
		t := fmt.Sprintf("%d%d:%d%d", hh1, hh2, mm1, mm2)
		tests[i] = testCaseA{f, t}
	}
	return tests
}

func solveA(tc testCaseA) string {
	b := []byte(tc.time)
	if b[3] > '5' {
		b[3] = '0'
	}
	if tc.format == 12 {
		if b[0] != '1' && b[1] == '0' {
			b[0] = '1'
		} else if b[0] > '1' || (b[0] == '1' && b[1] > '2') {
			b[0] = '0'
		}
	} else {
		if b[0] > '2' || (b[0] == '2' && b[1] > '3') {
			b[0] = '0'
		}
	}
	return string(b)
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsA()
	for i, tc := range tests {
		input := fmt.Sprintf("%d\n%s\n", tc.format, tc.time)
		exp := solveA(tc)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != exp {
			fmt.Printf("test %d failed: expected %q got %q\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
