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

const mod int64 = 1000000007

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

// bruteForce computes the answer for small inputs.
// b[i] = a[i mod n] for 0 <= i < l.
// Count subsequences bi1..bix with:
//   1 <= x <= k
//   floor(ij/n)+1 = floor(ij+1/n) for consecutive elements (consecutive periods)
//   non-decreasing values
func bruteForce(n int, l int64, k int, a []int) int64 {
	L := int(l)
	b := make([]int, L)
	for i := 0; i < L; i++ {
		b[i] = a[i%n]
	}

	// dp[i][length] = number of valid subsequences of given length ending at position i
	// Only allocate for positions we need
	dp := make([][]int64, L)
	for i := 0; i < L; i++ {
		maxLen := k
		if maxLen > L {
			maxLen = L
		}
		dp[i] = make([]int64, maxLen+1)
		dp[i][1] = 1
	}

	ans := int64(L) % mod // all length-1 subsequences

	maxK := k
	if maxK > L {
		maxK = L
	}
	for length := 2; length <= maxK; length++ {
		for j := 0; j < L; j++ {
			pj := j / n
			for i := 0; i < j; i++ {
				pi := i / n
				if pi+1 != pj {
					continue
				}
				if b[i] > b[j] {
					continue
				}
				dp[j][length] = (dp[j][length] + dp[i][length-1]) % mod
			}
			ans = (ans + dp[j][length]) % mod
		}
	}

	return ans % mod
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for caseNum := 0; caseNum < 100; caseNum++ {
		n := rng.Intn(4) + 1
		k := rng.Intn(4) + 1
		if n*k > 1000000 {
			k = 1000000 / n
		}
		// Keep l small enough for brute force but large enough to test multi-period logic
		// l should be at least n to have one full period, and up to about 5*n
		l := int64(n) + int64(rng.Intn(4*n+1))
		if l > 24 {
			l = 24
		}
		if l < 1 {
			l = 1
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(5) + 1
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, l, k))
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteString("\n")
		input := sb.String()

		expected := bruteForce(n, l, k, arr)

		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
		if got != fmt.Sprintf("%d", expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", caseNum+1, expected, got, input)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed")
}
