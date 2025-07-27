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

func randomString(rng *rand.Rand, m int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	bs := make([]rune, m)
	for i := range bs {
		bs[i] = letters[rng.Intn(len(letters))]
	}
	return string(bs)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(6) + 1
	m := rng.Intn(4) + 1
	set := make(map[string]struct{})
	arr := make([]string, n)
	for i := 0; i < n; i++ {
		s := randomString(rng, m)
		for {
			if _, ok := set[s]; !ok {
				break
			}
			s = randomString(rng, m)
		}
		set[s] = struct{}{}
		arr[i] = s
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, s := range arr {
		sb.WriteString(s)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	expected := "1304B.go"
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp, err := run(expected, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error running expected solution: %v\n", err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected:\n%s\ngot:\n%s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
