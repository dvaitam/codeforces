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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func indexInSequence(n int64, k int) int64 {
	used := make(map[int64]bool)
	var p int64 = 1
	var idx int64
	for {
		nums := make([]int64, 0, k)
		for len(nums) < k {
			if !used[p] {
				nums = append(nums, p)
				used[p] = true
			}
			p++
		}
		for _, v := range nums {
			idx++
			if v == n {
				return idx
			}
		}
		var sum int64
		for _, v := range nums {
			sum += v
		}
		idx++
		if sum == n {
			return idx
		}
		used[sum] = true
		if idx > n*2 {
			break
		}
	}
	return -1
}

func genCase(rng *rand.Rand) (int64, int) {
	n := int64(rng.Intn(200) + 1)
	k := rng.Intn(4) + 1
	return n, k
}

func formatInput(cases [][2]int64) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, c := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", c[0], c[1]))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases [][2]int64
	cases = append(cases, [2]int64{1, 1})
	for len(cases) < 100 {
		n, k := genCase(rng)
		cases = append(cases, [2]int64{n, int64(k)})
	}

	input := formatInput(cases)
	var expOut strings.Builder
	for _, c := range cases {
		expOut.WriteString(fmt.Sprintf("%d\n", indexInSequence(c[0], int(c[1]))))
	}
	expected := strings.TrimSpace(expOut.String())

	out, err := run(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "verification failed: %v\n", err)
		os.Exit(1)
	}
	if strings.TrimSpace(out) != expected {
		fmt.Fprintf(os.Stderr, "verification failed: expected:\n%s\n got:\n%s\n", expected, out)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
