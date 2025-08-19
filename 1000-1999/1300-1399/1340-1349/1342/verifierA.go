package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

type testCase struct {
	x, y, a, b int64
}

func (tc testCase) Input() string {
	return fmt.Sprintf("1\n%d %d\n%d %d\n", tc.x, tc.y, tc.a, tc.b)
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func expected(tc testCase) int64 {
	x, y, a, b := tc.x, tc.y, tc.a, tc.b
	if x < y {
		x, y = y, x
	}
	diff := x - y
	pair := b
	if pair > 2*a {
		pair = 2 * a
	}
	return diff*a + y*pair
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(1))
	tests := make([]testCase, 100)
	for i := range tests {
		tests[i] = testCase{
			x: rng.Int63n(1e9 + 1),
			y: rng.Int63n(1e9 + 1),
			a: rng.Int63n(1e9) + 1,
			b: rng.Int63n(1e9) + 1,
		}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		want := expected(tc)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\ninput:\n%s", i+1, err, input)
			return
		}
		got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Printf("case %d: non-integer output %q\n", i+1, out)
			return
		}
		if got != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %d\ngot: %d\n", i+1, input, want, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
