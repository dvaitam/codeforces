package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
   "sort"
)

// Fenwick tree (BIT) for int64
type BIT struct {
   n    int
   tree []int64
}

func NewBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]int64, n+1)}
}

// Add v at position i (1-based)
func (b *BIT) Add(i int, v int64) {
   for ; i <= b.n; i += i & -i {
       b.tree[i] += v
   }
}

// Sum returns prefix sum up to i (1-based)
func (b *BIT) Sum(i int) int64 {
   var s int64
   for ; i > 0; i -= i & -i {
       s += b.tree[i]
   }
   return s
}

// FindByPrefix returns smallest index i such that prefix sum >= k (k>=1)
func (b *BIT) FindByPrefix(k int64) int {
   pos := 0
   // largest power of two >= n
   maxLog := bits.Len(uint(b.n))
   for d := 1 << (maxLog - 1); d > 0; d >>= 1 {
       if pos+d <= b.n && b.tree[pos+d] < k {
           k -= b.tree[pos+d]
           pos += d
       }
   }
   return pos + 1
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   type op struct{ tp int; d int64 }
   ops := make([]op, n)
   vals := make([]int64, 0, n)
   for i := 0; i < n; i++ {
       var tp int
       var d int64
       fmt.Fscan(reader, &tp, &d)
       ops[i] = op{tp, d}
       if d < 0 {
           d = -d
       }
       vals = append(vals, d)
   }
   // coordinate compress
   sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
   uniq := vals[:1]
   for i := 1; i < len(vals); i++ {
       if vals[i] != vals[i-1] {
           uniq = append(uniq, vals[i])
       }
   }
   m := len(uniq)
   index := make(map[int64]int, m)
   for i, v := range uniq {
       index[v] = i + 1
   }
   bitCnt := NewBIT(m)
   bitSum := NewBIT(m)
   var totalCnt, lightCnt int
   var sumAll int64
   for _, o := range ops {
       tp, d := o.tp, o.d
       var x int64
       if d < 0 {
           x = -d
       } else {
           x = d
       }
       idx := index[x]
       if d > 0 {
           // add spell
           totalCnt++
           bitCnt.Add(idx, 1)
           bitSum.Add(idx, x)
           sumAll += x
           if tp == 1 {
               lightCnt++
           }
       } else {
           // remove spell
           totalCnt--
           bitCnt.Add(idx, -1)
           bitSum.Add(idx, -x)
           sumAll -= x
           if tp == 1 {
               lightCnt--
           }
       }
       // tokens applied = min(lightCnt, totalCnt-1)
       var tokens int
       if totalCnt > 0 {
           tokens = lightCnt
           if tokens > totalCnt-1 {
               tokens = totalCnt - 1
           }
       }
       // p = number of smallest kept uninterested = totalCnt - tokens
       p := totalCnt - tokens
       var sumSmall int64
       if p > 0 {
           // find sum of smallest p values
           idx0 := bitCnt.FindByPrefix(int64(p))
           sumBefore := bitSum.Sum(idx0 - 1)
           cntBefore := bitCnt.Sum(idx0 - 1)
           rem := int64(p) - cntBefore
           sumSmall = sumBefore + uniq[idx0-1]*rem
       }
       sumTop := sumAll - sumSmall
       ans := sumAll + sumTop
       fmt.Fprintln(writer, ans)
   }
}
