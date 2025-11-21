package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const testCount = 160

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "763B.go")
	tmp, err := os.CreateTemp("", "oracle763B")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()
	os.Remove(path)
	cmd := exec.Command("go", "build", "-o", path, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	return path, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

type rect struct {
	x1, y1, x2, y2 int
}

func genCase(r *rand.Rand) (string, []rect) {
	n := 1 + r.Intn(2000)
	rects := make([]rect, n)
	for i := 0; i < n; i++ {
		baseX := r.Intn(200) * 4
		baseY := r.Intn(200) * 4
		w := r.Intn(5)*2 + 1
		h := r.Intn(5)*2 + 1
		rects[i] = rect{baseX, baseY, baseX + w, baseY + h}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, rec := range rects {
		fmt.Fprintf(&sb, "%d %d %d %d\n", rec.x1, rec.y1, rec.x2, rec.y2)
	}
	return sb.String(), rects
}

func parseOutput(out string, n int) ([]int, error) {
	lines := strings.Fields(out)
	if len(lines) < 1 {
		return nil, fmt.Errorf("empty output")
	}
	if strings.ToUpper(lines[0]) == "NO" {
		return nil, fmt.Errorf("unexpected NO")
	}
	if strings.ToUpper(lines[0]) != "YES" {
		return nil, fmt.Errorf("expected YES or NO, got %s", lines[0])
	}
	if len(lines) != n+1 {
		return nil, fmt.Errorf("expected %d colors, got %d", n, len(lines)-1)
	}
	ans := make([]int, n)
	for i := 0; i < n; i++ {
		val, err := strconv.Atoi(lines[i+1])
		if err != nil {
			return nil, err
		}
		if val < 1 || val > 4 {
			return nil, fmt.Errorf("color %d out of range", val)
		}
		ans[i] = val
	}
	return ans, nil
}

func validateColors(rects []rect, colors []int) error {
	n := len(rects)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if touch(rects[i], rects[j]) {
				if colors[i] == colors[j] {
					return fmt.Errorf("rectangles %d and %d touch but share color %d", i+1, j+1, colors[i])
				}
			}
		}
	}
	return nil
}

func touch(a, b rect) bool {
	// check overlapping projection on one axis and exact touch on other
	if a.x1 == b.x1 && a.y1 == b.y1 && a.x2 == b.x2 && a.y2 == b.y2 {
		return true
	}
	// vertical sides touching
	if max(a.y1, b.y1) < min(a.y2, b.y2) {
		if a.x2 == b.x1 || b.x2 == a.x1 {
			return true
		}
	}
	// horizontal sides touching
	if max(a.x1, b.x1) < min(a.x2, b.x2) {
		if a.y2 == b.y1 || b.y2 == a.y1 {
			return true
		}
	}
	return false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	r := rand.New(rand.NewSource(1))
	for t := 0; t < testCount; t++ {
		input, rects := genCase(r)
		expectStr, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		if _, err := parseOutput(expectStr, len(rects)); err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotStr, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		colors, err := parseOutput(gotStr, len(rects))
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", t+1, input, err)
			os.Exit(1)
		}
		if err := validateColors(rects, colors); err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", t+1, input, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
