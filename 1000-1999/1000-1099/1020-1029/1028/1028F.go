package main

import (
   "bufio"
   "fmt"
   "os"
)

// pt represents a point with integer coordinates.
type pt struct { x, y int64 }

var (
   // m groups points by their squared distance from origin.
   m   map[int64][]pt
   // z counts, for each direction key, how many points already have a reflection partner on that line.
   z   map[int64]int64
   // num is the total number of points in the set.
   num int64
)

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

// sim computes a unique key for the direction of vector (x, y) reduced by its gcd.
func sim(x, y int64) int64 {
   g := gcd(x, y)
   x /= g
   y /= g
   // combine into one key: x * 1e6 + y (since coordinates <=112904, 1e6 is safe multiplier)
   return x*1000000 + y
}

// upd adds (t==1) or removes (t==-1) point x in the group v and updates global counters.
func upd(v *[]pt, x pt, t int64) {
   num += t
   key := sim(x.x, x.y)
   if t == 1 {
       z[key]++
       for _, a := range *v {
           z[sim(a.x+x.x, a.y+x.y)] += 2
       }
       *v = append(*v, x)
   } else {
       // removal
       z[key]--
       // remove one occurrence of x from v (order does not matter)
       for i, a := range *v {
           if a.x == x.x && a.y == x.y {
               last := len(*v) - 1
               (*v)[i] = (*v)[last]
               *v = (*v)[:last]
               break
           }
       }
       for _, a := range *v {
           z[sim(a.x+x.x, a.y+x.y)] -= 2
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var q int
   fmt.Fscan(reader, &q)
   m = make(map[int64][]pt, q)
   z = make(map[int64]int64, q)
   num = 0
   for i := 0; i < q; i++ {
       var t, x, y int64
       fmt.Fscan(reader, &t, &x, &y)
       d := x*x + y*y
       switch t {
       case 1:
           v := m[d]
           upd(&v, pt{x, y}, 1)
           m[d] = v
       case 2:
           v := m[d]
           upd(&v, pt{x, y}, -1)
           m[d] = v
       case 3:
           key := sim(x, y)
           res := num - z[key]
           fmt.Fprintln(writer, res)
       }
   }
