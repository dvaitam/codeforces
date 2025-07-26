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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(8) + 2
	q := rng.Intn(8) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		fmt.Fprintf(&sb, "%d ", p)
	}
	sb.WriteByte('\n')
	cnt := 0
	for i := 0; i < q; i++ {
		t := rng.Intn(3) + 1
		v := rng.Intn(n) + 1
		fmt.Fprintf(&sb, "%d %d\n", t, v)
		if t == 3 {
			cnt++
		}
	}
	return sb.String(), cnt
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, outLines := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		lines := strings.Fields(out)
		if len(lines) != outLines {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d outputs got %d\ninput:\n%s", i+1, outLines, len(lines), in)
			os.Exit(1)
		}
		for _, s := range lines {
			lw := strings.ToLower(strings.TrimSpace(s))
			if lw != "black" && lw != "white" {
				fmt.Fprintf(os.Stderr, "case %d failed: bad output %q\ninput:\n%s", i+1, s, in)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
