package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var V1, V2, V3, Vm int
	fmt.Fscan(in, &V1, &V2, &V3, &Vm)

	// Smallest car possible range based on son and Masha
	minS := Vm
	if minS < V3 {
		minS = V3
	}
	maxS := 2 * V3
	if maxS > 2*Vm {
		maxS = 2 * Vm
	}

	found := false
	for s := maxS; s >= minS && s > 0; s-- {
		// find middle
		maxM := 2 * V2
		if maxM > 2*s {
			maxM = 2 * s
		}
		for m := maxM; m > s; m-- {
			// m > s
			if V2 > m {
				continue
			}
			if 2*V2 < m {
				continue
			}
			// find largest car
			maxL := 2 * V1
			if maxL > 2*m {
				maxL = 2 * m
			}
			for l := maxL; l > m; l-- {
				if V1 > l {
					continue
				}
				if 2*V1 < l {
					continue
				}
				fmt.Fprintf(out, "%d %d %d\n", l, m, s)
				found = true
				break
			}
			if found {
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		fmt.Fprintln(out, -1)
	}
}
