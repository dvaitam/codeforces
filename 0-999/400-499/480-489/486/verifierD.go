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

const mod = 1000000007

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
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

func bfsConnected(mask int, adj [][]int) bool {
	n := len(adj)
	var start int
	for i := 0; i < n; i++ {
		if mask&(1<<i) != 0 {
			start = i
			break
		}
	}
	vis := make([]bool, n)
	q := []int{start}
	vis[start] = true
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for _, v := range adj[u] {
			if mask&(1<<v) == 0 || vis[v] {
				continue
			}
			vis[v] = true
			q = append(q, v)
		}
	}
	for i := 0; i < n; i++ {
		if mask&(1<<i) != 0 && !vis[i] {
			return false
		}
	}
	return true
}

func expected(d int, a []int, adj [][]int) int {
	n := len(a)
	ans := 0
	for mask := 1; mask < (1 << n); mask++ {
		minV := 3000
		maxV := -1
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				if a[i] < minV {
					minV = a[i]
				}
				if a[i] > maxV {
					maxV = a[i]
				}
			}
		}
		if maxV-minV > d {
			continue
		}
		if bfsConnected(mask, adj) {
			ans++
			if ans >= mod {
				ans -= mod
			}
		}
	}
	return ans
}

func genTree(rng *rand.Rand, n int) [][]int {
	adj := make([][]int, n)
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		adj[i] = append(adj[i], p)
		adj[p] = append(adj[p], i)
	}
	return adj
}

func verifyCase(bin string, d int, a []int, adj [][]int) error {
	n := len(a)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", d, n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	for i := 1; i < n; i++ {
		u := i + 1
		for _, v := range adj[i] {
			if v < i {
				sb.WriteString(fmt.Sprintf("%d %d\n", u, v+1))
			}
		}
	}
	got, err := run(bin, sb.String())
	if err != nil {
		return err
	}
	var res int
	if _, err := fmt.Sscan(got, &res); err != nil {
		return fmt.Errorf("invalid output")
	}
	want := expected(d, a, adj)
	if res%mod != want%mod {
		return fmt.Errorf("expected %d got %d", want%mod, res%mod)
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

	const cases = 100
	for i := 0; i < cases; i++ {
		n := rng.Intn(6) + 2 // 2..7
		d := rng.Intn(5)
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(6)
		}
		adj := genTree(rng, n)
		if err := verifyCase(bin, d, a, adj); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
