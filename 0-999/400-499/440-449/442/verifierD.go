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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func solveD(r *bufio.Reader) string {
	var n int
	fmt.Fscan(r, &n)
	parent := make([]int, n+2)
	for i := 1; i <= n; i++ {
		fmt.Fscan(r, &parent[i+1])
	}
	N := n + 1
	size := make([]int, N+1)
	dp := make([]int, N+1)
	bestNon := make([]int, N+1)
	heavy1 := make([]int, N+1)
	heavy2 := make([]int, N+1)
	for i := 1; i <= N; i++ {
		bestNon[i] = -1
	}
	var buf strings.Builder
	for i := 1; i <= n; i++ {
		v := i + 1
		size[v] = 1
		dp[v] = 0
		prev := v
		u := parent[v]
		for u != 0 {
			size[u]++
			heavyChanged := false
			if u == 1 {
				if heavy1[u] == 0 {
					heavy1[u] = prev
					heavyChanged = true
				} else if heavy2[u] == 0 {
					if size[prev] > size[heavy1[u]] {
						heavy2[u] = heavy1[u]
						heavy1[u] = prev
					} else {
						heavy2[u] = prev
					}
					heavyChanged = true
				} else {
					var minh int
					if size[heavy1[u]] < size[heavy2[u]] {
						minh = heavy1[u]
					} else {
						minh = heavy2[u]
					}
					if size[prev] > size[minh] {
						if minh == heavy1[u] {
							heavy1[u] = prev
						} else {
							heavy2[u] = prev
						}
						heavyChanged = true
						bestNon[u] = max(bestNon[u], dp[minh])
					} else {
						bestNon[u] = max(bestNon[u], dp[prev])
					}
				}
			} else {
				if heavy1[u] == 0 {
					heavy1[u] = prev
					heavyChanged = true
				} else if size[prev] > size[heavy1[u]] {
					old := heavy1[u]
					heavy1[u] = prev
					heavyChanged = true
					bestNon[u] = max(bestNon[u], dp[old])
				} else {
					bestNon[u] = max(bestNon[u], dp[prev])
				}
			}
			oldDp := dp[u]
			bestH := -1
			if heavy1[u] != 0 {
				bestH = max(bestH, dp[heavy1[u]])
			}
			if u == 1 && heavy2[u] != 0 {
				bestH = max(bestH, dp[heavy2[u]])
			}
			ndp := bestH
			if u == 1 {
				ndp = max(ndp, bestNon[u])
			} else {
				if bestNon[u] >= 0 {
					ndp = max(ndp, bestNon[u]+1)
				}
			}
			if ndp < 0 {
				ndp = 0
			}
			dp[u] = ndp
			if dp[u] == oldDp && !heavyChanged {
				break
			}
			prev = u
			u = parent[u]
		}
		buf.WriteString(fmt.Sprintf("%d", dp[1]+1))
		if i < n {
			buf.WriteByte(' ')
		}
	}
	buf.WriteByte('\n')
	return buf.String()
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 1; i <= n; i++ {
		if i > 1 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", rng.Intn(i)+1)
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		expect := solveD(bufio.NewReader(strings.NewReader(tc)))
		out, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %sinput:\n%s", i+1, expect, out, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
