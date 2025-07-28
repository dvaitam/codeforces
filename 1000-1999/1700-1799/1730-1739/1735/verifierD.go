package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func buildRef() (string, error) {
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1735D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return ref, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	k := rng.Intn(4) + 1
	used := make(map[int64]bool)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	count := 0
	for count < n {
		val := int64(0)
		card := make([]int, k)
		for i := 0; i < k; i++ {
			x := rng.Intn(3)
			card[i] = x
			val = val*3 + int64(x)
		}
		if used[val] {
			continue
		}
		used[val] = true
		for i, x := range card {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", x)
		}
		sb.WriteByte('\n')
		count++
	}
	return sb.String()
}

func genCases() []string {
	rng := rand.New(rand.NewSource(1735))
	cases := make([]string, 0, 100)
	cases = append(cases, "1 1\n0\n")
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	cases := genCases()
	for i, tc := range cases {
		exp, err := runBinary(ref, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on case %d: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "wrong answer on case %d\ninput:\n%sexpected: %s\ngot: %s\n", i+1, tc, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed!\n", len(cases))
}
