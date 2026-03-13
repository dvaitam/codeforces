package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

// computeAnswer computes the correct answer for CF 1156A.
// Shapes: 1=circle, 2=triangle, 3=square, inscribed one inside another.
// Adjacent triangle-square or square-triangle = infinite intersection points.
// Adjacent circle-triangle or triangle-circle = 3 intersection points.
// Adjacent circle-square or square-circle = 4 intersection points.
func computeAnswer(a []int) (bool, int) {
	ans := 0
	for i := 1; i < len(a); i++ {
		x, y := a[i-1], a[i]
		if (x == 2 && y == 3) || (x == 3 && y == 2) {
			return true, 0 // infinite
		}
		if (x == 1 && y == 2) || (x == 2 && y == 1) {
			ans += 3
		}
		if (x == 1 && y == 3) || (x == 3 && y == 1) {
			ans += 4
		}
	}
	// Deduplication: when a circle (1) is between a triangle (2) and square (3),
	// one intersection point is shared between the two adjacent pairs.
	for i := 0; i+2 < len(a); i++ {
		if a[i] == 3 && a[i+1] == 1 && a[i+2] == 2 {
			ans--
		}
	}
	return false, ans
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

func genCase(r *rand.Rand) (string, []int) {
	n := r.Intn(99-2+1) + 2 // 2..99
	arr := make([]int, n)
	arr[0] = r.Intn(3) + 1
	for i := 1; i < n; i++ {
		for {
			x := r.Intn(3) + 1
			if x != arr[i-1] {
				arr[i] = x
				break
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String(), arr
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		input, arr := genCase(r)
		infinite, count := computeAnswer(arr)

		var want string
		if infinite {
			want = "Infinite"
		} else {
			want = fmt.Sprintf("Finite\n%d", count)
		}

		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}

		// Normalize: case-insensitive first word, flexible whitespace
		wantNorm := strings.TrimSpace(want)
		gotNorm := strings.TrimSpace(got)

		if !strings.EqualFold(wantNorm, gotNorm) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, wantNorm, gotNorm)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}
