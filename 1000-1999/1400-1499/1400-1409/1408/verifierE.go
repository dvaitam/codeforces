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
	m := rng.Intn(3) + 1
	n := rng.Intn(3) + 1
	a := make([]int, m)
	b := make([]int, n)
	for i := 0; i < m; i++ {
		a[i] = rng.Intn(10) + 1
	}
	for j := 0; j < n; j++ {
		b[j] = rng.Intn(10) + 1
	}
	sets := make([][]int, m)
	for i := 0; i < m; i++ {
		sz := rng.Intn(n + 1)
		perm := rand.Perm(n)
		sets[i] = make([]int, sz)
		for j := 0; j < sz; j++ {
			sets[i][j] = perm[j] + 1
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", m, n))
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(a[i]))
	}
	sb.WriteByte('\n')
	for j := 0; j < n; j++ {
		if j > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(b[j]))
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		sb.WriteString(fmt.Sprint(len(sets[i])))
		for _, v := range sets[i] {
			sb.WriteByte(' ')
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	refBin := "ref1408E.bin"
	if err := exec.Command("go", "build", "-o", refBin, "1408E.go").Run(); err != nil {
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
