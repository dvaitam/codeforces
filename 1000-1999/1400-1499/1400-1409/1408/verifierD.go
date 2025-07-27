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

func runProgram(bin, input string) (string, error) {
	if _, err := os.Stat(bin); err == nil && !strings.Contains(bin, "/") {
		bin = "./" + bin
	}
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

func randomCase(rng *rand.Rand) string {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	type pt struct{ x, y int }
	a := make([]pt, n)
	b := make([]pt, m)
	for i := 0; i < n; i++ {
		a[i] = pt{rng.Intn(6), rng.Intn(6)}
	}
	for j := 0; j < m; j++ {
		b[j] = pt{rng.Intn(6), rng.Intn(6)}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", a[i].x, a[i].y)
	}
	for j := 0; j < m; j++ {
		fmt.Fprintf(&sb, "%d %d\n", b[j].x, b[j].y)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	refBin := "ref1408D.bin"
	if err := exec.Command("go", "build", "-o", refBin, "1408D.go").Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input := randomCase(rng)
		want, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
			os.Exit(1)
		}
		out, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(want) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", t+1, want, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
