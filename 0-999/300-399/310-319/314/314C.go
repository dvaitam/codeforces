package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

// Fenwick tree for sum modulo MOD
type Fenwick struct {
   n   int
   bit []int
}

func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, bit: make([]int, n+1)}
}

// add v at index i (1-based)
func (f *Fenwick) Add(i int, v int) {
   for x := i; x <= f.n; x += x & -x {
       f.bit[x] += v
       if f.bit[x] >= MOD {
           f.bit[x] -= MOD
       }
   }
}

// sum of [1..i]
func (f *Fenwick) Sum(i int) int {
   s := 0
   for x := i; x > 0; x -= x & -x {
       s += f.bit[x]
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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   maxv := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i] > maxv {
           maxv = a[i]
       }
   }
   // Fenwick over values up to maxv
   fw := NewFenwick(maxv)
   last := make([]int, maxv+1)
   ans := 0
   for _, v := range a {
       // sum of contributions ending with u <= v
       s := fw.Sum(v)
       cur := int((int64(s+1) * int64(v)) % MOD)
       // delta excludes duplicates from previous same value
       delta := cur - last[v]
       if delta < 0 {
           delta += MOD
       }
       last[v] = cur
       // add delta to fw at position v
       fw.Add(v, delta)
       ans += delta
       if ans >= MOD {
           ans -= MOD
       }
   }
   fmt.Fprintln(writer, ans)
