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

func generateCase(rng *rand.Rand) (string, float64) {
	a := float64(rng.Intn(5) + 1)
	v := float64(rng.Intn(10) + 1)
	l := float64(rng.Intn(9) + 2)
	d := float64(rng.Intn(int(l-1)) + 1)
	w := float64(rng.Intn(10) + 1)
	input := fmt.Sprintf("%d %d %d %d %d\n", int(a), int(v), int(l), int(d), int(w))
	expected := solve(a, v, l, d, w)
	return input, expected
}

func solve(a, v, l, d, w float64) float64 {
	t := 0.0
	if w > v {
		w = v
	}
	l3 := w * w / (2 * a)
	var v1 float64
	if l3 <= d {
		t += w / a
		l1 := d - l3
		l3 = math.Sqrt(a*l1 + w*w)
		if l3 > v {
			l3 = v
		}
		t += 2 * (l3 - w) / a
		t += (l1 - (l3*l3-w*w)/a) / l3
		v1 = w
	} else {
		t += math.Sqrt(2 * d / a)
		v1 = a * t
	}
	l = l - d
	l2 := (v*v - v1*v1) / (2 * a)
	if l2 <= l {
		t += (v - v1) / a
		t += (l - l2) / v
	} else {
		t += (math.Sqrt(2*a*l+v1*v1) - v1) / a
	}
	return t
}

func runCase(bin, input string, expected float64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got float64
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &got); err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	if math.Abs(got-expected) > 1e-4 {
		return fmt.Errorf("expected %.5f got %.5f", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
