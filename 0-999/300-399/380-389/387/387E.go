package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// BIT implements a Fenwick tree for ints (1-based)
type BIT struct {
   n    int
   tree []int
}

// NewBIT creates a BIT of size n
func NewBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]int, n+1)}
}

// Add adds v at position i
func (b *BIT) Add(i, v int) {
   for ; i <= b.n; i += i & -i {
       b.tree[i] += v
   }
}

// Sum returns prefix sum up to i
func (b *BIT) Sum(i int) int {
   s := 0
   for ; i > 0; i -= i & -i {
       s += b.tree[i]
   }
   return s
}

// FindByOrder finds the smallest index i such that Sum(i) >= k (1 <= k <= total sum)
func (b *BIT) FindByOrder(k int) int {
   idx := 0
   bitMask := 1
   // highest power of two <= n
   for bitMask<<1 <= b.n {
       bitMask <<= 1
   }
   for d := bitMask; d > 0; d >>= 1 {
       next := idx + d
       if next <= b.n && b.tree[next] < k {
           idx = next
           k -= b.tree[next]
       }
   }
   return idx + 1
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   fmt.Fscan(reader, &n, &k)
   p := make([]int, n+1)
   pos := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &p[i])
       pos[p[i]] = i
   }
   bvals := make([]int, k)
   isKeeper := make([]bool, n+1)
   for i := 0; i < k; i++ {
       fmt.Fscan(reader, &bvals[i])
       isKeeper[bvals[i]] = true
   }
   // blockers: keepers sorted by value
   type KV struct{ val, pos int }
   bElems := make([]KV, 0, k)
   for _, v := range bvals {
       bElems = append(bElems, KV{v, pos[v]})
   }
   sort.Slice(bElems, func(i, j int) bool {
       return bElems[i].val < bElems[j].val
   })
   // removable elements sorted by value implicitly
   remElems := make([]KV, 0, n-k)
   for v := 1; v <= n; v++ {
       if !isKeeper[v] {
           remElems = append(remElems, KV{v, pos[v]})
       }
   }

   remBIT := NewBIT(n)
   for i := 1; i <= n; i++ {
       remBIT.Add(i, 1)
   }
   blockerBIT := NewBIT(n)

   ans := int64(0)
   idxB := 0
   totalB := len(bElems)
   for _, item := range remElems {
       x := item.val
       px := item.pos
       // activate blockers with value < x
       for idxB < totalB && bElems[idxB].val < x {
           blockerBIT.Add(bElems[idxB].pos, 1)
           idxB++
       }
       // find left boundary
       sumLeft := blockerBIT.Sum(px - 1)
       var leftBound int
       if sumLeft == 0 {
           leftBound = 0
       } else {
           leftBound = blockerBIT.FindByOrder(sumLeft)
       }
       // find right boundary
       sumBefore := blockerBIT.Sum(px)
       totalBl := blockerBIT.Sum(n)
       var rightBound int
       if sumBefore < totalBl {
           rightBound = blockerBIT.FindByOrder(sumBefore + 1)
       } else {
           rightBound = n + 1
       }
       // count remaining elements in (leftBound, rightBound)
       l := leftBound + 1
       r := rightBound - 1
       count := 0
       if l <= r {
           count = remBIT.Sum(r) - remBIT.Sum(leftBound)
       }
       ans += int64(count)
       // remove this element
       remBIT.Add(px, -1)
   }
   fmt.Fprintln(writer, ans)
}
