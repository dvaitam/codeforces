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
	x int64
	y int64
}

func expected(x, y int64) string {
	if x-y == 1 {
		return "NO"
	}
	return "YES"
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 100)
	// edge cases around x-y==1
	for i := int64(1); i <= 50; i++ {
		tests = append(tests, testCase{x: i + 1, y: i}) // difference 1 -> NO
	}
	for i := int64(1); i <= 50; i++ {
		x := rng.Int63n(1e6) + 2
		y := rng.Int63n(x - 1)
		if x-y == 1 {
			x++ // ensure difference not 1
		}
		tests = append(tests, testCase{x: x, y: y})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("1\n%d %d\n", t.x, t.y)
		want := expected(t.x, t.y)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: x=%d y=%d expected %s got %s\n", i+1, t.x, t.y, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
