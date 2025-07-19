package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

// Fenwick tree for prefix sums and find by order
type BIT struct {
   n    int
   tree []int
}

// NewBIT creates a Fenwick tree of size n
func NewBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]int, n+1)}
}

// Update adds delta at position i
func (b *BIT) Update(i, delta int) {
   for i <= b.n {
       b.tree[i] += delta
       i += i & -i
   }
}

// Sum returns the prefix sum up to i
func (b *BIT) Sum(i int) int {
   if i > b.n {
       i = b.n
   }
   s := 0
   for i > 0 {
       s += b.tree[i]
       i -= i & -i
   }
   return s
}

// FindByOrder finds smallest index i such that sum(i) >= k (1-indexed)
func (b *BIT) FindByOrder(k int) int {
   idx := 0
   // largest power of two <= n
   bitMask := 1 << (bits.Len(uint(b.n)) - 1)
   for d := bitMask; d > 0; d >>= 1 {
       ni := idx + d
       if ni <= b.n && b.tree[ni] < k {
           idx = ni
           k -= b.tree[ni]
       }
   }
   return idx + 1
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       b := make([]int, n+2)
       for i := 1; i <= n; i++ {
           fmt.Fscan(reader, &b[i])
       }
       ok := true
       for i := 2; i <= n; i++ {
           if b[i] < b[i-1] {
               ok = false
               break
           }
       }
       if b[1] > n {
           ok = false
       }
       if !ok {
           writer.WriteString("No\n")
           continue
       }
       k := n + 1
       p := -1
       l, r := 1, n+1
       nxt := make([]int, n+2)
       vis := make([]bool, n+2)
       // initialize
       b[k] = k
       if b[1] >= 1 && b[1] <= n {
           vis[b[1]] = true
       }
       for i := 1; i <= n; i++ {
           if b[i] < b[i+1] {
               if b[i+1] <= n {
                   nxt[i] = b[i+1]
                   vis[nxt[i]] = true
               } else {
                   p = i
               }
           }
       }
       for l <= n && b[l] <= n {
           l++
       }
       a := make([]int, n+2)
       // check function returns -1,0,1
       check := func(kval, pval int) int {
           cnt := 0
           bit := NewBIT(n)
           for i := n; i >= 1; i-- {
               if nxt[i] != 0 {
                   a[i] = a[nxt[i]]
               } else if i >= kval || i == pval {
                   cnt++
                   a[i] = cnt
               } else {
                   c := bit.Sum(b[i])
                   if c == 0 {
                       return -1
                   }
                   idx := bit.FindByOrder(c)
                   a[i] = a[idx]
                   bit.Update(idx, -1)
               }
               if !vis[i] {
                   bit.Update(i, 1)
               }
           }
           total := bit.Sum(n)
           if total > 0 {
               idx := bit.FindByOrder(total)
               if idx > b[1] {
                   return 1
               }
           }
           writer.WriteString("Yes\n")
           for i := 1; i <= n; i++ {
               writer.WriteString(fmt.Sprintf("%d ", a[i]))
           }
           writer.WriteString("\n")
           return 0
       }
       found := false
       for l <= r {
           mid := (l + r) >> 1
           res := check(mid, p)
           if res == 0 {
               found = true
               break
           }
           if ^res != 0 {
               l = mid + 1
           } else {
               r = mid - 1
           }
       }
       if !found {
           writer.WriteString("No\n")
       }
   }
}
