package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	h   int
	cnt int64
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	h := make([]int, n)
	maxVal := -1
	maxIdx := -1
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &h[i])
		if h[i] > maxVal {
			maxVal = h[i]
			maxIdx = i
		}
	}

	// Rotate array so max element is at index 0
	rotatedH := make([]int, n)
	for i := 0; i < n; i++ {
		rotatedH[i] = h[(maxIdx+i)%n]
	}

	// Monotonic stack
	stack := make([]pair, 0, n)
	stack = append(stack, pair{rotatedH[0], 1})
	
	var res int64

	// Iterate from 1 to n-1
	for i := 1; i < n; i++ {
		cur := rotatedH[i]
		// Pop elements strictly smaller than current
		for len(stack) > 0 && stack[len(stack)-1].h < cur {
			res += stack[len(stack)-1].cnt
			stack = stack[:len(stack)-1]
		}
		
		// Handle equal elements
		if len(stack) > 0 && stack[len(stack)-1].h == cur {
			top := stack[len(stack)-1]
			res += top.cnt
			stack = stack[:len(stack)-1]
			
			// If there is a larger element below the equal block, current sees it too
			if len(stack) > 0 {
				res++
			}
			
			// Push merged block
			stack = append(stack, pair{cur, top.cnt + 1})
		} else {
			// Strictly smaller (already popped) or stack empty (not possible since max at 0)
			// Actually stack is never empty because rotatedH[0] is maxVal and others <= maxVal
			// If stack not empty, current sees the larger element at top
			if len(stack) > 0 {
				res++
			}
			stack = append(stack, pair{cur, 1})
		}
	}

	// Post-processing for wrap-around visibility
	// All elements remaining in the stack (except the bottom-most max block) 
	// can see the max element(s) via the "right" arc (wrap around).
	
	// Logic:
	// 1. Groups at index i >= 2 in the stack are separated from Max by stack[1] (which is > stack[i]).
	//    So they didn't see Max from the left. They definitely see Max from the right.
	//    Add stack[i].cnt.
	// 2. Group at index 1 (stack[1]) saw Max from the left.
	//    Do they see Max from the right? Yes.
	//    Is it a distinct pair?
	//    - If stack[0].cnt == 1 (unique max), then Left Max is the SAME as Right Max. Pair already counted.
	//    - If stack[0].cnt > 1 (multiple maxes), Left Max != Right Max (different indices). New pair.
	
	if len(stack) > 0 {
		// Check the max block
		firstGroup := stack[0]
		
		// Iterate remaining groups
		for i := 1; i < len(stack); i++ {
			if i == 1 && firstGroup.cnt == 1 {
				continue
			}
			res += stack[i].cnt
		}
	}

	fmt.Println(res)
}

