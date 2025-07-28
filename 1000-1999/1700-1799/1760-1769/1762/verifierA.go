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

func expectedA(arr []int) int {
	sum := 0
	for _, v := range arr {
		sum += v
	}
	if sum%2 == 0 {
		return 0
	}
	minOps := 1 << 30
	for _, x := range arr {
		cnt := 0
		y := x
		if y%2 == 0 {
			for y%2 == 0 {
				y /= 2
				cnt++
			}
		} else {
			for y%2 == 1 {
				y /= 2
				cnt++
			}
		}
		if cnt < minOps {
			minOps = cnt
		}
	}
	return minOps
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func genTestA(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(1000000) + 1
	}
	exp := expectedA(arr)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), fmt.Sprintf("%d", exp)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := genTestA(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}
