package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type order struct {
	dir   string
	price int
	qty   int
}

type caseB struct {
	n      int
	s      int
	orders []order
}

func genCase(rng *rand.Rand) caseB {
	n := rng.Intn(20) + 1
	s := rng.Intn(5) + 1
	orders := make([]order, n)
	for i := 0; i < n; i++ {
		dir := "B"
		if rng.Intn(2) == 0 {
			dir = "S"
		}
		var price int
		if dir == "B" {
			price = rng.Intn(40) + 1
		} else {
			price = rng.Intn(50) + 50
		}
		qty := rng.Intn(20) + 1
		orders[i] = order{dir, price, qty}
	}
	return caseB{n, s, orders}
}

func formatInput(tc caseB) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.s))
	for _, ord := range tc.orders {
		sb.WriteString(fmt.Sprintf("%s %d %d\n", ord.dir, ord.price, ord.qty))
	}
	return sb.String()
}

func expected(tc caseB) []string {
	buy := make(map[int]int)
	sell := make(map[int]int)
	for _, o := range tc.orders {
		if o.dir == "B" {
			buy[o.price] += o.qty
		} else {
			sell[o.price] += o.qty
		}
	}
	buyPrices := make([]int, 0, len(buy))
	for p := range buy {
		buyPrices = append(buyPrices, p)
	}
	sellPrices := make([]int, 0, len(sell))
	for p := range sell {
		sellPrices = append(sellPrices, p)
	}
	sort.Ints(buyPrices)
	sort.Ints(sellPrices)
	res := make([]string, 0)
	limitSell := tc.s
	if limitSell > len(sellPrices) {
		limitSell = len(sellPrices)
	}
	for i := limitSell - 1; i >= 0; i-- {
		p := sellPrices[i]
		res = append(res, fmt.Sprintf("S %d %d", p, sell[p]))
	}
	limitBuy := tc.s
	if limitBuy > len(buyPrices) {
		limitBuy = len(buyPrices)
	}
	for i := len(buyPrices) - 1; i >= len(buyPrices)-limitBuy; i-- {
		p := buyPrices[i]
		res = append(res, fmt.Sprintf("B %d %d", p, buy[p]))
	}
	return res
}

func runCase(bin string, tc caseB) error {
	input := formatInput(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	expectedLines := expected(tc)
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(out.String())))
	gotLines := make([]string, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			gotLines = append(gotLines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read output: %v", err)
	}
	if len(gotLines) != len(expectedLines) {
		return fmt.Errorf("expected %d lines got %d", len(expectedLines), len(gotLines))
	}
	for i, exp := range expectedLines {
		fields := strings.Fields(gotLines[i])
		var dir string
		var price, qty int
		if len(fields) != 3 {
			return fmt.Errorf("bad line %d: %q", i+1, gotLines[i])
		}
		if _, err := fmt.Sscan(gotLines[i], &dir, &price, &qty); err != nil {
			return fmt.Errorf("bad numbers on line %d: %v", i+1, err)
		}
		eFields := strings.Fields(exp)
		eDir := eFields[0]
		var ePrice, eQty int
		fmt.Sscan(exp, &eDir, &ePrice, &eQty)
		if dir != eDir || price != ePrice || qty != eQty {
			return fmt.Errorf("line %d expected '%s' got '%s'", i+1, exp, gotLines[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, formatInput(tc))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
