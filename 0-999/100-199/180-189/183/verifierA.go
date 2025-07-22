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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 1
	opts := []string{"UL", "UR", "DL", "DR", "ULDR"}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	cntUL, cntUR, cntDL, cntDR := 0, 0, 0, 0
	for i := 0; i < n; i++ {
		s := opts[rng.Intn(len(opts))]
		fmt.Fprintln(&sb, s)
		switch s {
		case "UL":
			cntUL++
		case "UR":
			cntUR++
		case "DL":
			cntDL++
		case "DR":
			cntDR++
		}
	}
	total := int64(n)
	freeU := total - int64(cntUR) - int64(cntDL)
	freeV := total - int64(cntUL) - int64(cntDR)
	ans := (freeU + 1) * (freeV + 1)
	return sb.String(), fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
