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

const maxDay = 1000005

var dp [maxDay]int

func solve(arr []int) int {
	used := make([]int, 0, len(arr)*2)
	best := 0
	for _, x := range arr {
		old := dp[x]
		valx := dp[x]
		if dp[x-1]+1 > valx {
			valx = dp[x-1] + 1
		}
		if valx != dp[x] {
			dp[x] = valx
			used = append(used, x)
		}
		if valx > best {
			best = valx
		}
		valx1 := dp[x+1]
		if old+1 > valx1 {
			valx1 = old + 1
		}
		if valx1 != dp[x+1] {
			dp[x+1] = valx1
			used = append(used, x+1)
		}
		if valx1 > best {
			best = valx1
		}
	}
	for _, idx := range used {
		dp[idx] = 0
	}
	return best
}

func runCase(bin string, a []int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(a)))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteString("\n")
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(&out, &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	want := solve(a)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func genCase(rng *rand.Rand) []int {
	n := rng.Intn(200) + 1
	a := make([]int, n)
	cur := rng.Intn(5) + 1
	for i := 0; i < n; i++ {
		cur += rng.Intn(10)
		a[i] = cur
	}
	return a
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
