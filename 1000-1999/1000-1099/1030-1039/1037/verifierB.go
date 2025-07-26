package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveCase(n int, s int, arr []int) string {
	sort.Ints(arr)
	mid := n / 2
	var ans int64
	if arr[mid] < s {
		for i := mid; i < n; i++ {
			if arr[i] < s {
				ans += int64(s - arr[i])
			}
		}
	} else {
		for i := 0; i <= mid; i++ {
			if arr[i] > s {
				ans += int64(arr[i] - s)
			}
		}
	}
	return fmt.Sprintf("%d", ans)
}

func genCase(r *rand.Rand) (string, string) {
	n := r.Intn(99)*2 + 1 // odd up to 199
	s := r.Intn(1_000_000_000) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = r.Intn(1_000_000_000)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, s))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	expect := solveCase(n, s, append([]int(nil), arr...))
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
