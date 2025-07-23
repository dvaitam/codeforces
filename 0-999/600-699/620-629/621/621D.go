package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type val struct {
	cat int
	val float64
}

func calc(base, expPart float64) val {
	if math.Abs(base-1.0) < 1e-9 {
		return val{1, 0}
	} else if base > 1 {
		return val{2, expPart + math.Log(math.Log(base))}
	} else {
		return val{0, -(expPart + math.Log(-math.Log(base)))}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var x, y, z float64
	fmt.Fscan(reader, &x, &y, &z)
	names := []string{"x^y^z", "x^z^y", "(x^y)^z", "(x^z)^y", "y^x^z", "y^z^x", "(y^x)^z", "(y^z)^x", "z^x^y", "z^y^x", "(z^x)^y", "(z^y)^x"}
	vals := make([]val, 12)
	lnX := math.Log(x)
	lnY := math.Log(y)
	lnZ := math.Log(z)
	vals[0] = calc(x, z*lnY)
	vals[1] = calc(x, y*lnZ)
	vals[2] = calc(x, lnY+lnZ)
	vals[3] = calc(x, lnZ+lnY)
	vals[4] = calc(y, z*lnX)
	vals[5] = calc(y, x*lnZ)
	vals[6] = calc(y, lnX+lnZ)
	vals[7] = calc(y, lnZ+lnX)
	vals[8] = calc(z, y*lnX)
	vals[9] = calc(z, x*lnY)
	vals[10] = calc(z, lnX+lnY)
	vals[11] = calc(z, lnY+lnX)
	best := 0
	for i := 1; i < 12; i++ {
		if vals[i].cat > vals[best].cat || (vals[i].cat == vals[best].cat && vals[i].val > vals[best].val) {
			best = i
		}
	}
	fmt.Println(names[best])
}
