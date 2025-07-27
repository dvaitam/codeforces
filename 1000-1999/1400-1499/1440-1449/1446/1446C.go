package main

import (
   "bufio"
   "fmt"
   "os"
)

// Node in binary trie
type node struct {
   ch  [2]int
   cnt int
}

var nodes []node
var tot int

// insert number x into trie
func insert(x int) {
   idx := 1 // root
   nodes[idx].cnt++
   for b := 30; b >= 0; b-- {
       bit := (x >> b) & 1
       nxt := nodes[idx].ch[bit]
       if nxt == 0 {
           tot++
           nodes = append(nodes, node{})
           nxt = tot
           nodes[idx].ch[bit] = nxt
       }
       idx = nxt
       nodes[idx].cnt++
   }
}

// f returns maximum size of good subsequence in trie rooted at idx, considering bits up to b
func f(idx, b int) int {
   if idx == 0 || nodes[idx].cnt == 0 {
       return 0
   }
   if nodes[idx].cnt == 1 || b < 0 {
       return 1
   }
   l := nodes[idx].ch[0]
   r := nodes[idx].ch[1]
   if l == 0 {
       return f(r, b-1)
   }
   if r == 0 {
       return f(l, b-1)
   }
   left := f(l, b-1)
   right := f(r, b-1)
   if left > right {
       return left + 1
   }
   return right + 1
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // init trie: reserve root at 1
   nodes = make([]node, 2)
   tot = 1
   for _, x := range a {
       insert(x)
   }
   best := f(1, 30)
   // need at least 2 remain
   if best < 2 {
       best = 2
   }
   // answer: removals
   fmt.Fprintln(out, n-best)
}
