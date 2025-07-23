package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type order struct {
	price int
	qty   int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, s int
	if _, err := fmt.Fscan(reader, &n, &s); err != nil {
		return
	}

	buy := make(map[int]int)
	sell := make(map[int]int)

	for i := 0; i < n; i++ {
		var dir string
		var p, q int
		fmt.Fscan(reader, &dir, &p, &q)
		if dir == "B" {
			buy[p] += q
		} else {
			sell[p] += q
		}
	}

	buys := make([]order, 0, len(buy))
	for price, qty := range buy {
		buys = append(buys, order{price, qty})
	}
	sells := make([]order, 0, len(sell))
	for price, qty := range sell {
		sells = append(sells, order{price, qty})
	}

	sort.Slice(buys, func(i, j int) bool { return buys[i].price > buys[j].price })
	sort.Slice(sells, func(i, j int) bool { return sells[i].price < sells[j].price })

	limitSell := s
	if limitSell > len(sells) {
		limitSell = len(sells)
	}
	for i := limitSell - 1; i >= 0; i-- {
		fmt.Fprintf(writer, "S %d %d\n", sells[i].price, sells[i].qty)
	}

	limitBuy := s
	if limitBuy > len(buys) {
		limitBuy = len(buys)
	}
	for i := 0; i < limitBuy; i++ {
		fmt.Fprintf(writer, "B %d %d\n", buys[i].price, buys[i].qty)
	}
}
