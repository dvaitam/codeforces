package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod int64 = 1000000007

func solveE(arr []int64) string {
	n := int64(len(arr))
	cnt := make([]int64, 60)
	for _, x := range arr {
		for b := 0; b < 60; b++ {
			if (x>>b)&1 == 1 {
				cnt[b]++
			}
		}
	}
	pow2 := make([]int64, 60)
	for b := 0; b < 60; b++ {
		pow2[b] = (1 << b) % mod
	}
	result := int64(0)
	for _, x := range arr {
		andVal := int64(0)
		orVal := int64(0)
		for b := 0; b < 60; b++ {
			if (x>>b)&1 == 1 {
				andVal = (andVal + cnt[b]*pow2[b]) % mod
				orVal = (orVal + n*pow2[b]) % mod
			} else {
				orVal = (orVal + cnt[b]*pow2[b]) % mod
			}
		}
		result = (result + andVal*orVal) % mod
	}
	return fmt.Sprint(result)
}

func genCases() []string {
	rand.Seed(5)
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(5) + 1
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			arr[j] = int64(rand.Intn(1 << 20))
		}
		sb := strings.Builder{}
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprint(n))
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(arr[j]))
		}
		sb.WriteByte('\n')
		cases[i] = sb.String()
	}
	return cases
}

func runCase(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, tc := range cases {
		lines := strings.Split(strings.TrimSpace(tc), "\n")
		var n int
		fmt.Sscan(lines[1], &n)
		arr := make([]int64, n)
		parts := strings.Fields(lines[2])
		for j := 0; j < n; j++ {
			fmt.Sscan(parts[j], &arr[j])
		}
		want := solveE(arr)
		got, err := runCase(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "Wrong answer on case %d\nInput:\n%sExpected: %s Got: %s\n", i+1, tc, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
