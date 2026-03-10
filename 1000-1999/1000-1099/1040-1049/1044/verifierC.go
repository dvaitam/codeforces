package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"log"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func buildOracle() (string, error) {
	src := os.Getenv("REFERENCE_SOURCE_PATH")
	if src == "" {
		log.Fatal("REFERENCE_SOURCE_PATH environment variable is not set")
	}
	bin := filepath.Join(os.TempDir(), "oracle1044C.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return bin, nil
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

func cross(o, a, b [2]int) int {
	return (a[0]-o[0])*(b[1]-o[1]) - (a[1]-o[1])*(b[0]-o[0])
}

func genConvexPoints(r *rand.Rand, n int) [][2]int {
	for {
		angles := make([]float64, n)
		for i := range angles {
			angles[i] = r.Float64() * 2 * math.Pi
		}
		sort.Float64s(angles)
		pts := make([][2]int, n)
		radius := 10000.0
		for i, a := range angles {
			x := int(radius * math.Cos(a))
			y := int(radius * math.Sin(a))
			pts[i] = [2]int{x, y}
		}
		// Check all points are distinct
		valid := true
		for i := 0; i < n && valid; i++ {
			for j := i + 1; j < n && valid; j++ {
				if pts[i] == pts[j] {
					valid = false
				}
			}
		}
		if !valid {
			continue
		}
		// Check no three points are collinear
		for i := 0; i < n && valid; i++ {
			for j := i + 1; j < n && valid; j++ {
				for k := j + 1; k < n && valid; k++ {
					if cross(pts[i], pts[j], pts[k]) == 0 {
						valid = false
					}
				}
			}
		}
		if valid {
			return pts
		}
	}
}

func genCase(r *rand.Rand) string {
	n := r.Intn(6) + 3 // 3..8
	pts := genConvexPoints(r, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, p := range pts {
		fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
	for i := 0; i < 100; i++ {
		input := genCase(r)
		want, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		wantFields := strings.Fields(want)
		gotFields := strings.Fields(got)
		if len(wantFields) != len(gotFields) {
			fmt.Printf("test %d failed: expected %d values, got %d\ninput:\n%sexpected: %s\ngot: %s\n", i+1, len(wantFields), len(gotFields), input, want, got)
			os.Exit(1)
		}
		mismatch := false
		for j := range wantFields {
			if wantFields[j] != gotFields[j] {
				mismatch = true
				break
			}
		}
		if mismatch {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}
