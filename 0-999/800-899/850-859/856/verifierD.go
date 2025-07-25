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

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef(src, out string) error {
	cmd := exec.Command("go", "build", "-o", out, src)
	return cmd.Run()
}

func generateCaseD(rng *rand.Rand) string {
	n := rng.Intn(6) + 3
	m := rng.Intn(5)
	parents := make([]int, n+1)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		parents[i] = p
		if i > 2 {
			sb.WriteByte(' ')
		}
		if i == 2 {
			fmt.Fprintf(&sb, "%d", p)
		} else {
			fmt.Fprintf(&sb, "%d", p)
		}
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		for v == u {
			v = rng.Intn(n) + 1
		}
		w := rng.Intn(10) + 1
		fmt.Fprintf(&sb, "%d %d %d\n", u, v, w)
	}
	return sb.String()
}

func runCaseD(bin, ref, input string) error {
	expect, err := runBinary(ref, input)
	if err != nil {
		return fmt.Errorf("reference error: %v", err)
	}
	got, err := runBinary(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, got)
	}
	if strings.TrimSpace(got) != strings.TrimSpace(expect) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expect), strings.TrimSpace(got))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref := "./refD.bin"
	if err := buildRef("856D.go", ref); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseD(rng)
		if err := runCaseD(bin, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
