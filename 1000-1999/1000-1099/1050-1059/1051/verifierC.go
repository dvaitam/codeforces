package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func possible(nums []int) bool {
	cnt := make([]int, 110)
	for _, v := range nums {
		cnt[v]++
	}
	calc := make([]int, 4)
	for _, c := range cnt {
		if c <= 2 {
			calc[c]++
		} else {
			calc[3]++
		}
	}
	return calc[1]%2 == 0 || calc[3] > 0
}

func isValid(nums []int, assign string) bool {
	if len(assign) != len(nums) {
		return false
	}
	cnt := make(map[int]int)
	for _, v := range nums {
		cnt[v]++
	}
	aNice, bNice := 0, 0
	for i, ch := range assign {
		if ch != 'A' && ch != 'B' {
			return false
		}
		if cnt[nums[i]] == 1 {
			if ch == 'A' {
				aNice++
			} else {
				bNice++
			}
		}
	}
	return aNice == bNice
}

func genCase(rng *rand.Rand) []int {
	n := rng.Intn(8) + 2
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = rng.Intn(10)
	}
	return nums
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		nums := genCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(nums)))
		for j, v := range nums {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		out, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, sb.String())
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if lines[0] == "NO" {
			if possible(nums) {
				fmt.Fprintf(os.Stderr, "case %d failed: expected YES got NO\n", i+1)
				os.Exit(1)
			}
			continue
		}
		if lines[0] != "YES" {
			fmt.Fprintf(os.Stderr, "case %d failed: expected YES/NO got %s\n", i+1, lines[0])
			os.Exit(1)
		}
		if len(lines) < 2 {
			fmt.Fprintf(os.Stderr, "case %d failed: missing assignment\n", i+1)
			os.Exit(1)
		}
		assign := strings.TrimSpace(lines[1])
		if !possible(nums) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected NO but got YES\n", i+1)
			os.Exit(1)
		}
		if !isValid(nums, assign) {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid assignment %s\n", i+1, assign)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
