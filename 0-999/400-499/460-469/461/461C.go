package main

import (
   "bufio"
   "fmt"
   "os"
)

type BIT struct {
   n   int
   bit []int64
}

func NewBIT(n int) *BIT {
   return &BIT{n: n, bit: make([]int64, n+1)}
}

func (b *BIT) Add(i int, v int64) {
   // i: 1-indexed
   for x := i; x <= b.n; x += x & -x {
       b.bit[x] += v
   }
}

func (b *BIT) Sum(i int) int64 {
   // sum [1..i]
   var s int64
   for x := i; x > 0; x -= x & -x {
       s += b.bit[x]
   }
   return s
}

func (b *BIT) RangeSum(l, r int) int64 {
   if l > r {
       return 0
   }
   return b.Sum(r) - b.Sum(l-1)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, q int
   fmt.Fscan(in, &n, &q)
   bit := NewBIT(n)
   for i := 1; i <= n; i++ {
       bit.Add(i, 1)
   }
   l, r := 0, n
   rev := false
   for qi := 0; qi < q; qi++ {
       var t int
       fmt.Fscan(in, &t)
       if t == 1 {
           var p int
           fmt.Fscan(in, &p)
           w := r - l
           if p <= w-p {
               // fold left p onto right
               for i := 0; i < p; i++ {
                   var x1, x2 int
                   if !rev {
                       x1 = l + i
                       x2 = l + (w-1-i)
                   } else {
                       x1 = r - 1 - i
                       x2 = r - 1 - (w-1-i)
                   }
                   // bit positions are 1-indexed
                   v := bit.RangeSum(x1+1, x1+1)
                   bit.Add(x2+1, v)
               }
               // drop left p
               if !rev {
                   l += p
               } else {
                   r -= p
               }
           } else {
               // fold right (w-p) onto left
               sz := w - p
               for i := 0; i < sz; i++ {
                   var x1, x2 int
                   // rightmost sz pos mapping to left
                   if !rev {
                       x1 = l + (w-1-i)
                       x2 = l + i
                   } else {
                       x1 = r - 1 - (w-1-i)
                       x2 = r - 1 - i
                   }
                   v := bit.RangeSum(x1+1, x1+1)
                   bit.Add(x2+1, v)
               }
               // drop right sz
               if !rev {
                   r -= (w - p)
               } else {
                   l += (w - p)
               }
               rev = !rev
           }
       } else if t == 2 {
           var li, ri int
           fmt.Fscan(in, &li, &ri)
           // sum over [li, ri)
           var a, b int
           if !rev {
               a = l + li
               b = l + ri - 1
           } else {
               a = r - 1 - li
               b = r - ri
           }
           if a > b {
               a, b = b, a
           }
           ans := bit.RangeSum(a+1, b+1)
           fmt.Fprintln(out, ans)
       }
   }
}
