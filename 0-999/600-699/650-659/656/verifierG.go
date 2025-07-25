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

func expected(opinions []string, T int) string {
	if len(opinions) == 0 {
		return "0"
	}
	I := len(opinions[0])
	counts := make([]int, I)
	for _, s := range opinions {
		for j, ch := range s {
			if ch == 'Y' {
				counts[j]++
			}
		}
	}
	ans := 0
	for _, c := range counts {
		if c >= T {
			ans++
		}
	}
	return fmt.Sprintf("%d", ans)
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

func genCase(rng *rand.Rand) ([]string, int, int, int) {
	F := rng.Intn(10) + 1
	I := rng.Intn(10) + 1
	T := rng.Intn(F) + 1
	opinions := make([]string, F)
	for i := 0; i < F; i++ {
		var sb strings.Builder
		for j := 0; j < I; j++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('Y')
			} else {
				sb.WriteByte('N')
			}
		}
		opinions[i] = sb.String()
	}
	return opinions, F, I, T
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	const cases = 100
	for i := 0; i < cases; i++ {
		opinions, F, I, T := genCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", F, I, T))
		for _, s := range opinions {
			sb.WriteString(s)
			sb.WriteByte('\n')
		}
		input := sb.String()
		want := expected(opinions, T)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected %s\ngot %s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
