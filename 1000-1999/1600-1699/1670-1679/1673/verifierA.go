package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	s string
}

func generateTests() []testCase {
	r := rand.New(rand.NewSource(1))
	tests := []testCase{{"a"}, {"z"}, {"aa"}, {"ab"}, {"abc"}}
	for len(tests) < 100 {
		n := r.Intn(10) + 1
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			b[i] = byte('a' + r.Intn(26))
		}
		tests = append(tests, testCase{s: string(b)})
	}
	return tests
}

func expected(s string) string {
	n := len(s)
	total := 0
	for i := 0; i < n; i++ {
		total += int(s[i]-'a') + 1
	}
	if n == 1 {
		return fmt.Sprintf("Bob %d", total)
	}
	if n%2 == 0 {
		return fmt.Sprintf("Alice %d", total)
	}
	first := int(s[0]-'a') + 1
	last := int(s[n-1]-'a') + 1
	bob := first
	if last < bob {
		bob = last
	}
	alice := total - bob
	if alice > bob {
		return fmt.Sprintf("Alice %d", alice-bob)
	}
	return fmt.Sprintf("Bob %d", bob-alice)
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("run failed: %v\n%s", err, errb.String())
	}
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("1\n%s\n", t.s)
		want := expected(t.s)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.ToUpper(strings.TrimSpace(got)) != strings.ToUpper(want) {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, want, strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
