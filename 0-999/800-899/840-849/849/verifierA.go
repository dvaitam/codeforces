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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(arr []int) string {
	n := len(arr)
	if n%2 == 1 && arr[0]%2 == 1 && arr[n-1]%2 == 1 {
		return "Yes"
	}
	return "No"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1 // 1..100
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(101)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), expected(arr)
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []struct{ in, out string }{}
	// deterministic edge cases
	cases = append(cases, struct{ in, out string }{"1\n1\n", "Yes"})
	cases = append(cases, struct{ in, out string }{"2\n1 3\n", "No"})
	cases = append(cases, struct{ in, out string }{"3\n1 2 3\n", "Yes"})
	cases = append(cases, struct{ in, out string }{"3\n2 3 5\n", "No"})

	for i := 0; i < 100; i++ {
		in, out := generateCase(rng)
		cases = append(cases, struct{ in, out string }{in, out})
	}

	for i, tc := range cases {
		got, err := runCandidate(bin, tc.in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(tc.out) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.out, got, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
