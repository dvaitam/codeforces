package main

import (
   "bufio"
   "fmt"
   "os"
)

// Node represents a segment tree node for bracket matching
type Node struct {
   pairs int // number of matched pairs
   open  int // number of unmatched '('
   close int // number of unmatched ')'
}

// merge combines two nodes: a corresponds to left segment, b to right segment
func merge(a, b Node) Node {
   // match some of a.open with b.close
   match := a.open
   if b.close < match {
       match = b.close
   }
   return Node{
       pairs: a.pairs + b.pairs + match,
       open:  a.open + b.open - match,
       close: a.close + b.close - match,
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // build segment tree of size N = next power of two
   N := 1
   for N < n {
       N <<= 1
   }
   size := 2 * N
   tree := make([]Node, size)
   // initialize leaves
   for i := 0; i < n; i++ {
       idx := N + i
       if s[i] == '(' {
           tree[idx] = Node{pairs: 0, open: 1, close: 0}
       } else {
           tree[idx] = Node{pairs: 0, open: 0, close: 1}
       }
   }
   // build internal nodes
   for i := N - 1; i > 0; i-- {
       tree[i] = merge(tree[2*i], tree[2*i+1])
   }
   // process queries
   var m int
   fmt.Fscan(reader, &m)
   for q := 0; q < m; q++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       // convert to 0-based
       l--
       r--
       // query on [l, r]
       l += N
       r += N
       // accumulators
       resL := Node{}
       resR := Node{}
       for l <= r {
           if l&1 == 1 {
               resL = merge(resL, tree[l])
               l++
           }
           if r&1 == 0 {
               resR = merge(tree[r], resR)
               r--
           }
           l >>= 1
           r >>= 1
       }
       res := merge(resL, resR)
       // each pair contributes 2 characters
       fmt.Fprintln(writer, res.pairs*2)
   }
}
