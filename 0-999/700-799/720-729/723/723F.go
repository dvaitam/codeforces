package main

import (
   "bufio"
   "fmt"
   "os"
)

type Pair struct{ u, v int }

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func readInt() int {
   var c byte
   var x int
   var neg bool
   for {
       b, err := reader.ReadByte()
       if err != nil {
           return x
       }
       c = b
       if c == '-' || (c >= '0' && c <= '9') {
           break
       }
   }
   if c == '-' {
       neg = true
       b, _ := reader.ReadByte()
       c = b
   }
   for c >= '0' && c <= '9' {
       x = x*10 + int(c-'0')
       b, err := reader.ReadByte()
       if err != nil {
           break
       }
       c = b
   }
   if neg {
       return -x
   }
   return x
}

func main() {
   defer writer.Flush()
   n := readInt()
   m := readInt()
   edges := make([]Pair, m)
   for i := 0; i < m; i++ {
       u := readInt()
       v := readInt()
       if u > v {
           u, v = v, u
       }
       edges[i] = Pair{u, v}
   }
   s := readInt()
   t := readInt()
   ds := readInt()
   dt := readInt()
   if s > t {
       s, t = t, s
       ds, dt = dt, ds
   }
   parent := make([]int, n+1)
   for i := 1; i <= n; i++ {
       parent[i] = i
   }
   var find func(int) int
   find = func(u int) int {
       if parent[u] != u {
           parent[u] = find(parent[u])
       }
       return parent[u]
   }
   ans := make([]Pair, 0, n)
   // initial union excluding s and t
   for _, e := range edges {
       u, v := e.u, e.v
       if u != s && v != s && u != t && v != t {
           ru := find(u)
           rv := find(v)
           if ru != rv {
               parent[ru] = rv
               ans = append(ans, Pair{u, v})
           }
       }
   }
   a := make([]int, n+1)
   b := make([]int, n+1)
   flag := false
   for _, e := range edges {
       u, v := e.u, e.v
       if u == s && v == t {
           flag = true
       } else {
           if u == s && v != t {
               a[find(v)] = v
           }
           if v == s && u != t {
               a[find(u)] = u
           }
           if u == t && v != s {
               b[find(v)] = v
           }
           if v == t && u != s {
               b[find(u)] = u
           }
       }
   }
   x := make([]Pair, 0)
   y := make([]Pair, 0)
   z := make([]Pair, 0)
   for i := 1; i <= n; i++ {
       ri := find(i)
       if ri == i && i != s && i != t {
           av := a[i]
           bv := b[i]
           if av != 0 && bv != 0 {
               x = append(x, Pair{av, bv})
           } else if av != 0 {
               y = append(y, Pair{s, av})
           } else if bv != 0 {
               z = append(z, Pair{bv, t})
           } else {
               fmt.Fprintln(writer, "No")
               return
           }
       }
   }
   ds -= len(y)
   dt -= len(z)
   ans = append(ans, y...)
   ans = append(ans, z...)
   if len(x) > 0 {
       ds--
       dt--
       last := x[len(x)-1]
       ans = append(ans, Pair{last.u, s})
       ans = append(ans, Pair{last.v, t})
       x = x[:len(x)-1]
       take := ds
       if take < 0 {
           take = 0
       }
       if take > len(x) {
           take = len(x)
       }
       ds -= take
       for i := 0; i < take; i++ {
           ans = append(ans, Pair{x[i].u, s})
       }
       rem := len(x) - take
       dt -= rem
       for i := take; i < len(x); i++ {
           ans = append(ans, Pair{x[i].v, t})
       }
   } else if flag {
       ans = append(ans, Pair{s, t})
       ds--
       dt--
   } else {
       fmt.Fprintln(writer, "No")
       return
   }
   if ds < 0 || dt < 0 {
       fmt.Fprintln(writer, "No")
       return
   }
   fmt.Fprintln(writer, "Yes")
   for _, p := range ans {
       fmt.Fprintln(writer, p.u, p.v)
   }
}
