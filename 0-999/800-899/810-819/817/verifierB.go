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

func choose(n int64, k int) int64 {
	if k == 0 {
		return 1
	}
	switch k {
	case 1:
		return n
	case 2:
		return n * (n - 1) / 2
	case 3:
		return n * (n - 1) * (n - 2) / 6
	}
	return 0
}

func expected(arr []int) int64 {
	sort.Ints(arr)
	v1 := arr[0]
	cnt1 := 0
	for cnt1 < len(arr) && arr[cnt1] == v1 {
		cnt1++
	}
	if cnt1 >= 3 {
		return choose(int64(cnt1), 3)
	}
	v2 := arr[cnt1]
	cnt2 := 0
	idx := cnt1
	for idx < len(arr) && arr[idx] == v2 {
		cnt2++
		idx++
	}
	if cnt1 == 2 {
		return choose(2, 2) * int64(cnt2)
	}
	if cnt2 >= 2 {
		return choose(int64(cnt2), 2)
	}
	cnt3 := 0
	v3 := arr[idx]
	for idx < len(arr) && arr[idx] == v3 {
		cnt3++
		idx++
	}
	return int64(cnt3)
}

func runBin(bin, input string) (string, error) {
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

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 3
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(1000)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	exp := fmt.Sprintf("%d", expected(arr))
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		out, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\ninput:\n%s\n", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:%s\ngot:%s\n", i+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
