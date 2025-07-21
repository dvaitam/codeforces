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

type TestCase struct {
	x, y, z, k int64
}

func pieces(x, y, z, k int64) int64 {
	caps := []int64{x - 1, y - 1, z - 1}
	sort.Slice(caps, func(i, j int) bool { return caps[i] < caps[j] })
	sumcaps := caps[0] + caps[1] + caps[2]
	if k >= sumcaps {
		return x * y * z
	}
	T := k
	low, high := int64(0), caps[2]
	for low < high {
		mid := (low + high) / 2
		var s int64
		for i := 0; i < 3; i++ {
			if caps[i] < mid {
				s += caps[i]
			} else {
				s += mid
			}
		}
		if s >= T {
			high = mid
		} else {
			low = mid + 1
		}
	}
	L := low
	xcuts := make([]int64, 3)
	var sumX int64
	for i := 0; i < 3; i++ {
		if caps[i] < L {
			xcuts[i] = caps[i]
		} else {
			xcuts[i] = L
		}
		sumX += xcuts[i]
	}
	E := sumX - T
	for i := 2; i >= 0 && E > 0; i-- {
		if xcuts[i] == L {
			xcuts[i]--
			E--
		}
	}
	res := int64(1)
	for i := 0; i < 3; i++ {
		res *= (xcuts[i] + 1)
	}
	return res
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		x := rng.Int63n(10) + 1
		y := rng.Int63n(10) + 1
		z := rng.Int63n(10) + 1
		k := rng.Int63n(30)
		input := fmt.Sprintf("%d %d %d %d\n", x, y, z, k)
		expected := fmt.Sprintf("%d", pieces(x, y, z, k))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
