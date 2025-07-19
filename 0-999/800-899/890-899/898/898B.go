package main

import (
	"bufio"
	"fmt"
	"os"
)

// abs returns the absolute value of x.
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int64) int64 {
	if b == 0 {
		return abs(a)
	}
	return gcd(b, a%b)
}

// extended computes x, y such that a*x + b*y = gcd(a, b).
func extended(a, b int64) (int64, int64) {
	if b == 0 {
		return 1, 0
	}
	x1, y1 := extended(b, a%b)
	// x = y1, y = x1 - (a/b)*y1
	return y1, x1 - (a/b)*y1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var c, a, b int64
	if _, err := fmt.Fscan(reader, &c, &a, &b); err != nil {
		return
	}
	// Check if equation a*x + b*y = c has integer solutions
	g := gcd(a, b)
	if c%g != 0 {
		fmt.Fprintln(writer, "NO")
		return
	}
	swapped := false
	// Ensure a >= b for adjustment logic
	if a < b {
		a, b = b, a
		swapped = true
	}
	// Get initial solution
	s, t := extended(a, b)
	x0 := (c / g) * s
	y0 := (c / g) * t
	x, y := x0, y0
	// Adjust to non-negative solution
	if x < 0 {
		k := int64(1)
		stepX := b / g
		stepY := a / g
		for x < 0 {
			x = x0 + k*stepX
			y = y0 - k*stepY
			k++
		}
	} else if y < 0 {
		k := int64(1)
		stepX := b / g
		stepY := a / g
		for y < 0 {
			x = x0 - k*stepX
			y = y0 + k*stepY
			k++
		}
	}
	// If no non-negative solution
	if x < 0 || y < 0 {
		fmt.Fprintln(writer, "NO")
		return
	}
	fmt.Fprintln(writer, "YES")
	if !swapped {
		fmt.Fprintf(writer, "%d %d\n", x, y)
	} else {
		fmt.Fprintf(writer, "%d %d\n", y, x)
	}
}
