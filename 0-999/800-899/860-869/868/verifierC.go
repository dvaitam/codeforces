package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	in  string
	out string
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func solveOracle(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return ""
	}
	maxMask := 1 << k
	exists := make([]bool, maxMask)
	for i := 0; i < n; i++ {
		mask := 0
		for j := 0; j < k; j++ {
			var x int
			fmt.Fscan(in, &x)
			if x == 1 {
				mask |= 1 << j
			}
		}
		exists[mask] = true
	}
	masks := make([]int, 0)
	for m := 0; m < maxMask; m++ {
		if exists[m] {
			masks = append(masks, m)
		}
	}
	m := len(masks)
	for subset := 1; subset < (1 << m); subset++ {
		total := 0
		counts := make([]int, k)
		for i := 0; i < m; i++ {
			if subset&(1<<i) != 0 {
				total++
				mask := masks[i]
				for j := 0; j < k; j++ {
					if mask&(1<<j) != 0 {
						counts[j]++
					}
				}
			}
		}
		valid := true
		for j := 0; j < k; j++ {
			if counts[j]*2 > total {
				valid = false
				break
			}
		}
		if valid {
			return "YES"
		}
	}
	return "NO"
}

func genCase(rng *rand.Rand) Test {
	n := rng.Intn(6) + 1
	k := rng.Intn(4) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i := 0; i < n; i++ {
		for j := 0; j < k; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", rng.Intn(2)))
		}
		sb.WriteByte('\n')
	}
	input := sb.String()
	out := solveOracle(input)
	return Test{input, out}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(3))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		got, err := run(bin, tc.in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
		if got != tc.out {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.out, got, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
