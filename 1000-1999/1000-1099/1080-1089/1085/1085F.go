package main

import (
   "bufio"
   "fmt"
   "os"
)

// BIT for 1-indexed array
type BIT struct {
   n    int
   tree []int
}

func NewBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]int, n+1)}
}

// add v at position i
func (b *BIT) Add(i, v int) {
   for x := i; x <= b.n; x += x & -x {
       b.tree[x] += v
   }
}

// Sum returns prefix sum [1..i]
func (b *BIT) Sum(i int) int {
   s := 0
   for x := i; x > 0; x -= x & -x {
       s += b.tree[x]
   }
   return s
}

// Kth returns smallest i such that Sum(i) >= k, assumes k>=1 and k<=Sum(n)
func (b *BIT) Kth(k int) int {
   pos, acc := 0, 0
   // compute max power of two <= n
   maxPow := 1
   for maxPow<<1 <= b.n {
       maxPow <<= 1
   }
   for pw := maxPow; pw > 0; pw >>= 1 {
       np := pos + pw
       if np <= b.n && acc+b.tree[np] < k {
           acc += b.tree[np]
           pos = np
       }
   }
   return pos + 1
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, q int
   fmt.Fscan(in, &n, &q)
   sbytes := make([]byte, n+1)
   var tmp string
   fmt.Fscan(in, &tmp)
   for i := 1; i <= n; i++ {
       sbytes[i] = tmp[i-1]
   }
   // map types: R=0, S=1, P=2
   idx := map[byte]int{'R': 0, 'S': 1, 'P': 2}
   win := [3]int{1, 2, 0}
   lose := [3]int{2, 0, 1}
   bits := make([]*BIT, 3)
   cnt := [3]int{}
   for t := 0; t < 3; t++ {
       bits[t] = NewBIT(n)
   }
   for i := 1; i <= n; i++ {
       t := idx[sbytes[i]]
       bits[t].Add(i, 1)
       cnt[t]++
   }
   // function to compute answer
   calc := func() int {
       res := 0
       for t := 0; t < 3; t++ {
           if cnt[lose[t]] == 0 {
               res += cnt[t]
           } else if cnt[win[t]] == 0 {
               // none
           } else {
               // both exist
               // lose positions
               l1 := bits[lose[t]].Kth(1)
               l2 := bits[lose[t]].Kth(cnt[lose[t]])
               w1 := bits[win[t]].Kth(1)
               w2 := bits[win[t]].Kth(cnt[win[t]])
               lo := min(l1, w1)
               hi := max(l2, w2)
               before := bits[t].Sum(lo - 1)
               after := bits[t].Sum(n) - bits[t].Sum(hi)
               res += before + after
           }
       }
       return res
   }
   // initial
   fmt.Fprintln(out, calc())
   for i := 0; i < q; i++ {
       var p int
       var c byte
       fmt.Fscan(in, &p, &c)
       old := idx[sbytes[p]]
       if byte(c) == sbytes[p] {
           fmt.Fprintln(out, calc())
           continue
       }
       // remove old
       bits[old].Add(p, -1)
       cnt[old]--
       // add new
       sbytes[p] = c
       t := idx[c]
       bits[t].Add(p, 1)
       cnt[t]++
       fmt.Fprintln(out, calc())
   }
}
