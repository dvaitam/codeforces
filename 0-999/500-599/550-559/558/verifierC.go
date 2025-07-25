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

const maxVal = 200000

type testCase struct {
	input    string
	expected string
}

func expectedAnswer(a []int) int {
	dist := make([]int, maxVal+1)
	cnt := make([]int, maxVal+1)
	for _, v := range a {
		visited := make(map[int]bool)
		val := v
		j := 0
		for val > 0 {
			x := val
			k := 0
			for x <= maxVal {
				if !visited[x] {
					dist[x] += j + k
					cnt[x]++
					visited[x] = true
				}
				x <<= 1
				k++
			}
			val >>= 1
			j++
		}
		if !visited[0] {
			dist[0] += j
			cnt[0]++
		}
	}
	ans := int(^uint(0) >> 1)
	n := len(a)
	for x := 0; x <= maxVal; x++ {
		if cnt[x] == n && dist[x] < ans {
			ans = dist[x]
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	arr := make([]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(1000) + 1
		sb.WriteString(fmt.Sprintf("%d ", arr[i]))
	}
	sb.WriteString("\n")
	expect := expectedAnswer(arr)
	return testCase{input: sb.String(), expected: fmt.Sprint(expect)}
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
	got := strings.TrimSpace(out.String())
	if got != tc.expected {
		return fmt.Errorf("expected %s got %s", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
