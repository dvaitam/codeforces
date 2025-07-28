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

const mod int64 = 998244353

type testCase struct {
	arr []int64
}

func expected(tc testCase) string {
	n := len(tc.arr)
	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] ^ tc.arr[i-1]
	}
	pow2 := make([]int64, 31)
	pow2[0] = 1
	for i := 1; i < 31; i++ {
		pow2[i] = (pow2[i-1] * 2) % mod
	}
	var ans int64
	for bit := 0; bit < 31; bit++ {
		cnt0, sum0 := int64(1), int64(0)
		cnt1, sum1 := int64(0), int64(0)
		for j := 1; j <= n; j++ {
			if ((prefix[j] >> bit) & 1) == 0 {
				tmp := int64(j)*cnt1 - sum1
				ans = (ans + ((tmp%mod+mod)%mod)*pow2[bit]) % mod
				cnt0++
				sum0 += int64(j)
			} else {
				tmp := int64(j)*cnt0 - sum0
				ans = (ans + ((tmp%mod+mod)%mod)*pow2[bit]) % mod
				cnt1++
				sum1 += int64(j)
			}
		}
	}
	ans %= mod
	if ans < 0 {
		ans += mod
	}
	return fmt.Sprintf("%d", ans)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 102)
	cases = append(cases, testCase{arr: []int64{1}})
	cases = append(cases, testCase{arr: []int64{1, 2, 3}})
	for len(cases) < 102 {
		n := rng.Intn(10) + 1
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Int63n(100)
		}
		cases = append(cases, testCase{arr: arr})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genCases()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(tc.arr)))
		for j := range tc.arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", tc.arr[j]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		want := expected(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
