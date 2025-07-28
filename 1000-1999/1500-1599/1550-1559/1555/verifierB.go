package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleB")
	cmd := exec.Command("g++", "-std=c++17", "-O2", "-o", oracle, "solB.cpp")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func expected(W, H, x1, y1, x2, y2, w, h int64) string {
	const inf int64 = 1 << 60
	ans := inf
	if x2-x1+w <= W {
		if v := w - x1; v > 0 && v < ans {
			ans = v
		} else if v := int64(0); v < ans {
			if w-x1 <= 0 {
				ans = 0
			}
		}
		if v := x2 - (W - w); v > 0 && v < ans {
			ans = v
		} else if v := int64(0); v < ans {
			if x2-(W-w) <= 0 {
				ans = 0
			}
		}
	}
	if y2-y1+h <= H {
		if v := h - y1; v > 0 && v < ans {
			ans = v
		} else if v := int64(0); v < ans {
			if h-y1 <= 0 {
				ans = 0
			}
		}
		if v := y2 - (H - h); v > 0 && v < ans {
			ans = v
		} else if v := int64(0); v < ans {
			if y2-(H-h) <= 0 {
				ans = 0
			}
		}
	}
	if ans == inf {
		return "-1"
	}
	return fmt.Sprint(ans)
}

func generate(rng *rand.Rand) (int64, int64, int64, int64, int64, int64, int64, int64) {
	W := int64(rng.Intn(1000) + 2)
	H := int64(rng.Intn(1000) + 2)
	x1 := int64(rng.Intn(int(W - 1)))
	x2 := x1 + int64(rng.Intn(int(W-x1))) + 1
	y1 := int64(rng.Intn(int(H - 1)))
	y2 := y1 + int64(rng.Intn(int(H-y1))) + 1
	w := int64(rng.Intn(int(W)) + 1)
	h := int64(rng.Intn(int(H)) + 1)
	return W, H, x1, y1, x2, y2, w, h
}

func runCase(bin, oracle string, params []int64) error {
	W, H, x1, y1, x2, y2, w, h := params[0], params[1], params[2], params[3], params[4], params[5], params[6], params[7]
	input := fmt.Sprintf("1\n%d %d\n%d %d %d %d\n%d %d\n", W, H, x1, y1, x2, y2, w, h)
	cmdO := exec.Command(oracle)
	cmdO.Stdin = strings.NewReader(input)
	var outO bytes.Buffer
	cmdO.Stdout = &outO
	if err := cmdO.Run(); err != nil {
		return fmt.Errorf("oracle run error: %v", err)
	}
	expected := strings.TrimSpace(outO.String())
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		p := make([]int64, 8)
		W, H, x1, y1, x2, y2, w, h := generate(rng)
		p[0], p[1], p[2], p[3], p[4], p[5], p[6], p[7] = W, H, x1, y1, x2, y2, w, h
		if err := runCase(bin, oracle, p); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
