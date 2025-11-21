package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	input  string
	expect string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA5.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", i+1, err, tc.input)
			os.Exit(1)
		}
		if err := check(tc.expect, strings.TrimSpace(out)); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, err, tc.input, tc.expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func check(expect, actual string) error {
	if expectedInt(expect) != actualInt(actual) {
		return fmt.Errorf("expected %s but got %s", expect, actual)
	}
	return nil
}

func expectedInt(s string) int {
	val, _ := strconv.Atoi(strings.TrimSpace(s))
	return val
}

func actualInt(s string) int {
	val, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return math.MaxInt
	}
	return val
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeTest(0.25*math.Pi, 0),
		makeTest(0.75*math.Pi, 1),
	}
	for i := 0; i < 200; i++ {
		theta := rand.Float64()*0.98*math.Pi + 0.01*math.Pi
		gate := rand.Intn(2)
		tests = append(tests, makeTest(theta, gate))
	}
	return tests
}

func makeTest(theta float64, gate int) testCase {
	input := fmt.Sprintf("%.10f %d\n", theta, gate)
	expect := fmt.Sprintf("%d", gate)
	return testCase{
		input:  input,
		expect: expect,
	}
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}
