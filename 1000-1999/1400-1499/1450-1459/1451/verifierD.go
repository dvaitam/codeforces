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

func expected(d, k int64) string {
	D2 := d * d
	k2 := k * k
	x := int64(math.Sqrt(float64(D2) / float64(2*k2)))
	for (x+1)*(x+1)*2*k2 <= D2 {
		x++
	}
	for x >= 0 && x*x*2*k2 > D2 {
		x--
	}
	rem := D2 - x*x*k2
	y := int64(math.Sqrt(float64(rem) / float64(k2)))
	for (y+1)*(y+1)*k2 <= rem {
		y++
	}
	for y >= 0 && y*y*k2 > rem {
		y--
	}
	if y > x {
		return "Ashish"
	}
	return "Utkarsh"
}

func runCase(bin string, d, k int64) error {
	input := fmt.Sprintf("1\n%d %d\n", d, k)
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
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	want := expected(d, k)
	if strings.ToLower(got) != strings.ToLower(want) {
		return fmt.Errorf("expected %s got %s", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		d := rng.Int63n(1000) + 1
		k := rng.Int63n(d) + 1
		if err := runCase(bin, d, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
