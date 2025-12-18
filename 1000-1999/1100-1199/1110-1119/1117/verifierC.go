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

func abs(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func solveReference(x1, y1, x2, y2 int64, n int64, s string) int64 {
	type pair struct {
		x, y int64
	}
	w := make([]pair, n+1)
	currX, currY := int64(0), int64(0)
	for i := 0; i < int(n); i++ {
		switch s[i] {
		case 'U':
			currY++
		case 'D':
			currY--
		case 'L':
			currX--
		case 'R':
			currX++
		}
		w[i+1] = pair{currX, currY}
	}

	targetX := x2 - x1
	targetY := y2 - y1

	check := func(k int64) bool {
		cnt := k / n
		rem := k % n

		windX := w[n].x*cnt + w[rem].x
		windY := w[n].y*cnt + w[rem].y

		dx := targetX - windX
		dy := targetY - windY

		return abs(dx)+abs(dy) <= k
	}

	var low int64 = 1
	var high int64 = 2000000000000000 // 2 * 10^15
	var ans int64 = -1

	if check(high) {
		for low <= high {
			mid := (low + high) / 2
			if check(mid) {
				ans = mid
				high = mid - 1
			} else {
				low = mid + 1
			}
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, int64) {
	// Use slightly larger range for testing robustness, though original was small.
	rangeVal := 100
	x1 := int64(rng.Intn(2*rangeVal+1) - rangeVal)
	y1 := int64(rng.Intn(2*rangeVal+1) - rangeVal)
	x2 := int64(rng.Intn(2*rangeVal+1) - rangeVal)
	y2 := int64(rng.Intn(2*rangeVal+1) - rangeVal)
	
	// Ensure start != end as per problem statement guarantee
	for x1 == x2 && y1 == y2 {
		x2 = int64(rng.Intn(2*rangeVal+1) - rangeVal)
	}

	n := int64(rng.Intn(100) + 1)
	var sBuilder strings.Builder
	moves := []byte{'U', 'D', 'L', 'R'}
	for i := int64(0); i < n; i++ {
		sBuilder.WriteByte(moves[rng.Intn(4)])
	}
	s := sBuilder.String()

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n%d %d\n%d\n%s\n", x1, y1, x2, y2, n, s)
	
	return sb.String(), solveReference(x1, y1, x2, y2, n, s)
}

func runCase(bin, input string, exp int64) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("failed to parse output: %v\nOutput: %q", err, out.String())
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	// Verify the problem failing case first
	// input:
	// 5 0 -6 10 8
	// LRULUULU
	// Note: user input format reading handles spaces/newlines flexibly. 
	// We will supply it in the strict format we decided on.
	failingInput := "5 0\n-6 10\n8\nLRULUULU\n"
	if err := runCase(os.Args[1], failingInput, 13); err != nil {
		fmt.Fprintf(os.Stderr, "Regression test (known failing case) failed: %v\n", err)
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}