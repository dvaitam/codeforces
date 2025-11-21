package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	mod   = 1_000_000_007
	maxN  = 400000
	free  = 0
	zero  = 1
	fixed = 2
	pair  = 3
)

var totalWays []int64

func init() {
	totalWays = make([]int64, maxN+1)
	totalWays[0] = 1
	if maxN >= 1 {
		totalWays[1] = 2
	}
	for i := 2; i <= maxN; i++ {
		totalWays[i] = (2*totalWays[i-1] + int64(i-1)*totalWays[i-2]) % mod
	}
}

func main() {
	data, _ := io.ReadAll(os.Stdin)
	idx := 0
	nextInt := func() int {
		sign := 1
		val := 0
		for idx < len(data) && (data[idx] < '0' || data[idx] > '9') && data[idx] != '-' {
			idx++
		}
		if idx < len(data) && data[idx] == '-' {
			sign = -1
			idx++
		}
		for idx < len(data) && data[idx] >= '0' && data[idx] <= '9' {
			val = val*10 + int(data[idx]-'0')
			idx++
		}
		return sign * val
	}

	t := nextInt()
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for ; t > 0; t-- {
		n := nextInt()
		a := make([]int, n+1)
		for i := 0; i <= n; i++ {
			a[i] = nextInt()
		}

		status := make([]byte, n+1)
		partner := make([]int, n+1)
		ok := true

		for i := 1; i <= n && ok; i++ {
			v := a[i]
			if v == -1 {
				continue
			}
			switch {
			case v == 0:
				ok = setZero(i, status, partner)
			case v == i:
				ok = setFixed(i, status, partner)
			case v >= 1 && v <= n:
				ok = setPair(i, v, status, partner)
			default:
				ok = false
			}
		}

		if ok && status[n] == zero {
			ok = false
		}

		freeCnt := 0
		if ok {
			for i := 1; i <= n; i++ {
				if status[i] == free {
					freeCnt++
				}
			}
		}

		var ans int64
		if !ok {
			ans = 0
		} else if status[n] == free {
			if freeCnt == 0 {
				ans = 0
			} else {
				ans = totalWays[freeCnt] - totalWays[freeCnt-1]
			}
		} else {
			ans = totalWays[freeCnt]
		}

		ans %= mod
		if ans < 0 {
			ans += mod
		}
		fmt.Fprintln(writer, ans)
	}
}

func setZero(i int, status []byte, partner []int) bool {
	switch status[i] {
	case zero:
		return true
	case fixed, pair:
		return false
	default:
		if partner[i] != 0 {
			return false
		}
		status[i] = zero
		return true
	}
}

func setFixed(i int, status []byte, partner []int) bool {
	switch status[i] {
	case fixed:
		return true
	case zero, pair:
		return false
	default:
		if partner[i] != 0 {
			return false
		}
		status[i] = fixed
		return true
	}
}

func setPair(i, j int, status []byte, partner []int) bool {
	if i == j {
		return setFixed(i, status, partner)
	}
	if j <= 0 || j >= len(status) {
		return false
	}
	if status[i] == zero || status[i] == fixed {
		return false
	}
	if status[j] == zero || status[j] == fixed {
		return false
	}
	if partner[i] != 0 && partner[i] != j {
		return false
	}
	if partner[j] != 0 && partner[j] != i {
		return false
	}
	partner[i] = j
	partner[j] = i
	status[i] = pair
	status[j] = pair
	return true
}
