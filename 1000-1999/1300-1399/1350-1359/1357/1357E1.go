package main

import (
	"bufio"
	"fmt"
	"math"
	"math/cmplx"
	"os"
)

// qft performs the quantum Fourier transform on the state vector.
// state represents amplitudes of basis states in little-endian order.
func qft(state []complex128) []complex128 {
	n := len(state)
	res := make([]complex128, n)
	norm := math.Sqrt(float64(n))
	for k := 0; k < n; k++ {
		var sum complex128
		for j := 0; j < n; j++ {
			angle := 2 * math.Pi * float64(j*k) / float64(n)
			sum += state[j] * cmplx.Exp(complex(0, angle))
		}
		res[k] = sum / complex(norm, 0)
	}
	return res
}

// powerQFT applies the QFT operator P times using the property that QFT^4 = I.
func powerQFT(state []complex128, P int) []complex128 {
	P %= 4
	for i := 0; i < P; i++ {
		state = qft(state)
	}
	return state
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var P, n int
	if _, err := fmt.Fscan(in, &P, &n); err != nil {
		return
	}
	size := 1 << n
	state := make([]complex128, size)
	for i := 0; i < size; i++ {
		var re, im float64
		fmt.Fscan(in, &re, &im)
		state[i] = complex(re, im)
	}
	state = powerQFT(state, P)

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i := 0; i < size; i++ {
		fmt.Fprintf(out, "%.10f %.10f\n", real(state[i]), imag(state[i]))
	}
}
