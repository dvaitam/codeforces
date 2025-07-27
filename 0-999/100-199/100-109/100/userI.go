package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var k int
	var x, y float64
	if _, err := fmt.Fscan(in, &k, &x, &y); err != nil {
		return
	}

	rad := float64(k) * math.Pi / 180.0
	c := math.Cos(rad)
	s := math.Sin(rad)

	nx := x*c - y*s
	ny := x*s + y*c

	fmt.Printf("%.10f %.10f\n", nx, ny)
}