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

func solveCase(arr []int64) []int64 {
	pos := make(map[int64]int)
	removed := make([]bool, len(arr))
	for i := 0; i < len(arr); i++ {
		v := arr[i]
		for {
			if j, ok := pos[v]; ok {
				delete(pos, v)
				removed[j] = true
				v *= 2
			} else {
				break
			}
		}
		arr[i] = v
		pos[v] = i
	}
	var res []int64
	for i := 0; i < len(arr); i++ {
		if !removed[i] {
			res = append(res, arr[i])
		}
	}
	return res
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	arr := make([]int64, n)
	input := fmt.Sprintf("%d\n", n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Intn(10) + 1)
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", arr[i])
	}
	input += "\n"
	outArr := solveCase(append([]int64(nil), arr...))
	res := fmt.Sprintf("%d\n", len(outArr))
	for i, v := range outArr {
		if i > 0 {
			res += " "
		}
		res += fmt.Sprintf("%d", v)
	}
	if len(outArr) > 0 {
		res += "\n"
	}
	return input, res
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
