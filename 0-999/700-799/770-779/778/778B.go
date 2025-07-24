package main

import (
	"bufio"
	"fmt"
	"os"
)

type Var struct {
	typ  int    // 0 const, 1 op
	bits string // for const
	op1  int
	op2  int
	op   string
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	vars := make([]Var, n)
	idx := make(map[string]int)

	for i := 0; i < n; i++ {
		var name, tmp, t1 string
		fmt.Fscan(in, &name, &tmp, &t1)
		if t1[0] == '0' || t1[0] == '1' {
			vars[i] = Var{typ: 0, bits: t1}
		} else {
			var op, t2 string
			fmt.Fscan(in, &op, &t2)
			o1 := -1
			if t1 != "?" {
				o1 = idx[t1]
			}
			o2 := -1
			if t2 != "?" {
				o2 = idx[t2]
			}
			vars[i] = Var{typ: 1, op1: o1, op2: o2, op: op}
		}
		idx[name] = i
	}

	minAns := make([]byte, m)
	maxAns := make([]byte, m)
	vals0 := make([]int, n)
	vals1 := make([]int, n)

	for bit := 0; bit < m; bit++ {
		sum0, sum1 := 0, 0
		for i := 0; i < n; i++ {
			v := vars[i]
			if v.typ == 0 {
				b := int(v.bits[bit] - '0')
				vals0[i], vals1[i] = b, b
			} else {
				var a0, a1, b0, b1 int
				if v.op1 == -1 {
					a0, a1 = 0, 1
				} else {
					a0, a1 = vals0[v.op1], vals1[v.op1]
				}
				if v.op2 == -1 {
					b0, b1 = 0, 1
				} else {
					b0, b1 = vals0[v.op2], vals1[v.op2]
				}
				switch v.op {
				case "AND":
					vals0[i] = a0 & b0
					vals1[i] = a1 & b1
				case "OR":
					vals0[i] = a0 | b0
					vals1[i] = a1 | b1
				case "XOR":
					vals0[i] = a0 ^ b0
					vals1[i] = a1 ^ b1
				}
			}
			sum0 += vals0[i]
			sum1 += vals1[i]
		}
		if sum0 <= sum1 {
			minAns[bit] = '0'
		} else {
			minAns[bit] = '1'
		}
		if sum1 > sum0 {
			maxAns[bit] = '1'
		} else {
			maxAns[bit] = '0'
		}
	}

	fmt.Fprintln(out, string(minAns))
	fmt.Fprintln(out, string(maxAns))
}
