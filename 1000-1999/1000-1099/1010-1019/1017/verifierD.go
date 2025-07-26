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
	n := rng.Intn(4) + 1
	m := rng.Intn(3) + 1
	q := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, q)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d ", rng.Intn(11))
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		var s strings.Builder
		for j := 0; j < n; j++ {
			s.WriteByte(byte('0' + rng.Intn(2)))
		}
		fmt.Fprintf(&sb, "%s\n", s.String())
	}
	for i := 0; i < q; i++ {
		var t strings.Builder
		for j := 0; j < n; j++ {
			t.WriteByte(byte('0' + rng.Intn(2)))
		}
		fmt.Fprintf(&sb, "%s %d\n", t.String(), rng.Intn(11))
	}
	return sb.String(), q
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, q := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		lines := strings.Fields(out)
		if len(lines) != q {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\ninput:\n%s", i+1, q, len(lines), in)
			os.Exit(1)
		}
		for _, v := range lines {
			if _, err := fmt.Sscan(v, new(int)); err != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: bad number %q\ninput:\n%s", i+1, v, in)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
