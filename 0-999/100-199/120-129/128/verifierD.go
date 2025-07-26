package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func runProgram(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func solve(arr []int) string {
	n := len(arr)
	if n < 3 {
		return "NO"
	}
	sort.Ints(arr)
	uniq := []int{arr[0]}
	cnt := []int{1}
	last := arr[0]
	for i := 1; i < n; i++ {
		if arr[i] == last {
			cnt[len(cnt)-1]++
		} else {
			if arr[i] != last+1 {
				return "NO"
			}
			last = arr[i]
			uniq = append(uniq, last)
			cnt = append(cnt, 1)
		}
	}
	ePrev := 0
	for _, c := range cnt {
		eCur := 2*c - ePrev
		if eCur < 0 {
			return "NO"
		}
		ePrev = eCur
	}
	if ePrev != 0 {
		return "NO"
	}
	return "YES"
}

func randomTest(rng *rand.Rand) []int {
	m := rng.Intn(5) + 3
	// create consecutive numbers with possible duplicates
	start := rng.Intn(10)
	arr := make([]int, m)
	for i := 0; i < m; i++ {
		arr[i] = start + rng.Intn(3)
	}
	return arr
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		arr := randomTest(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		expected := solve(append([]int(nil), arr...))
		out, err := runProgram(bin, sb.String())
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			return
		}
		if strings.TrimSpace(out) != expected {
			fmt.Printf("Test %d failed: expected %s got %s\n", i+1, expected, strings.TrimSpace(out))
			return
		}
	}
	fmt.Println("All tests passed")
}
