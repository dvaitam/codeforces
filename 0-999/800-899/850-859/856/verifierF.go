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

func generateCaseF(rng *rand.Rand) string {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	C := rng.Int63n(10)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, C)
	for i := 0; i < n; i++ {
		a := rng.Int63n(20)
		b := a + rng.Int63n(5) + 1
		fmt.Fprintf(&sb, "%d %d\n", a, b)
	}
	for i := 0; i < m; i++ {
		c := rng.Int63n(20)
		d := c + rng.Int63n(5) + 1
		fmt.Fprintf(&sb, "%d %d\n", c, d)
	}
	return sb.String()
}

func runCaseF(bin, ref, input string) error {
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref := "./refF.bin"
	if err := buildRef("856F.go", ref); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseF(rng)
		if err := runCaseF(bin, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
