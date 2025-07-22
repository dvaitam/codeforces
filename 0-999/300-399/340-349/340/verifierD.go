package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func expectedAnswer(a []int) int {
	dp := make([]int, 0, len(a))
	for _, x := range a {
		i := sort.Search(len(dp), func(i int) bool { return dp[i] >= x })
		if i == len(dp) {
			dp = append(dp, x)
		} else {
			dp[i] = x
		}
	}
	return len(dp)
}

func generateCase(rng *rand.Rand) []int {
	n := rng.Intn(100) + 1
	perm := rand.Perm(n)
	for i := range perm {
		perm[i]++ // 1..n
	}
	return perm
}

func runCase(bin string, arr []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := strconv.Itoa(expectedAnswer(arr))
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		arr := generateCase(rng)
		if err := runCase(bin, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
