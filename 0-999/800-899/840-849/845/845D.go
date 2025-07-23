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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	currentSpeed := 0
	var speedLimits []int
	noOvertake := 0
	ans := 0

	for i := 0; i < n; i++ {
		var t int
		fmt.Fscan(in, &t)
		switch t {
		case 1:
			var s int
			fmt.Fscan(in, &s)
			currentSpeed = s
			for len(speedLimits) > 0 && currentSpeed > speedLimits[len(speedLimits)-1] {
				ans++
				speedLimits = speedLimits[:len(speedLimits)-1]
			}
		case 2:
			if noOvertake > 0 {
				ans += noOvertake
				noOvertake = 0
			}
		case 3:
			var limit int
			fmt.Fscan(in, &limit)
			if currentSpeed > limit {
				ans++
			} else {
				speedLimits = append(speedLimits, limit)
			}
		case 4:
			noOvertake = 0
		case 5:
			speedLimits = speedLimits[:0]
		case 6:
			noOvertake++
		}
	}

	fmt.Fprintln(out, ans)
}
