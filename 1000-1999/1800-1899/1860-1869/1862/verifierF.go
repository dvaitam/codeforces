package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func can(t int64, w, f int64, sum int, dp []bool) bool {
	maxW := w * t
	if maxW > int64(sum) {
		maxW = int64(sum)
	}
	maxF := f * t
	for i := int64(0); i <= maxW; i++ {
		if dp[i] {
			if int64(sum)-i <= maxF {
				return true
			}
		}
	}
	return false
}

func solveF(w, f int64, s []int) string {
	sum := 0
	for _, v := range s {
		sum += v
	}
	dp := make([]bool, sum+1)
	dp[0] = true
	for _, v := range s {
		for j := sum; j >= v; j-- {
			if dp[j-v] {
				dp[j] = true
			}
		}
	}
	low := int64(0)
	high := int64(sum)
	if w < f {
		if int64(sum)%w != 0 {
			if int64(sum)/w+1 > high {
				high = int64(sum)/w + 1
			}
		} else {
			if int64(sum)/w > high {
				high = int64(sum) / w
			}
		}
	} else {
		if int64(sum)%f != 0 {
			if int64(sum)/f+1 > high {
				high = int64(sum)/f + 1
			}
		} else {
			if int64(sum)/f > high {
				high = int64(sum) / f
			}
		}
	}
	for low < high {
		mid := (low + high) / 2
		if can(mid, w, f, sum, dp) {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return fmt.Sprint(low)
}

func genCases() []string {
	rand.Seed(6)
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		w := int64(rand.Intn(5) + 1)
		f := int64(rand.Intn(5) + 1)
		n := rand.Intn(6) + 1
		s := make([]int, n)
		for j := 0; j < n; j++ {
			s[j] = rand.Intn(10) + 1
		}
		sb := strings.Builder{}
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", w, f))
		sb.WriteString(fmt.Sprint(n))
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(s[j]))
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
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, tc := range cases {
		lines := strings.Split(strings.TrimSpace(tc), "\n")
		var w, f int64
		fmt.Sscan(lines[1], &w, &f)
		var n int
		fmt.Sscan(lines[2], &n)
		parts := strings.Fields(lines[3])
		s := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Sscan(parts[j], &s[j])
		}
		want := solveF(w, f, s)
		got, err := runCase(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "Wrong answer on case %d\nInput:\n%sExpected: %s Got: %s\n", i+1, tc, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
