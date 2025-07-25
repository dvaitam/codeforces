package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func solveC(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var C, Hr, Hb, Wr, Wb int64
	if _, err := fmt.Fscan(in, &C, &Hr, &Hb, &Wr, &Wb); err != nil {
		return ""
	}
	w1, h1 := Wr, Hr
	w2, h2 := Wb, Hb
	if w1 > w2 {
		w1, w2 = w2, w1
		h1, h2 = h2, h1
	}
	limit := int64(math.Sqrt(float64(C))) + 1
	var maxJoy int64
	maxX := min64(C/w1, limit)
	for x := int64(0); x <= maxX; x++ {
		rem := C - x*w1
		y := rem / w2
		joy := x*h1 + y*h2
		if joy > maxJoy {
			maxJoy = joy
		}
	}
	maxY := min64(C/w2, limit)
	for y := int64(0); y <= maxY; y++ {
		rem := C - y*w2
		x := rem / w1
		joy := x*h1 + y*h2
		if joy > maxJoy {
			maxJoy = joy
		}
	}
	return fmt.Sprint(maxJoy)
}

func genTestC(rng *rand.Rand) string {
	C := rng.Int63n(1000) + 1
	Hr := rng.Int63n(100) + 1
	Hb := rng.Int63n(100) + 1
	Wr := rng.Int63n(100) + 1
	Wb := rng.Int63n(100) + 1
	return fmt.Sprintf("%d %d %d %d %d\n", C, Hr, Hb, Wr, Wb)
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 1; i <= 100; i++ {
		in := genTestC(rng)
		expect := solveC(in)
		got, err := run(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
