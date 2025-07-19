package main

import (
	"bufio"
	"os"
	"strconv"
)

func readInt(r *bufio.Reader) int64 {
	var x int64
	var sign int64 = 1
	b, err := r.ReadByte()
	for err == nil && b != '-' && (b < '0' || b > '9') {
		b, err = r.ReadByte()
	}
	if b == '-' {
		sign = -1
		b, _ = r.ReadByte()
	}
	for ; err == nil && b >= '0' && b <= '9'; b, err = r.ReadByte() {
		x = x*10 + int64(b-'0')
	}
	return x * sign
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	n := readInt(reader)
	q := readInt(reader)
	halfRow := n / 2
	halfBoard := (n * n) / 2
	evenN := (n % 2) == 0

	for k := int64(0); k < q; k++ {
		i := readInt(reader)
		j := readInt(reader)
		var ans int64
		if evenN {
			ans = ((i-1)*n)/2 + j/2 + j%2
			if (i+j)%2 != 0 {
				ans += halfBoard
			}
		} else {
			if (i+j)%2 == 0 {
				if i%2 == 1 {
					ans = ((i-1)/2)*n + j/2 + 1
				} else {
					ans = ((i-2)/2)*n + halfRow + 1 + j/2
				}
			} else {
				if i%2 == 1 {
					ans = ((i-1)/2)*n + j/2 + halfBoard + 1
				} else {
					ans = ((i-2)/2)*n + halfRow + j/2 + 1 + halfBoard + 1
				}
			}
		}
		writer.WriteString(strconv.FormatInt(ans, 10))
		writer.WriteByte('\n')
	}
}
