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

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveCase(n, k int, arr []int) int {
	candies := 0
	total := 0
	for i := 0; i < n; i++ {
		candies += arr[i]
		give := 8
		if candies < 8 {
			give = candies
		}
		total += give
		candies -= give
		if total >= k {
			return i + 1
		}
	}
	return -1
}

func genCase(rng *rand.Rand) (int, int, []int) {
	n := rng.Intn(100) + 1
	k := rng.Intn(10000) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(100) + 1
	}
	return n, k, arr
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, k, arr := genCase(rng)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
		input := sb.String()
		want := solveCase(n, k, arr)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\ninput:\n%s", i+1, err, input)
			return
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Printf("case %d: non-integer output %q\n", i+1, out)
			return
		}
		if got != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %d\ngot: %d\n", i+1, input, want, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
