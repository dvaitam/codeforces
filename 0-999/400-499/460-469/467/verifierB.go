package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
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
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveB(r *bufio.Reader) string {
	var n, m, k int
	fmt.Fscan(r, &n, &m, &k)
	armies := make([]int, m+1)
	for i := 0; i < m+1; i++ {
		fmt.Fscan(r, &armies[i])
	}
	fedor := armies[m]
	friends := 0
	for i := 0; i < m; i++ {
		if bits.OnesCount(uint(fedor^armies[i])) <= k {
			friends++
		}
	}
	return fmt.Sprint(friends)
}

func generateCaseB(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	m := rng.Intn(20) + 1
	k := rng.Intn(n + 1)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	maxVal := 1 << n
	for i := 0; i < m+1; i++ {
		val := rng.Intn(maxVal)
		fmt.Fprintf(&sb, "%d ", val)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseB(rng)
		expect := solveB(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
