package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	must   int
	s      [20]int
	length int
	ans    [][4]int
	cur    [4]int
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
)

func part(pos, cnt int) {
	if pos == length {
		if cnt == 4 {
			var ip [4]int
			copy(ip[:], cur[:])
			ans = append(ans, ip)
		}
		return
	}
	if cnt >= 4 {
		return
	}
	curval := 0
	for i := 0; i < 3; i++ {
		if pos+i >= length {
			break
		}
		curval = curval*10 + s[pos+i]
		if curval > 255 {
			break
		}
		if i > 0 && s[pos] == 0 {
			break
		}
		cur[cnt] = curval
		part(pos+i+1, cnt+1)
	}
}

func goRec(pos, mask int) {
	if pos >= 2 && pos <= 6 && mask == must {
		for iter := 0; iter < 2; iter++ {
			length = 2*pos - iter
			for i := 0; i < pos; i++ {
				s[length-1-i] = s[i]
			}
			part(0, 0)
		}
	}
	if pos == 6 {
		return
	}
	for d := 0; d < 10; d++ {
		if must&(1<<d) != 0 {
			s[pos] = d
			goRec(pos+1, mask|1<<d)
		}
	}
}

func main() {
	defer writer.Flush()
	var cnt int
	if _, err := fmt.Fscan(reader, &cnt); err != nil {
		return
	}
	for i := 0; i < cnt; i++ {
		var d int
		fmt.Fscan(reader, &d)
		must |= 1 << d
	}
	goRec(0, 0)
	fmt.Fprintln(writer, len(ans))
	for _, ip := range ans {
		fmt.Fprintf(writer, "%d.%d.%d.%d\n", ip[0], ip[1], ip[2], ip[3])
	}
}
