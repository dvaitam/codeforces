package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// BIT implements a Fenwick tree for prefix sums over ints.
type BIT struct {
   n int
   t []int
}

// NewBIT returns a BIT of size n (0..n-1).
func NewBIT(n int) *BIT {
   return &BIT{n: n, t: make([]int, n)}
}

// Add v at index i.
func (b *BIT) Add(i, v int) {
   for x := i; x < b.n; x |= x + 1 {
       b.t[x] += v
   }
}

// Sum returns the prefix sum [0..i].
func (b *BIT) Sum(i int) int {
   s := 0
   for x := i; x >= 0; x = (x&(x+1)) - 1 {
       s += b.t[x]
   }
   return s
}

// FindByPrefix finds the smallest index idx such that Sum(idx) >= k.
// Assumes all values are non-negative and total sum >= k.
func (b *BIT) FindByPrefix(k int) int {
   idx := -1
   // largest power of two <= n
   mask := 1
   for mask<<1 <= b.n {
       mask <<= 1
   }
   for ; mask > 0; mask >>= 1 {
       nxt := idx + mask
       if nxt < b.n && b.t[nxt] < k {
           k -= b.t[nxt]
           idx = nxt
       }
   }
   return idx + 1
}

type Interval struct {
   start, end int
   idx        int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   fmt.Fscan(reader, &n, &k)
   intervals := make([]Interval, n)
   coords := make([]int, 0, 2*n+1)
   coords = append(coords, -1)
   for i := 0; i < n; i++ {
       var s, l int
       fmt.Fscan(reader, &s, &l)
       intervals[i] = Interval{start: s, end: s + l, idx: i}
       coords = append(coords, s, s+l)
   }
   sort.Ints(coords)
   // unique
   m := 1
   for i := 1; i < len(coords); i++ {
       if coords[i] != coords[m-1] {
           coords[m] = coords[i]
           m++
       }
   }
   coords = coords[:m]
   // map value to compressed index
   comp := make(map[int]int, len(coords))
   for i, v := range coords {
       comp[v] = i
   }
   // sort intervals by end asc, start desc
   sort.Slice(intervals, func(i, j int) bool {
       if intervals[i].end != intervals[j].end {
           return intervals[i].end < intervals[j].end
       }
       return intervals[i].start > intervals[j].start
   })

   ans := make([]int, 0, k)
   // check if with mid buses we can schedule >= k intervals
   var check func(mid int, save bool) bool
   check = func(mid int, save bool) bool {
       bit := NewBIT(len(coords))
       // initialize mid buses available at time -1
       bit.Add(comp[-1], mid)
       cnt := 0
       if save {
           ans = ans[:0]
       }
       for _, iv := range intervals {
           ci := comp[iv.start]
           total := bit.Sum(ci)
           if total > 0 {
               // find the bus with largest time <= start
               pos := bit.FindByPrefix(total)
               // use this bus
               bit.Add(pos, -1)
               ce := comp[iv.end]
               bit.Add(ce, 1)
               cnt++
               if save && len(ans) < k {
                   ans = append(ans, iv.idx+1) // 1-based
               }
               if cnt >= k && !save {
                   // can early exit when not saving
                   return true
               }
           }
       }
       return cnt >= k
   }

   // binary search minimal buses
   l, r := 1, n
   for r-l > 1 {
       mid := (l + r) / 2
       if check(mid, false) {
           r = mid
       } else {
           l = mid
       }
   }
   if check(l, false) {
       r = l
   }
   // final with saving
   check(r, true)
   // output
   fmt.Fprintln(writer, r)
   for i, v := range ans {
       if i > 0 {
           writer.WriteString(" ")
       }
       writer.WriteString(fmt.Sprint(v))
   }
   writer.WriteByte('\n')
}
