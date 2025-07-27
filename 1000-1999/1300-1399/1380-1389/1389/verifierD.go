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

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func solveCase(n int, k int64, l1, r1, l2, r2 int64) int64 {
	if l1 > l2 {
		l1, l2 = l2, l1
		r1, r2 = r2, r1
	}
	intersection := max64(0, min64(r1, r2)-max64(l1, l2))
	unionLen := max64(r1, r2) - min64(l1, l2)
	gap := int64(0)
	if r1 < l2 {
		gap = l2 - r1
	} else if r2 < l1 {
		gap = l1 - r2
	}
	base := int64(n) * intersection
	if k <= base {
		return 0
	}
	need := k - base
	extraWithin := unionLen - intersection
	if gap > 0 {
		extraWithin = unionLen
	}
	ans := int64(1<<62 - 1)
	for i := 1; i <= n; i++ {
		cost := int64(0)
		if gap > 0 {
			cost += gap * int64(i)
		}
		gainLimit := extraWithin * int64(i)
		take := min64(need, gainLimit)
		cost += take
		remain := need - take
		cost += 2 * remain
		if cost < ans {
			ans = cost
		}
	}
	return ans
}

func run(bin, input string) (string, error) {
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

func genTest(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	k := int64(rng.Intn(20) + 1)
	l1 := int64(rng.Intn(10))
	r1 := l1 + int64(rng.Intn(10))
	l2 := int64(rng.Intn(10))
	r2 := l2 + int64(rng.Intn(10))
	input := fmt.Sprintf("1\n%d %d\n%d %d\n%d %d\n", n, k, l1, r1, l2, r2)
	exp := fmt.Sprintf("%d", solveCase(n, k, l1, r1, l2, r2))
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, exp := genTest(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
