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

type testCase struct {
	t, s, x int64
}

func expected(tc testCase) string {
	if tc.x == tc.t || (tc.x >= tc.t+tc.s && ((tc.x-tc.t)%tc.s == 0 || (tc.x-tc.t-1)%tc.s == 0)) {
		return "YES"
	}
	return "NO"
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d %d %d\n", tc.t, tc.s, tc.x)
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	want := expected(tc)
	if got != want {
		return fmt.Errorf("expected %s got %s", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	// some fixed edge cases
	cases = append(cases,
		testCase{0, 2, 0},
		testCase{0, 2, 1},
		testCase{0, 2, 2},
		testCase{0, 2, 3},
		testCase{1, 2, 1},
		testCase{1, 2, 2},
		testCase{1, 2, 3},
		testCase{1, 2, 4},
		testCase{5, 3, 8},
		testCase{5, 3, 9},
	)

	// generate random cases
	for i := 0; i < 200; i++ {
		t := rng.Int63n(1000)
		s := rng.Int63n(999) + 2 // s >= 2
		x := rng.Int63n(2000)
		cases = append(cases, testCase{t, s, x})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
