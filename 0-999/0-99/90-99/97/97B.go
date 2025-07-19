package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// P represents a pair of integers
type P struct { x, y int }

var v []P

// solve recursively combines pairs in v between indices b and e
func solve(b, e int) {
   if e-b <= 1 {
       return
   }
   m := (b + e) / 2
   solve(b, m)
   solve(m, e)
   // combine x from middle with all y in [b,e)
   for i := b; i < e; i++ {
       v = append(v, P{v[m].x, v[i].y})
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   v = make([]P, 0, n)
   for i := 0; i < n; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       v = append(v, P{x, y})
   }
   // initial sort
   sort.Slice(v, func(i, j int) bool {
       if v[i].x != v[j].x {
           return v[i].x < v[j].x
       }
       return v[i].y < v[j].y
   })
   // recursive combination
   solve(0, n)
   // final sort
   sort.Slice(v, func(i, j int) bool {
       if v[i].x != v[j].x {
           return v[i].x < v[j].x
       }
       return v[i].y < v[j].y
   })
   // unique
   u := make([]P, 0, len(v))
   for i, p := range v {
       if i == 0 || p != v[i-1] {
           u = append(u, p)
       }
   }
   // output
   fmt.Fprintln(writer, len(u))
   for _, p := range u {
       fmt.Fprintln(writer, p.x, p.y)
   }
}
