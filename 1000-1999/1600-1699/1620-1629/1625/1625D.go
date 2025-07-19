package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Node represents a trie node with two children and a dp pair (length, index).
type Node struct {
   next   [2]int
   dpLen  int
   dpIdx  int
}

var (
   nodes []Node
   tsz   int
)

func newNode() Node {
   return Node{next: [2]int{-1, -1}, dpLen: 0, dpIdx: 0}
}

// Add inserts value x with dp pair (y, idx) into the trie.
func Add(x, y, idx int) {
   v := 0
   for bit := 30; bit >= 0; bit-- {
       z := (x >> bit) & 1
       if nodes[v].next[z] == -1 {
           nodes = append(nodes, newNode())
           nodes[v].next[z] = tsz
           tsz++
       }
       v = nodes[v].next[z]
       // update dp pair if better
       if y > nodes[v].dpLen || (y == nodes[v].dpLen && idx > nodes[v].dpIdx) {
           nodes[v].dpLen = y
           nodes[v].dpIdx = idx
       }
   }
}

// Get finds the best dp pair for x with constraint k.
func Get(x, k int) (bestLen, bestIdx int) {
   bestLen, bestIdx = 0, 0
   v := 0
   for bit := 30; bit >= 0; bit-- {
       zk := (k >> bit) & 1
       zx := (x >> bit) & 1
       if zk == 0 {
           opp := nodes[v].next[zx^1]
           if opp != -1 {
               nd := nodes[opp]
               if nd.dpLen > bestLen || (nd.dpLen == bestLen && nd.dpIdx > bestIdx) {
                   bestLen = nd.dpLen
                   bestIdx = nd.dpIdx
               }
           }
           v = nodes[v].next[zx]
       } else {
           v = nodes[v].next[zx^1]
       }
       if v == -1 {
           break
       }
   }
   if v != -1 {
       nd := nodes[v]
       if nd.dpLen > bestLen || (nd.dpLen == bestLen && nd.dpIdx > bestIdx) {
           bestLen = nd.dpLen
           bestIdx = nd.dpIdx
       }
   }
   return
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, k int
   fmt.Fscan(in, &n, &k)
   a := make([]struct{ val, idx int }, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i].val)
       a[i].idx = i
   }
   if k == 0 {
       // any sequence works
       fmt.Fprintln(out, n)
       for i := 1; i <= n; i++ {
           if i < n {
               fmt.Fprint(out, i, " ")
           } else {
               fmt.Fprintln(out, i)
           }
       }
       return
   }
   // sort by value
   sort.Slice(a, func(i, j int) bool { return a[i].val < a[j].val })

   // initialize trie
   nodes = make([]Node, 1, 31*(n+1))
   nodes[0] = newNode()
   tsz = 1

   parent := make([]int, n+1)
   ansLen, ansIdx := 0, 0
   // process
   for i := 0; i < n; i++ {
       mxLen, mxIdx := Get(a[i].val, k)
       parent[i+1] = mxIdx
       curLen := mxLen + 1
       if curLen > ansLen || (curLen == ansLen && i+1 > ansIdx) {
           ansLen = curLen
           ansIdx = i + 1
       }
       Add(a[i].val, curLen, i+1)
   }
   if ansLen <= 1 {
       fmt.Fprintln(out, -1)
       return
   }
   fmt.Fprintln(out, ansLen)
   // reconstruct sequence
   seq := make([]int, 0, ansLen)
   for j := ansIdx; j > 0; j = parent[j] {
       seq = append(seq, j)
   }
   // reverse seq
   for i, j := 0, len(seq)-1; i < j; i, j = i+1, j-1 {
       seq[i], seq[j] = seq[j], seq[i]
   }
   for i, si := range seq {
       // original index + 1
       out.WriteString(fmt.Sprintf("%d", a[si-1].idx+1))
       if i+1 < len(seq) {
           out.WriteByte(' ')
       }
   }
   out.WriteByte('\n')
}
