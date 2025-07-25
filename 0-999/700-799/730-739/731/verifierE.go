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

func solveE(a []int64) string {
	n := len(a) - 1
	pref := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] + a[i]
	}
	dp := make([]int64, n+2)
	mx := make([]int64, n+3)
	const inf = int64(1) << 60
	dp[n] = 0
	mx[n+1] = -inf
	mx[n] = pref[n] - dp[n]
	for j := n - 1; j >= 1; j-- {
		dp[j] = mx[j+1]
		v := pref[j] - dp[j]
		if v > mx[j+1] {
			mx[j] = v
		} else {
			mx[j] = mx[j+1]
		}
	}
	return fmt.Sprintf("%d\n", dp[1])
}

func genCaseE(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	arr := make([]int64, n+1)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 1; i <= n; i++ {
		arr[i] = int64(rng.Intn(21) - 10)
		if i > 1 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	return sb.String(), solveE(arr)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, expect := genCaseE(rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))))
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, in, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
