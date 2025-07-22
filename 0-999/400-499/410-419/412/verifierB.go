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

func expectedB(arr []int, k int) int {
	b := append([]int(nil), arr...)
	sort.Ints(b)
	return b[len(b)-k]
}

func runCase(bin string, n, k int, arr []int) error {
	input := fmt.Sprintf("%d %d\n", n, k)
	for i, v := range arr {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprint(v)
	}
	input += "\n"
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(strings.TrimSpace(out.String())), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expect := expectedB(arr, k)
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func genCase(rng *rand.Rand) (int, int, []int) {
	n := rng.Intn(100) + 1
	k := rng.Intn(n) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(32768-16+1) + 16
	}
	return n, k, arr
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, k, arr := genCase(rng)
		if err := runCase(bin, n, k, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d %d\n%v\n", i+1, err, n, k, arr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
