package main

import (
   "bufio"
   "fmt"
   "os"
)

// 2D BIT for XOR with parity separation
type BIT2D struct {
   n int
   t [][]uint64
}

func NewBIT2D(n int) *BIT2D {
   t := make([][]uint64, n+2)
   for i := range t {
       t[i] = make([]uint64, n+2)
   }
   return &BIT2D{n: n + 1, t: t}
}

// point update: t[x][y] ^= v
func (b *BIT2D) update(x, y int, v uint64) {
   for i := x; i <= b.n; i += i & -i {
       row := b.t[i]
       for j := y; j <= b.n; j += j & -j {
           row[j] ^= v
       }
   }
}

// prefix query XOR over [1..x][1..y]
func (b *BIT2D) query(x, y int) uint64 {
   var res uint64
   for i := x; i > 0; i -= i & -i {
       row := b.t[i]
       for j := y; j > 0; j -= j & -j {
           res ^= row[j]
       }
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m int
   fmt.Fscan(in, &n, &m)
   // BITs separated by parity of x and y
   bits := [2][2]*BIT2D{}
   for px := 0; px < 2; px++ {
       for py := 0; py < 2; py++ {
           bits[px][py] = NewBIT2D(n + 1)
       }
   }
   // helper: point update in difference array d
   pointUpdate := func(x, y int, v uint64) {
       if x <= 0 || y <= 0 {
           return
       }
       if x > n || y > n {
           // we use BIT size n+1, allow x==n+1
           if x > n+1 || y > n+1 {
               return
           }
       }
       px, py := x&1, y&1
       bits[px][py].update(x, y, v)
   }
   // prefix sum S(x,y)
   prefix := func(x, y int) uint64 {
       if x <= 0 || y <= 0 {
           return 0
       }
       if x > n+1 {
           x = n+1
       }
       if y > n+1 {
           y = n+1
       }
       px, py := x&1, y&1
       return bits[px][py].query(x, y)
   }
   for k := 0; k < m; k++ {
       var typ int
       fmt.Fscan(in, &typ)
       if typ == 1 {
           var x0, y0, x1, y1 int
           fmt.Fscan(in, &x0, &y0, &x1, &y1)
           res := prefix(x1, y1) ^ prefix(x0-1, y1) ^ prefix(x1, y0-1) ^ prefix(x0-1, y0-1)
           fmt.Fprintln(out, res)
       } else {
           var x0, y0, x1, y1 int
           var v uint64
           fmt.Fscan(in, &x0, &y0, &x1, &y1, &v)
           // difference array corner updates
           pointUpdate(x0, y0, v)
           pointUpdate(x0, y1+1, v)
           pointUpdate(x1+1, y0, v)
           pointUpdate(x1+1, y1+1, v)
       }
   }
}
