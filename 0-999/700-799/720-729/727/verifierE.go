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
	ref := "refE.bin"
	cmd := exec.Command("go", "build", "-o", ref, "727E.go")
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

func randString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(rng.Intn(26) + 'a')
	}
	return string(b)
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	k := rng.Intn(3) + 1
	g := rng.Intn(5) + n
	names := make([]string, g)
	used := make(map[string]bool)
	for i := 0; i < g; i++ {
		for {
			s := randString(rng, k)
			if !used[s] {
				used[s] = true
				names[i] = s
				break
			}
		}
	}
	perm := rng.Perm(g)[:n]
	cd := make([]string, n)
	for i := 0; i < n; i++ {
		cd[i] = names[perm[i]]
	}
	big := strings.Join(cd, "")
	shift := rng.Intn(n * k)
	big = big[shift:] + big[:shift]
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	sb.WriteString(big)
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", g)
	for i := 0; i < g; i++ {
		sb.WriteString(names[i])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d wrong answer\nexpected:\n%s\ngot:\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
