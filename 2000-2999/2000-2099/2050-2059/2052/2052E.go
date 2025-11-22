package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

func parseExpr(s string) (bool, *big.Int) {
	if len(s) == 0 {
		return false, nil
	}
	res := big.NewInt(0)
	sign := 1
	i := 0
	for {
		start := i
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
		}
		if start == i {
			return false, nil
		}
		// leading zero check
		if s[start] == '0' && i-start > 1 {
			return false, nil
		}
		if i-start > 10 {
			return false, nil
		}
		var num int64
		for k := start; k < i; k++ {
			num = num*10 + int64(s[k]-'0')
		}
		tmp := big.NewInt(num)
		if sign == 1 {
			res.Add(res, tmp)
		} else {
			res.Sub(res, tmp)
		}

		if i == len(s) {
			break
		}
		if s[i] != '+' && s[i] != '-' {
			return false, nil
		}
		if s[i] == '+' {
			sign = 1
		} else {
			sign = -1
		}
		i++
		if i == len(s) { // operator at end invalid
			return false, nil
		}
	}
	return true, res
}

func parseEquality(s string) (bool, *big.Int, *big.Int) {
	eq := -1
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			if eq != -1 {
				return false, nil, nil // more than one '='
			}
			eq = i
		} else if (s[i] < '0' || s[i] > '9') && s[i] != '+' && s[i] != '-' {
			return false, nil, nil
		}
	}
	if eq == -1 || eq == 0 || eq == len(s)-1 {
		return false, nil, nil
	}
	validL, left := parseExpr(s[:eq])
	if !validL {
		return false, nil, nil
	}
	validR, right := parseExpr(s[eq+1:])
	if !validR {
		return false, nil, nil
	}
	return true, left, right
}

func moveDigit(s string, from, to int) string {
	removed := s[from]
	t := s[:from] + s[from+1:]
	if to < 0 {
		to = 0
	}
	if to > len(t) {
		to = len(t)
	}
	return t[:to] + string(removed) + t[to:]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var str string
	if _, err := fmt.Fscan(in, &str); err != nil {
		return
	}

	// First, check if already correct.
	if ok, l, r := parseEquality(str); ok && l.Cmp(r) == 0 {
		fmt.Println("Correct")
		return
	}

	// Try moving one digit.
	for i := 0; i < len(str); i++ {
		if str[i] < '0' || str[i] > '9' {
			continue
		}
		for to := 0; to <= len(str); to++ {
			if to == i {
				continue
			}
			// Adjust position if removing before insertion point
			toPos := to
			if to > i {
				toPos = to - 1
			}
			ns := moveDigit(str, i, toPos)
			if ok, l, r := parseEquality(ns); ok && l.Cmp(r) == 0 {
				fmt.Println(ns)
				return
			}
		}
	}

	fmt.Println("Impossible")
}
