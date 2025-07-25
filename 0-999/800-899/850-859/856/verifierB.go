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

func generateCaseB(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		l := rng.Intn(5) + 1
		for j := 0; j < l; j++ {
			sb.WriteByte(byte('a' + rng.Intn(3)))
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCaseB(bin, ref, input string) error {
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref := "./refB.bin"
	if err := buildRef("856B.go", ref); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseB(rng)
		if err := runCaseB(bin, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
