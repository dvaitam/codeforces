package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Friend struct {
	k    int64
	x    int64
	mask int
}

func solveCase(n, m int, b int64, fs []Friend) int64 {
	sort.Slice(fs, func(i, j int) bool { return fs[i].k < fs[j].k })
	full := (1 << m) - 1
	const INF int64 = 1<<63 - 1
	dp := make([]int64, 1<<m)
	for i := 1; i <= full; i++ {
		dp[i] = INF
	}
	dp[0] = 0
	ans := INF
	for i := 0; i < n; {
		curK := fs[i].k
		j := i
		for j < n && fs[j].k == curK {
			f := fs[j]
			for mask := full; ; mask-- {
				prev := dp[mask]
				if prev < INF {
					nm := mask | f.mask
					cost := prev + f.x
					if cost < dp[nm] {
						dp[nm] = cost
					}
				}
				if mask == 0 {
					break
				}
			}
			j++
		}
		if dp[full] < INF {
			total := dp[full] + curK*b
			if total < ans {
				ans = total
			}
		}
		i = j
	}
	if ans == INF {
		return -1
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(5) + 1
	m := rng.Intn(4) + 1
	b := int64(rng.Intn(5) + 1)
	friends := make([]Friend, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, b)
	for i := 0; i < n; i++ {
		xi := int64(rng.Intn(10) + 1)
		ki := int64(rng.Intn(5) + 1)
		mi := rng.Intn(m) + 1
		probs := rng.Perm(m)[:mi]
		mask := 0
		for _, p := range probs {
			mask |= 1 << p
		}
		friends[i] = Friend{k: ki, x: xi, mask: mask}
		fmt.Fprintf(&sb, "%d %d %d\n", xi, ki, mi)
		for idx, p := range probs {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", p+1)
		}
		sb.WriteByte('\n')
	}
	expect := solveCase(n, m, b, friends)
	return sb.String(), expect
}

func runCase(bin, input string, expected int64) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output")
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
