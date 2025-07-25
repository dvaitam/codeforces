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

type testCaseB struct {
	s string
}

func generateCase(rng *rand.Rand) (string, testCaseB) {
	length := rng.Intn(20) + 1
	var sb strings.Builder
	sb.WriteByte(byte('1' + rng.Intn(9)))
	for i := 1; i < length; i++ {
		sb.WriteByte(byte('0' + rng.Intn(10)))
	}
	s := sb.String()
	return s + "\n", testCaseB{s: s}
}

func expected(tc testCaseB) string {
	var sb strings.Builder
	sb.WriteString(tc.s)
	for i := len(tc.s) - 1; i >= 0; i-- {
		sb.WriteByte(tc.s[i])
	}
	return sb.String()
}

func runCase(bin, input string, tc testCaseB) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	want := expected(tc)
	if got != want {
		return fmt.Errorf("expected %s got %s", want, got)
	}
	return nil
}

func main() {
	var bin string
	switch len(os.Args) {
	case 2:
		bin = os.Args[1]
	case 3:
		if os.Args[1] == "--" {
			bin = os.Args[2]
		}
	}
	if bin == "" {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go [--] /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, tc := generateCase(rng)
		if err := runCase(bin, input, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
