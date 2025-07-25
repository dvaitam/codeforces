package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

// Fenwick tree for int64 values
type Fenwick struct {
   n int
   t []int64
}

func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, t: make([]int64, n+1)}
}

func (f *Fenwick) Add(i int, v int64) {
   for x := i; x <= f.n; x += x & -x {
       f.t[x] += v
   }
}

func (f *Fenwick) Sum(i int) int64 {
   var s int64
   for x := i; x > 0; x -= x & -x {
       s += f.t[x]
   }
   return s
}

// find smallest idx such that sum(idx) >= target
func (f *Fenwick) LowerBound(target int64) int {
   idx := 0
   var sum int64
   // max bit: highest power of two <= n
   bit := 1
   for bit<<1 <= f.n {
       bit <<= 1
   }
   for step := bit; step > 0; step >>= 1 {
       nxt := idx + step
       if nxt <= f.n && sum+f.t[nxt] < target {
           sum += f.t[nxt]
           idx = nxt
       }
   }
   return idx + 1
}

// FenwickMod for int64 mod MOD
type FenwickMod struct {
   n int
   t []int64
}

func NewFenwickMod(n int) *FenwickMod {
   return &FenwickMod{n: n, t: make([]int64, n+1)}
}

func (f *FenwickMod) Add(i int, v int64) {
   v %= MOD
   if v < 0 {
       v += MOD
   }
   for x := i; x <= f.n; x += x & -x {
       f.t[x] += v
       if f.t[x] >= MOD {
           f.t[x] -= MOD
       }
   }
}

func (f *FenwickMod) Sum(i int) int64 {
   var s int64
   for x := i; x > 0; x -= x & -x {
       s += f.t[x]
       if s >= MOD {
           s -= MOD
       }
   }
   return s
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, q int
   fmt.Fscan(reader, &n, &q)
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   w := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &w[i])
   }
   // B[i] = a[i] - i
   B := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       B[i] = a[i] - int64(i)
   }
   bitW := NewFenwick(n)
   bitWB := NewFenwickMod(n)
   for i := 1; i <= n; i++ {
       bitW.Add(i, w[i])
       // B[i] may be negative
       bi := B[i] % MOD
       if bi < 0 {
           bi += MOD
       }
       bitWB.Add(i, w[i]%MOD*bi%MOD)
   }
   for qi := 0; qi < q; qi++ {
       var x, y int64
       fmt.Fscan(reader, &x, &y)
       if x < 0 {
           id := int(-x)
           neww := y
           delta := neww - w[id]
           w[id] = neww
           bitW.Add(id, delta)
           bi := B[id] % MOD
           if bi < 0 {
               bi += MOD
           }
           bitWB.Add(id, delta%MOD*bi%MOD)
       } else {
           l := int(x)
           r := int(y)
           // total weight
           tot := bitW.Sum(r) - bitW.Sum(l-1)
           half := (tot + 1) / 2
           base := bitW.Sum(l-1)
           m := bitW.LowerBound(base + half)
           if m > r {
               m = r
           }
           // sums
           sumWl := bitW.Sum(m) - bitW.Sum(l-1)
           sumWr := bitW.Sum(r) - bitW.Sum(m)
           sumWB_l := (bitWB.Sum(m) - bitWB.Sum(l-1) + MOD) % MOD
           sumWB_r := (bitWB.Sum(r) - bitWB.Sum(m) + MOD) % MOD
           med := B[m] % MOD
           if med < 0 {
               med += MOD
           }
           swl := sumWl % MOD
           swr := sumWr % MOD
           c1 := (med*swl%MOD - sumWB_l + MOD) % MOD
           c2 := (sumWB_r - med*swr%MOD + MOD) % MOD
           res := (c1 + c2) % MOD
           fmt.Fprintln(writer, res)
       }
   }
}
