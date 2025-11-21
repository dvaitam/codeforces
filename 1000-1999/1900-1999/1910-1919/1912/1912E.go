package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// reverseNumber returns the value obtained by reversing decimal digits of x.
func reverseNumber(x int64) int64 {
	var res int64
	for x > 0 {
		res = res*10 + x%10
		x /= 10
	}
	return res
}

// reverseString takes a non-negative integer and returns its decimal string reversed.
func reverseString(x int64) string {
	s := strconv.FormatInt(x, 10)
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

// buildExpression constructs the required expression for given p and q.
func buildExpression(p, q int64) string {
	// Pick a small positive shift so that neither operand ends with zero and
	// the resulting delta has even parity (so it can be fixed with gadgets).
	var shift int64 = -1
	for s := int64(1); s <= 1_000_000; s++ {
		var A, B int64
		if p >= 0 {
			A, B = p+s, s
		} else {
			A, B = s, s-p
		}
		if A%10 == 0 || B%10 == 0 {
			continue
		}
		delta := q - (reverseNumber(B) - reverseNumber(A))
		if delta%2 == 0 {
			shift = s
			break
		}
	}

	// Given the problem guarantee, the loop above should always succeed.
	if shift == -1 {
		shift = 1
	}

	var A, B int64
	if p >= 0 {
		A, B = p+shift, shift
	} else {
		A, B = shift, shift-p
	}

	delta := q - (reverseNumber(B) - reverseNumber(A))
	t := delta / 2 // amount we need to add to the reverse value

	// Break t into one or two pieces so that none ends with zero.
	parts := make([]int64, 0)
	if t > 0 {
		if t%10 != 0 {
			parts = append(parts, t)
		} else {
			parts = append(parts, 1, t-1)
		}
	} else if t < 0 {
		u := -t
		if u%10 != 0 {
			parts = append(parts, -u)
		} else {
			parts = append(parts, -int64(1), -(u - 1))
		}
	}

	var sb strings.Builder
	sb.WriteString(strconv.FormatInt(A, 10))
	sb.WriteByte('-')
	sb.WriteString(strconv.FormatInt(B, 10))

	for _, v := range parts {
		sb.WriteByte('+')
		xStr := reverseString(abs64(v))
		if v > 0 {
			// Gadget with forward value 0 and reverse value +2*xStr (i.e., +v*2).
			sb.WriteString(xStr)
			sb.WriteString("*1+")
			sb.WriteString(xStr)
			sb.WriteString("*1-")
			sb.WriteString(xStr)
			sb.WriteString("*2")
		} else {
			// Gadget with forward value 0 and reverse value -2*xStr.
			sb.WriteString(xStr)
			sb.WriteString("*2-")
			sb.WriteString(xStr)
			sb.WriteString("*1-")
			sb.WriteString(xStr)
			sb.WriteString("*1")
		}
	}

	return sb.String()
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var p, q int64
	fmt.Fscan(in, &p, &q)

	expr := buildExpression(p, q)
	fmt.Print(expr)
}
