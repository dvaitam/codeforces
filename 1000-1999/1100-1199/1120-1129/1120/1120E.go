package main

import (
	"bufio"
	"fmt"
	"os"
)

type state struct {
	carry int
	diff  int
}

func digitSum(x int) int {
	sum := 0
	for x > 0 {
		sum += x % 10
		x /= 10
	}
	return sum
}

func isSolvableBase(a int) bool {
	tmp := a
	for tmp%2 == 0 {
		tmp /= 2
	}
	for tmp%5 == 0 {
		tmp /= 5
	}
	if tmp > 1 {
		return true
	}
	return a == 2 || a == 4 || a == 8
}

func encodeKey(c, d, span int) int64 {
	return int64(c)*int64(span) + int64(d)
}

func bfsSolve(a int, limit int, sums []int) (string, bool) {
	span := 2*limit + 1
	encode := func(c, d int) int64 {
		return encodeKey(c, d+limit, span)
	}

	stateCarry := []int{0}
	stateDiff := []int{0}
	parent := []int{-1}
	digit := []byte{0}

	queue := make([]int, 0, 1<<14)
	queue = append(queue, 0)

	visited := map[int64]int{encode(0, 0): 0}

	for len(queue) > 0 {
		idx := queue[0]
		queue = queue[1:]

		c := stateCarry[idx]
		d := stateDiff[idx]

		for x := 0; x <= 9; x++ {
			t := a*x + c
			y := t % 10
			c2 := t / 10
			delta := x - a*y - a*(sums[c2]-sums[c])
			d2 := d + delta
			if d2 < -limit || d2 > limit {
				continue
			}
			key := encode(c2, d2)
			if _, ok := visited[key]; ok {
				continue
			}
			visited[key] = len(stateCarry)
			stateCarry = append(stateCarry, c2)
			stateDiff = append(stateDiff, d2)
			parent = append(parent, idx)
			digit = append(digit, byte(x))
			queue = append(queue, len(stateCarry)-1)

			if d2 == 0 && len(stateCarry)-1 != 0 {
				digits := make([]byte, 0)
				cur := len(stateCarry) - 1
				nonZero := false
				for cur != 0 {
					digits = append(digits, digit[cur])
					if digit[cur] != 0 {
						nonZero = true
					}
					cur = parent[cur]
					if len(digits) > 500000 {
						break
					}
				}
				if !nonZero || len(digits) == 0 || len(digits) > 500000 {
					continue
				}
				for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
					digits[i], digits[j] = digits[j], digits[i]
				}
				res := make([]byte, len(digits))
				for i, v := range digits {
					res[i] = byte('0' + v)
				}
				trim := 0
				for trim < len(res)-1 && res[trim] == '0' {
					trim++
				}
				return string(res[trim:]), true
			}
		}
	}
	return "", false
}

func findNumber(a int) (string, bool) {
	sums := make([]int, a+1)
	for i := 0; i <= a; i++ {
		sums[i] = digitSum(i)
	}

	for limit := 9 * a; limit <= 500000; limit *= 2 {
		if limit == 0 {
			limit = 1
		}
		if res, ok := bfsSolve(a, limit, sums); ok {
			return res, true
		}
	}
	return "", false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var a int
	if _, err := fmt.Fscan(in, &a); err != nil {
		return
	}

	if !isSolvableBase(a) {
		fmt.Println(-1)
		return
	}

	number, ok := findNumber(a)
	if !ok {
		fmt.Println(-1)
		return
	}

	fmt.Println(number)
}
