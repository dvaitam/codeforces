package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const MAXN = 100000

type node struct {
	c          [2]*node
	starNum    int
	isSeries   bool
	isParallel bool
	K          int64
}

var nodes [MAXN*2 + 5]node
var ans [MAXN + 5]int64

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

func solveE(R int64, expr string) []int64 {
	nodesCnt := 0
	N := 0
	type stItem struct {
		nd *node
		op byte
	}
	stack := make([]stItem, 0, 256)
	stack = append(stack, stItem{nil, 0})
	for i := 0; i < len(expr); i++ {
		c := expr[i]
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
	res := make([]int64, N)
	for i := 1; i <= N; i++ {
		res[i-1] = ans[i]
	}
	return res
}

func genExpr(depth int) string {
	if depth == 0 {
		return "*"
	}
	if rand.Intn(2) == 0 {
		return "*"
	}
	left := genExpr(depth - 1)
	right := genExpr(depth - 1)
	if rand.Intn(2) == 0 {
		return "(" + left + "S" + right + ")"
	}
	return "(" + left + "P" + right + ")"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	t := 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		R := rand.Intn(10) + 1
		expr := genExpr(3)
		fmt.Fprintln(&input, R, expr)
		res := solveE(int64(R), expr)
		line := "REVOLTING"
		for _, v := range res {
			line += fmt.Sprintf(" %d", v)
		}
		expected[i] = line
	}
	cmd := exec.Command(bin)
	cmd.Stdin = &input
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("binary error:", err)
		fmt.Print(string(out))
		return
	}
	outputs := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(outputs) != t {
		fmt.Printf("expected %d lines, got %d\n", t, len(outputs))
		fmt.Print(string(out))
		return
	}
	for i := 0; i < t; i++ {
		if strings.TrimSpace(outputs[i]) != expected[i] {
			fmt.Printf("Test %d failed: expected %s got %s\n", i+1, expected[i], strings.TrimSpace(outputs[i]))
			return
		}
	}
	fmt.Println("All tests passed!")
}
