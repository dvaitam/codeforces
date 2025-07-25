package main

import (
   "bufio"
   "fmt"
   "os"
)

type Node struct {
   vec     []int
   extra   int
}

var (
   m    int
   tree []Node
)

func merge(a, b Node) Node {
   // a: left, b: right
   // b.extra pops apply to a.vec
   ra := b.extra
   av := a.vec
   if ra >= len(av) {
       // all a.vec popped
       ra2 := ra - len(av)
       // leftover pops = a.extra + ra2
       return Node{vec: append([]int(nil), b.vec...), extra: a.extra + ra2}
   }
   // partial pops
   newLeft := make([]int, len(av)-ra)
   copy(newLeft, av[:len(av)-ra])
   // combine
   resVec := append(newLeft, b.vec...)
   return Node{vec: resVec, extra: a.extra}
}

func build(n int) {
   size := 1
   for size < n {
       size <<= 1
   }
   tree = make([]Node, 2*size)
}

func update(n int, pos int, nd Node) {
   // tree base at size
   size := len(tree) / 2
   idx := pos + size - 1
   tree[idx] = nd
   for idx >>= 1; idx > 0; idx >>= 1 {
       tree[idx] = merge(tree[2*idx], tree[2*idx+1])
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &m)
   // initialize tree
   size := 1
   for size < m {
       size <<= 1
   }
   tree = make([]Node, 2*size)
   // process queries
   for i := 0; i < m; i++ {
       var p, t, x int
       fmt.Fscan(in, &p, &t)
       var nd Node
       if t == 0 {
           nd = Node{vec: nil, extra: 1}
       } else {
           fmt.Fscan(in, &x)
           nd = Node{vec: []int{x}, extra: 0}
       }
       update(m, p, nd)
       root := tree[1]
       if len(root.vec) == 0 {
           out.WriteString("-1 ")
       } else {
           out.WriteString(fmt.Sprintf("%d ", root.vec[len(root.vec)-1]))
       }
   }
}
