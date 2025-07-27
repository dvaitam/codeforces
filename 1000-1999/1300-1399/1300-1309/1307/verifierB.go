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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(input string) string {
	in := strings.NewReader(input)
	var t int
	fmt.Fscan(in, &t)
	var out strings.Builder
	for casei := 0; casei < t; casei++ {
		var n int
		var x int64
		fmt.Fscan(in, &n, &x)
		arr := make([]int64, n)
		var maxa int64
		eq := false
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			if arr[i] == x {
				eq = true
			}
			if arr[i] > maxa {
				maxa = arr[i]
			}
		}
		var ans int64
		if eq {
			ans = 1
		} else {
			hops := (x + maxa - 1) / maxa
			if hops < 2 {
				hops = 2
			}
			ans = hops
		}
		fmt.Fprintf(&out, "%d", ans)
		if casei+1 < t {
			out.WriteByte('\n')
		}
	}
	return out.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(5) + 1
		x := int64(rng.Intn(20) + 1)
		fmt.Fprintf(&sb, "%d %d\n", n, x)
		for j := 0; j < n; j++ {
			val := int64(rng.Intn(20) + 1)
			if j+1 == n {
				fmt.Fprintf(&sb, "%d\n", val)
			} else {
				fmt.Fprintf(&sb, "%d ", val)
			}
		}
	}
	in := sb.String()
	return in, solve(in)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
