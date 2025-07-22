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

func solveB(n int, a []int64) []int64 {
	m := (n - 1) * 2
	const UNVIS = int64(-2)
	const INF = int64(-1)
	dp := make([]int64, m)
	rev := make([][]int, m)
	for i := range dp {
		dp[i] = UNVIS
	}
	q := make([]int, 0, m)
	for pos := 2; pos <= n; pos++ {
		idx := (pos - 2) * 2
		for dir := 0; dir < 2; dir++ {
			u := idx + dir
			var vpos int64
			if dir == 0 {
				vpos = int64(pos) + a[pos]
			} else {
				vpos = int64(pos) - a[pos]
			}
			vdir := dir ^ 1
			if vpos >= 2 && vpos <= int64(n) {
				vidx := (int(vpos) - 2) * 2
				v := vidx + vdir
				rev[v] = append(rev[v], u)
			} else if vpos <= 0 || vpos > int64(n) {
				dp[u] = a[pos]
				q = append(q, u)
			}
		}
	}
	for head := 0; head < len(q); head++ {
		v := q[head]
		for _, u := range rev[v] {
			if dp[u] != UNVIS {
				continue
			}
			pos := u/2 + 2
			if dp[v] == INF {
				dp[u] = INF
			} else {
				dp[u] = a[pos] + dp[v]
			}
			q = append(q, u)
		}
	}
	for i := 0; i < m; i++ {
		if dp[i] == UNVIS {
			dp[i] = INF
		}
	}
	ans := make([]int64, n-1)
	for i := 1; i <= n-1; i++ {
		pos1 := 1 + i
		u := (pos1-2)*2 + 1
		if dp[u] == INF {
			ans[i-1] = -1
		} else {
			ans[i-1] = dp[u] + int64(i)
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, []int64) {
	n := rng.Intn(18) + 2
	a := make([]int64, n+1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 2; i <= n; i++ {
		a[i] = int64(rng.Intn(10) + 1)
		if i > 2 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(a[i]))
	}
	sb.WriteString("\n")
	ans := solveB(n, a)
	expLines := make([]string, len(ans))
	for i, v := range ans {
		expLines[i] = fmt.Sprint(v)
	}
	return sb.String(), ans
}

func runCase(bin, input string, exp []int64) error {
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
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != len(exp) {
		return fmt.Errorf("expected %d lines got %d", len(exp), len(lines))
	}
	for i, line := range lines {
		var val int64
		if _, err := fmt.Sscan(strings.TrimSpace(line), &val); err != nil {
			return fmt.Errorf("bad int on line %d: %v", i+1, err)
		}
		if val != exp[i] {
			return fmt.Errorf("line %d expected %d got %d", i+1, exp[i], val)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
