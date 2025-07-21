package main

import (
   "bufio"
   "fmt"
   "os"
)

// Fenwick implements a 1-indexed Fenwick tree for sums and k-th order statistic
type Fenwick struct {
   n    int
   tree []int
   log  int
}

// NewFenwick creates a Fenwick tree of size n
func NewFenwick(n int) *Fenwick {
   // compute highest power of two >= n
   log := 1
   for (1 << log) <= n {
       log++
   }
   return &Fenwick{n: n, tree: make([]int, n+1), log: log}
}

// Add adds v at position i
func (f *Fenwick) Add(i, v int) {
   for idx := i; idx <= f.n; idx += idx & -idx {
       f.tree[idx] += v
   }
}

// Sum returns sum of [1..i]
func (f *Fenwick) Sum(i int) int {
   s := 0
   for idx := i; idx > 0; idx -= idx & -idx {
       s += f.tree[idx]
   }
   return s
}

// FindKth returns smallest i such that Sum(i) >= k, or f.n+1 if not found
func (f *Fenwick) FindKth(k int) int {
   pos := 0
   for d := 1 << (f.log - 1); d > 0; d >>= 1 {
       nxt := pos + d
       if nxt <= f.n && f.tree[nxt] < k {
           pos = nxt
           k -= f.tree[nxt]
       }
   }
   // pos is such that prefix sum up to pos < original k, so answer is pos+1
   return pos + 1
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   const MAXB = 17
   // per-bit data
   zeroFT := make([]*Fenwick, MAXB)
   oneSegSum := make([]int64, MAXB)
   hasOne := make([][]bool, MAXB)
   for b := 0; b < MAXB; b++ {
       zeroFT[b] = NewFenwick(n)
       hasOne[b] = make([]bool, n+2)
       // build zeros Fenwick and one segments sum
       var cur int
       for i := 1; i <= n; i++ {
           if (a[i]>>b)&1 == 0 {
               zeroFT[b].Add(i, 1)
               if cur > 0 {
                   oneSegSum[b] += int64(cur) * int64(cur+1) / 2
                   cur = 0
               }
           } else {
               hasOne[b][i] = true
               cur++
           }
       }
       if cur > 0 {
           oneSegSum[b] += int64(cur) * int64(cur+1) / 2
       }
   }
   // process queries
   for qi := 0; qi < m; qi++ {
       var p, v int
       fmt.Fscan(reader, &p, &v)
       old := a[p]
       if old != v {
           for b := 0; b < MAXB; b++ {
               ob := (old >> b) & 1
               nb := (v >> b) & 1
               if ob == nb {
                   continue
               }
               // handle bit b change at position p
               ft := zeroFT[b]
               // find prev zero and next zero positions
               // count zeros before p
               cL := ft.Sum(p - 1)
               var l0 int
               if cL > 0 {
                   l0 = ft.FindKth(cL)
               }
               totalZ := ft.Sum(n)
               // count zeros <= p
               var cP int
               if ob == 0 {
                   // currently zero at p
                   cP = cL + 1
               } else {
                   // currently one, so zeros <=p is same as cL
                   cP = cL
               }
               // next zero is cP+1 th
               var r0 int
               if cP+1 <= totalZ {
                   r0 = ft.FindKth(cP + 1)
               } else {
                   r0 = n + 1
               }
               // segment lengths
               if ob == 0 && nb == 1 {
                   // zero -> one: merge left and right segments
                   lenL := p - l0 - 1
                   lenR := r0 - p - 1
                   oldSum := int64(lenL) * int64(lenL+1) / 2
                   oldSum += int64(lenR) * int64(lenR+1) / 2
                   newLen := lenL + 1 + lenR
                   newSum := int64(newLen) * int64(newLen+1) / 2
                   oneSegSum[b] += newSum - oldSum
                   // remove zero at p
                   ft.Add(p, -1)
                   hasOne[b][p] = true
               } else if ob == 1 && nb == 0 {
                   // one -> zero: split segment
                   segLen := r0 - l0 - 1
                   oldSum := int64(segLen) * int64(segLen+1) / 2
                   lenL := p - l0 - 1
                   lenR := r0 - p - 1
                   newSum := int64(lenL) * int64(lenL+1) / 2
                   newSum += int64(lenR) * int64(lenR+1) / 2
                   oneSegSum[b] += newSum - oldSum
                   // add zero at p
                   ft.Add(p, 1)
                   hasOne[b][p] = false
               }
           }
           a[p] = v
       }
       // compute answer
       var ans int64
       for b := 0; b < MAXB; b++ {
           if oneSegSum[b] != 0 {
               ans += oneSegSum[b] * int64(1<<b)
           }
       }
       fmt.Fprintln(writer, ans)
   }
}
