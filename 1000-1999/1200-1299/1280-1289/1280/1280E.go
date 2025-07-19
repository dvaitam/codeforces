package main

import (
	"bufio"
	"fmt"
	"os"
)

const MAXN = 100000

// node represents an operation or a resistor (*) in the circuit parse tree.
type node struct {
	c          [2]*node
	starNum    int
	isSeries   bool
	isParallel bool
	K          int64
}

var nodes [MAXN*2 + 5]node
var ans [MAXN + 5]int64

// getK computes K for each node: 1 for leaf, min for series, sum for parallel.
func (nd *node) getK() {
	if nd.starNum != 0 {
		nd.K = 1
	} else if nd.isSeries {
		nd.c[0].getK()
		nd.c[1].getK()
		if nd.c[0].K < nd.c[1].K {
			nd.K = nd.c[0].K
		} else {
			nd.K = nd.c[1].K
		}
	} else if nd.isParallel {
		nd.c[0].getK()
		nd.c[1].getK()
		nd.K = nd.c[0].K + nd.c[1].K
	}
}

// color distributes the total R*K along the tree, assigning values to leaves.
func (nd *node) color(RK int64) {
	if nd.starNum != 0 {
		ans[nd.starNum] = RK
	} else if nd.isSeries {
		if nd.c[0].K < nd.c[1].K {
			nd.c[0].color(RK)
			nd.c[1].color(0)
		} else {
			nd.c[1].color(RK)
			nd.c[0].color(0)
		}
	} else if nd.isParallel {
		nd.c[0].color(RK)
		nd.c[1].color(RK)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var R int64
		fmt.Fscan(reader, &R)
		// read rest of line as expression
		S, _ := reader.ReadString('\n')
		// reset for this test
		nodesCnt := 0
		N := 0
		type stItem struct {
			nd *node
			op byte
		}
		stack := make([]stItem, 0, 256)
		stack = append(stack, stItem{nil, 0})
		// parse expression
		for i := 0; i < len(S); i++ {
			c := S[i]
			if c == ' ' || c == '\n' || c == '\r' {
				continue
			} else if c == '(' {
				stack = append(stack, stItem{nil, 0})
			} else if c == ')' {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				curNode := top.nd
				prev := &stack[len(stack)-1]
				if prev.op == 0 {
					prev.nd = curNode
				} else {
					newNode := &nodes[nodesCnt]
					nodesCnt++
					*newNode = node{}
					newNode.c[0] = prev.nd
					newNode.c[1] = curNode
					if prev.op == 'S' {
						newNode.isSeries = true
					}
					if prev.op == 'P' {
						newNode.isParallel = true
					}
					prev.nd = newNode
					prev.op = 0
				}
			} else if c == '*' {
				curNode := &nodes[nodesCnt]
				nodesCnt++
				*curNode = node{}
				N++
				curNode.starNum = N
				prev := &stack[len(stack)-1]
				if prev.op == 0 {
					prev.nd = curNode
				} else {
					newNode := &nodes[nodesCnt]
					nodesCnt++
					*newNode = node{}
					newNode.c[0] = prev.nd
					newNode.c[1] = curNode
					if prev.op == 'S' {
						newNode.isSeries = true
					}
					if prev.op == 'P' {
						newNode.isParallel = true
					}
					prev.nd = newNode
					prev.op = 0
				}
			} else if c == 'S' || c == 'P' {
				stack[len(stack)-1].op = c
			}
		}
		root := stack[0].nd
		root.getK()
		root.color(R * root.K)
		writer.WriteString("REVOLTING")
		for i := 1; i <= N; i++ {
			writer.WriteByte(' ')
			writer.WriteString(fmt.Sprintf("%d", ans[i]))
		}
		writer.WriteByte('\n')
	}
}
