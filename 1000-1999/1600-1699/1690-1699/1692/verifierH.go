package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func genTest() []byte {
	t := rand.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rand.Intn(20) + 1
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(rand.Intn(50) + 1))
		}
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func check(input []byte, output string) error {
	in := bufio.NewReader(bytes.NewReader(input))
	out := bufio.NewReader(strings.NewReader(output))
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return fmt.Errorf("failed to read t: %v", err)
	}
	for tc := 1; tc <= t; tc++ {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return fmt.Errorf("failed to read n in test %d: %v", tc, err)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		// compute maximum possible gain
		maxDiff := -1 << 60
		for l := 0; l < n; l++ {
			count := make(map[int]int)
			for r := l; r < n; r++ {
				count[arr[r]]++
				for _, c := range count {
					diff := 2*c - (r - l + 1)
					if diff > maxDiff {
						maxDiff = diff
					}
				}
			}
		}
		var a, l, r int
		if _, err := fmt.Fscan(out, &a, &l, &r); err != nil {
			return fmt.Errorf("failed to read output for test %d: %v", tc, err)
		}
		if l < 1 || r < l || r > n {
			return fmt.Errorf("invalid segment on test %d", tc)
		}
		cnt := 0
		for i := l - 1; i < r; i++ {
			if arr[i] == a {
				cnt++
			}
		}
		candDiff := 2*cnt - (r - l + 1)
		if candDiff != maxDiff {
			return fmt.Errorf("expected diff %d, got %d", maxDiff, candDiff)
		}
	}
	// ensure no extra data? ignore
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := genTest()
		got, err := run(cand, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		if err := check(input, got); err != nil {
			fmt.Printf("wrong answer on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", string(input))
			fmt.Println("output:\n", got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
