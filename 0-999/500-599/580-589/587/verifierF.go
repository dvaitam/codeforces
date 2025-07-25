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

func randString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func run(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref := "587F.go"
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for caseNum := 0; caseNum < 100; caseNum++ {
		n := rng.Intn(4) + 1
		q := rng.Intn(4) + 1
		strs := make([]string, n)
		for i := 0; i < n; i++ {
			strs[i] = randString(rng, rng.Intn(5)+1)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
		for i := 0; i < n; i++ {
			sb.WriteString(strs[i])
			if i+1 < n {
				sb.WriteByte('\n')
			} else {
				sb.WriteByte('\n')
			}
		}
		for i := 0; i < q; i++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			k := rng.Intn(n) + 1
			sb.WriteString(fmt.Sprintf("%d %d %d\n", l, r, k))
		}
		input := sb.String()
		expected, err := run(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error running reference: %v\n", err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", caseNum+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
