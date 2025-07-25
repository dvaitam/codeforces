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

func buildRef() (string, error) {
	ref := "refB.bin"
	cmd := exec.Command("go", "build", "-o", ref, "727B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randName(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(rng.Intn(26) + 'a')
	}
	return string(b)
}

func formatPrice(cents int64) string {
	dollars := cents / 100
	centsPart := cents % 100
	ds := fmt.Sprintf("%d", dollars)
	var sb strings.Builder
	n := len(ds)
	first := n % 3
	if first == 0 {
		first = 3
	}
	sb.WriteString(ds[:first])
	for i := first; i < n; i += 3 {
		sb.WriteByte('.')
		sb.WriteString(ds[i : i+3])
	}
	if centsPart > 0 {
		sb.WriteByte('.')
		sb.WriteString(fmt.Sprintf("%02d", centsPart))
	}
	return sb.String()
}

func genCase(rng *rand.Rand) string {
	var sb strings.Builder
	length := 0
	for {
		nameLen := rng.Intn(5) + 1
		priceCents := rng.Int63n(1_000_000_00) + 1
		price := formatPrice(priceCents)
		name := randName(rng, nameLen)
		if length+len(name)+len(price) > 950 && length > 0 {
			break
		}
		sb.WriteString(name)
		sb.WriteString(price)
		length += len(name) + len(price)
		if length > 900 || rng.Intn(3) == 0 {
			break
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCase(rng)
		exp, err := run(ref, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out, err := run(candidate, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d wrong answer\nexpected:\n%s\ngot:\n%s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
