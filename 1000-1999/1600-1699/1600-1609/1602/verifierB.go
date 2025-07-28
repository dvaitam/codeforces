package main

import (
	"bufio"
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(reader *bufio.Reader) string {
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return ""
	}
	var sb strings.Builder
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := range arr {
			fmt.Fscan(reader, &arr[i])
		}
		var q int
		fmt.Fscan(reader, &q)
		states := make([][]int, 0, n+1)
		first := append([]int(nil), arr...)
		states = append(states, first)
		for step := 1; step <= n; step++ {
			prev := states[step-1]
			freq := make([]int, n+1)
			for _, v := range prev {
				freq[v]++
			}
			cur := make([]int, n)
			same := true
			for i, v := range prev {
				cur[i] = freq[v]
				if cur[i] != prev[i] {
					same = false
				}
			}
			states = append(states, cur)
			if same {
				break
			}
		}
		maxStep := len(states) - 1
		for ; q > 0; q-- {
			var x, k int
			fmt.Fscan(reader, &x, &k)
			if k > maxStep {
				k = maxStep
			}
			fmt.Fprintf(&sb, "%d\n", states[k][x-1])
		}
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) string {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for ; t > 0; t-- {
		n := rng.Intn(20) + 1
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rng.Intn(n)+1)
		}
		sb.WriteByte('\n')
		q := rng.Intn(20) + 1
		fmt.Fprintf(&sb, "%d\n", q)
		for i := 0; i < q; i++ {
			x := rng.Intn(n) + 1
			k := rng.Intn(n*2 + 1)
			fmt.Fprintf(&sb, "%d %d\n", x, k)
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
