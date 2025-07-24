package main

import (
	"bufio"
	"fmt"
	"os"
)

func extendedGCD(a, b int64) (int64, int64, int64) {
	if b == 0 {
		return a, 1, 0
	}
	g, x1, y1 := extendedGCD(b, a%b)
	return g, y1, x1 - (a/b)*y1
}

func ceilDiv(a, b int64) int64 {
	if b < 0 {
		a, b = -a, -b
	}
	if a >= 0 {
		return (a + b - 1) / b
	}
	return a / b
}

func floorDiv(a, b int64) int64 {
	if b < 0 {
		a, b = -a, -b
	}
	if a >= 0 {
		return a / b
	}
	return -((-a + b - 1) / b)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var a, h, w int64
	if _, err := fmt.Fscan(reader, &a, &h, &w); err != nil {
		return
	}

	if h < a || w < a {
		fmt.Fprintln(writer, -1)
		return
	}

	A := h + a
	B := w + a
	C := w - h

	g, x0, y0 := extendedGCD(A, B)
	if C%g != 0 {
		fmt.Fprintln(writer, -1)
		return
	}

	factor := C / g
	m0 := x0 * factor
	n0 := -y0 * factor

	Bdiv := B / g
	Adiv := A / g

	lowN := ceilDiv(1-n0, Adiv)
	highN := floorDiv(h/a-n0, Adiv)
	lowM := ceilDiv(1-m0, Bdiv)
	highM := floorDiv(w/a-m0, Bdiv)

	low := lowN
	if lowM > low {
		low = lowM
	}
	high := highN
	if highM < high {
		high = highM
	}

	if low > high {
		fmt.Fprintln(writer, -1)
		return
	}

	t := high
	n := n0 + Adiv*t
	_ = m0 + Bdiv*t // m value not needed explicitly
	x := float64(h-n*a) / float64(n+1)
	fmt.Fprintf(writer, "%.10f\n", x)
}
