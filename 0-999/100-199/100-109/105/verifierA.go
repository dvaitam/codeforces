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
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func randomString(rng *rand.Rand, min, max int) string {
	l := rng.Intn(max-min+1) + min
	b := make([]byte, l)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	k := rng.Intn(99) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d 0.%02d\n", n, m, k)
	names := make(map[string]struct{})
	for i := 0; i < n; i++ {
		name := randomString(rng, 1, 5)
		for {
			if _, ok := names[name]; !ok {
				break
			}
			name = randomString(rng, 1, 5)
		}
		names[name] = struct{}{}
		lvl := rng.Intn(10000)
		fmt.Fprintf(&sb, "%s %d\n", name, lvl)
	}
	used := make(map[string]struct{})
	for i := 0; i < m; i++ {
		name := randomString(rng, 1, 5)
		for {
			if _, ok := used[name]; !ok {
				break
			}
			name = randomString(rng, 1, 5)
		}
		used[name] = struct{}{}
		fmt.Fprintf(&sb, "%s\n", name)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	refPath := "./105A.go"
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		want, err := runBinary(refPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal solver failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\ninput:\n%s\n", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
