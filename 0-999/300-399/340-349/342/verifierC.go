package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveC(r, h int) string {
	const diam = 7.0
	dz := r / 7
	if dz <= 0 {
		return "0\n"
	}
	cols := (2 * r) / 7
	if cols <= 0 {
		return "0\n"
	}
	R0 := float64(r) - diam/2
	total2D := 0
	rr := R0 * R0
	for i := 0; i < cols; i++ {
		x := -float64(r) + diam/2 + float64(i)*diam
		sq := x * x
		if sq > rr {
			continue
		}
		yMax := float64(h) + math.Sqrt(rr-sq)
		kMax := int(math.Floor((yMax-diam/2)/diam + 1e-9))
		if kMax >= 0 {
			total2D += kMax + 1
		}
	}
	result := int64(total2D) * int64(dz)
	return fmt.Sprintf("%d\n", result)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func generateCase(rng *rand.Rand) string {
	r := rng.Intn(1000) + 1
	h := rng.Intn(1000) + 1
	return fmt.Sprintf("%d %d\n", r, h)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		in := generateCase(rng)
		candOut, cErr := runBinary(bin, in)
		if cErr != nil {
			fmt.Printf("test %d: runtime error: %v\n", t+1, cErr)
			os.Exit(1)
		}
		r := 0
		h := 0
		fmt.Sscanf(strings.TrimSpace(in), "%d %d", &r, &h)
		expect := solveC(r, h)
		if strings.TrimSpace(candOut) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:%sexpected:%sactual:%s\n", t+1, in, expect, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
