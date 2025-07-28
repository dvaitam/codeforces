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

const maxVal = 1000000

var spf [maxVal + 1]int

func init() {
	for i := 2; i <= maxVal; i++ {
		if spf[i] == 0 {
			for j := i; j <= maxVal; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}
}

func squareFree(x int) int {
	res := 1
	for x > 1 {
		p := spf[x]
		if p == 0 {
			p = x
		}
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt++
		}
		if cnt%2 == 1 {
			res *= p
		}
	}
	return res
}

func solveCase(nums []int, queries []int64) []int {
	freq := make(map[int]int)
	for _, x := range nums {
		f := squareFree(x)
		freq[f]++
	}
	ans0 := 0
	merge := 0
	for v, c := range freq {
		if c > ans0 {
			ans0 = c
		}
		if v == 1 || c%2 == 0 {
			merge += c
		}
	}
	ans1 := ans0
	if merge > ans1 {
		ans1 = merge
	}
	res := make([]int, len(queries))
	for i, w := range queries {
		if w == 0 {
			res[i] = ans0
		} else {
			res[i] = ans1
		}
	}
	return res
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = rng.Intn(50) + 1
	}
	q := rng.Intn(5) + 1
	queries := make([]int64, q)
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 {
			queries[i] = 0
		} else {
			queries[i] = int64(rng.Intn(5) + 1)
		}
	}
	input := fmt.Sprintf("1\n%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", nums[i])
	}
	input += "\n"
	input += fmt.Sprintf("%d\n", q)
	for i := 0; i < q; i++ {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", queries[i])
	}
	input += "\n"
	ans := solveCase(nums, queries)
	var sb strings.Builder
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	return input, sb.String()
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
