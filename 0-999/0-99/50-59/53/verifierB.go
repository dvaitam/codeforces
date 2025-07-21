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

func bestRect(h, w int64) (int64, int64) {
	var bestH, bestW, bestArea int64
	for x := 0; x < 62; x++ {
		h0 := int64(1) << x
		if h0 > h {
			break
		}
		wMin := (4*h0 + 5 - 1) / 5
		wMax := (5 * h0) / 4
		if wMin < 1 {
			wMin = 1
		}
		if wMax > w {
			wMax = w
		}
		if wMin <= wMax {
			w0 := wMax
			area := h0 * w0
			if area > bestArea || (area == bestArea && h0 > bestH) {
				bestArea = area
				bestH = h0
				bestW = w0
			}
		}
	}
	for x := 0; x < 62; x++ {
		w0 := int64(1) << x
		if w0 > w {
			break
		}
		hMin := (4*w0 + 5 - 1) / 5
		hMax := (5 * w0) / 4
		if hMin < 1 {
			hMin = 1
		}
		if hMax > h {
			hMax = h
		}
		if hMin <= hMax {
			h0 := hMax
			area := h0 * w0
			if area > bestArea || (area == bestArea && h0 > bestH) {
				bestArea = area
				bestH = h0
				bestW = w0
			}
		}
	}
	return bestH, bestW
}

func generateCase(rng *rand.Rand) (string, string) {
	h := rng.Int63n(1_000_000_000) + 1
	w := rng.Int63n(1_000_000_000) + 1
	h0, w0 := bestRect(h, w)
	in := fmt.Sprintf("%d %d\n", h, w)
	out := fmt.Sprintf("%d %d\n", h0, w0)
	return in, out
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if outStr != exp {
		return fmt.Errorf("expected %q got %q", exp, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
