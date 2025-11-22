package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var s, t string
		fmt.Fscan(in, &s)
		fmt.Fscan(in, &t)

		if len(s) != len(t) {
			fmt.Fprintln(out, -1)
			continue
		}

		n := len(s)
		zerosS := 0
		zerosT := 0
		onesT := 0
		for i := 0; i < n; i++ {
			if s[i] == '0' {
				zerosS++
			}
			if t[i] == '0' {
				zerosT++
			} else if t[i] == '1' {
				onesT++
			}
		}
		q := n - zerosT - onesT

		needZeroFromQ := zerosS - zerosT
		if needZeroFromQ < 0 || needZeroFromQ > q {
			fmt.Fprintln(out, -1)
			continue
		}

		fixedMismatch := 0 // positions where s=0 and t fixed to 1
		zeroQuestion := 0  // positions where s=0 and t='?'
		for i := 0; i < n; i++ {
			if s[i] == '0' {
				if t[i] == '1' {
					fixedMismatch++
				} else if t[i] == '?' {
					zeroQuestion++
				}
			}
		}

		// We need to turn "needZeroFromQ" of the '?' into '0'.
		// Prefer using positions where s has '0' to avoid mismatches.
		zerosUsedOnZeroPositions := needZeroFromQ
		if zerosUsedOnZeroPositions > zeroQuestion {
			zerosUsedOnZeroPositions = zeroQuestion
		}

		extraZeroToOne := zeroQuestion - zerosUsedOnZeroPositions // these become mismatches s=0, t=1
		ans := fixedMismatch + extraZeroToOne
		fmt.Fprintln(out, ans)
	}
}
