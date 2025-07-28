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

func expectedD(arr []int) int {
	pos := 1
	for i, v := range arr {
		if v == 0 {
			pos = i + 1
		}
	}
	return pos
}

func genTestD(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		if rng.Intn(5) == 0 {
			arr[i] = 0
		} else {
			arr[i] = rng.Intn(100) + 1
		}
	}
	exp := expectedD(arr)
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
	return sb.String(), fmt.Sprintf("%d %d", exp, exp)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input, expect := genTestD(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n%s", t+1, err, got)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", t+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}
