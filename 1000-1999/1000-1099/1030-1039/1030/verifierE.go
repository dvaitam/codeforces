package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func solve(arr []uint64) string {
	n := len(arr)
	c := make([]int, n+1)
	for i := 1; i <= n; i++ {
		c[i] = bits.OnesCount64(arr[i-1])
	}
	cnt := [2]int64{1, 0}
	par := 0
	for i := 1; i <= n; i++ {
		par ^= c[i] & 1
		cnt[par]++
	}
	ans := cnt[0]*(cnt[0]-1)/2 + cnt[1]*(cnt[1]-1)/2
	const K = 128
	for r := 1; r <= n; r++ {
		sum := 0
		mx := 0
		for l := r; l >= 1 && r-l < K; l-- {
			sum += c[l]
			if c[l] > mx {
				mx = c[l]
			}
			if sum%2 == 0 && mx*2 > sum {
				ans--
			}
		}
	}
	return fmt.Sprintf("%d\n", ans)
}

func buildCase(arr []uint64) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(arr)))
	sb.WriteByte('\n')
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatUint(v, 10))
	}
	sb.WriteByte('\n')
	return testCase{input: sb.String(), expected: solve(arr)}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(30) + 1
	arr := make([]uint64, n)
	for i := range arr {
		arr[i] = rng.Uint64() % (1 << 20)
	}
	return buildCase(arr)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if out.String() != tc.expected {
		return fmt.Errorf("expected %q got %q", tc.expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		buildCase([]uint64{1, 2, 3}),
		buildCase([]uint64{7, 7, 7}),
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
