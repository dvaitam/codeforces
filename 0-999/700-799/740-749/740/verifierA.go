package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(reader *bufio.Reader) string {
	var n, a, b, c int64
	if _, err := fmt.Fscan(reader, &n, &a, &b, &c); err != nil {
		return ""
	}
	const INF int64 = math.MaxInt64 / 4
	dp := [4]int64{0, INF, INF, INF}
	for it := 0; it < 4; it++ {
		old := dp
		for r := 0; r < 4; r++ {
			curr := old[r]
			if curr >= INF {
				continue
			}
			nr := (r + 1) & 3
			if val := curr + a; val < dp[nr] {
				dp[nr] = val
			}
			nr2 := (r + 2) & 3
			if val := curr + b; val < dp[nr2] {
				dp[nr2] = val
			}
			nr3 := (r + 3) & 3
			if val := curr + c; val < dp[nr3] {
				dp[nr3] = val
			}
		}
	}
	rem := int((4 - n%4) % 4)
	return fmt.Sprint(dp[rem])
}

func generateCase(rng *rand.Rand) string {
	n := rng.Int63n(1_000_000_000) + 1
	a := rng.Int63n(1_000_000_000) + 1
	b := rng.Int63n(1_000_000_000) + 1
	c := rng.Int63n(1_000_000_000) + 1
	return fmt.Sprintf("%d %d %d %d\n", n, a, b, c)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		expect := solve(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
