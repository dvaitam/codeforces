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

	var q int
	fmt.Fscan(in, &q)

	for k := 0; k < q; k++ {
		var b, w int
		fmt.Fscan(in, &b, &w)

		if b > 3*w+1 || w > 3*b+1 {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "YES")
			if b >= w {
				// We need more black cells than white cells (or equal).
				// Construct a vertical spine at x=2 starting at y=2.
				// (2, y) is White if y is even. (2, y+1) is Black.
				// We place w pairs of (White, Black).
				// We attach extra Black cells to the White cells of the spine.
				rem := b - w
				for j := 0; j < w; j++ {
					y := 2 + 2*j
					fmt.Fprintln(out, 2, y)   // White
					fmt.Fprintln(out, 2, y+1) // Black

					// Add extra Black neighbors to the White cell (2, y)
					if rem > 0 {
						fmt.Fprintln(out, 1, y) // Left neighbor (Black)
						rem--
					}
					if rem > 0 {
						fmt.Fprintln(out, 3, y) // Right neighbor (Black)
						rem--
					}
				}
				// If we still need one more Black cell, attach to top of first White cell
				if rem > 0 {
					fmt.Fprintln(out, 2, 1)
				}
			} else {
				// We need more white cells than black cells.
				// Construct a vertical spine at x=2 starting at y=3.
				// (2, y) is Black if y is odd. (2, y+1) is White.
				// We place b pairs of (Black, White).
				// We attach extra White cells to the Black cells of the spine.
				rem := w - b
				for j := 0; j < b; j++ {
					y := 3 + 2*j
					fmt.Fprintln(out, 2, y)   // Black
					fmt.Fprintln(out, 2, y+1) // White

					// Add extra White neighbors to the Black cell (2, y)
					if rem > 0 {
						fmt.Fprintln(out, 1, y) // Left neighbor (White)
						rem--
					}
					if rem > 0 {
						fmt.Fprintln(out, 3, y) // Right neighbor (White)
						rem--
					}
				}
				// If we still need one more White cell, attach to top of first Black cell
				if rem > 0 {
					fmt.Fprintln(out, 2, 2)
				}
			}
		}
	}
}