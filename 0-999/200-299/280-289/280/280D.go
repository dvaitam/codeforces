package main

import (
   "bufio"
   "fmt"
   "os"
)

const (
   MAXK   = 20
   NEG_INF = -1000000000
)

// Node represents a segment tree node
type Node struct {
   sum int
   f   [MAXK+1][4]int
}

var (
   n, m int
   arr  []int
   size int
   tree []Node
   identity Node
)

// merge two nodes
func merge(a, b Node) Node {
   var c Node
   c.sum = a.sum + b.sum
   // init f with NEG_INF
   for i := 0; i <= MAXK; i++ {
       for mask := 0; mask < 4; mask++ {
           c.f[i][mask] = NEG_INF
       }
   }
   // combine
   for sa := 0; sa <= MAXK; sa++ {
       for ma := 0; ma < 4; ma++ {
           va := a.f[sa][ma]
           if va <= NEG_INF { continue }
           for sb := 0; sb <= MAXK-sa; sb++ {
               for mb := 0; mb < 4; mb++ {
                   vb := b.f[sb][mb]
                   if vb <= NEG_INF { continue }
                   ns := sa + sb
                   // if connect across boundary
                   if (ma&2) != 0 && (mb&1) != 0 {
                       ns--
                   }
                   if ns < 0 || ns > MAXK { continue }
                   nm := (ma & 1) | (mb & 2)
                   sumv := va + vb
                   if sumv > c.f[ns][nm] {
                       c.f[ns][nm] = sumv
                   }
               }
           }
       }
   }
   return c
}

// build segment tree
func build() {
   // compute size as power of two
   size = 1
   for size < n {
       size <<= 1
   }
   tree = make([]Node, size*2)
   // identity node: sum=0, f[0][0]=0
   for i := 0; i <= MAXK; i++ {
       for mask := 0; mask < 4; mask++ {
           identity.f[i][mask] = NEG_INF
       }
   }
   identity.sum = 0
   identity.f[0][0] = 0
   // init leaves
   for i := 0; i < n; i++ {
       nd := &tree[size+i]
       // initialize
       nd.sum = arr[i]
       for j := 0; j <= MAXK; j++ {
           for mask := 0; mask < 4; mask++ {
               nd.f[j][mask] = NEG_INF
           }
       }
       nd.f[0][0] = 0
       if MAXK >= 1 {
           nd.f[1][3] = arr[i]
       }
   }
   // other leaves as identity
   for i := n; i < size; i++ {
       tree[size+i] = identity
   }
   // build internal
   for i := size - 1; i > 0; i-- {
       tree[i] = merge(tree[i<<1], tree[i<<1|1])
   }
}

// update position pos (0-indexed) to value v
func update(pos, v int) {
   idx := pos + size
   // set leaf
   nd := &tree[idx]
   nd.sum = v
   for j := 0; j <= MAXK; j++ {
       for mask := 0; mask < 4; mask++ {
           nd.f[j][mask] = NEG_INF
       }
   }
   nd.f[0][0] = 0
   if MAXK >= 1 {
       nd.f[1][3] = v
   }
   // update up
   for idx >>= 1; idx > 0; idx >>= 1 {
       tree[idx] = merge(tree[idx<<1], tree[idx<<1|1])
   }
}

// query interval [l,r) 0-indexed
func query(l, r int) Node {
   l += size
   r += size
   leftRes := identity
   rightRes := identity
   for l < r {
       if l&1 == 1 {
           leftRes = merge(leftRes, tree[l])
           l++
       }
       if r&1 == 1 {
           r--
           rightRes = merge(tree[r], rightRes)
       }
       l >>= 1
       r >>= 1
   }
   return merge(leftRes, rightRes)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n)
   arr = make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &arr[i])
   }
   build()
   fmt.Fscan(reader, &m)
   for i := 0; i < m; i++ {
       var typ int
       fmt.Fscan(reader, &typ)
       if typ == 0 {
           var pos, v int
           fmt.Fscan(reader, &pos, &v)
           update(pos-1, v)
       } else {
           var l, r, k int
           fmt.Fscan(reader, &l, &r, &k)
           res := query(l-1, r)
           // get max for <=k
           ans := 0
           for s := 1; s <= k; s++ {
               for mask := 0; mask < 4; mask++ {
                   v := res.f[s][mask]
                   if v > ans {
                       ans = v
                   }
               }
           }
           // also consider 0 segments => 0
           fmt.Fprintln(writer, ans)
       }
   }
}
