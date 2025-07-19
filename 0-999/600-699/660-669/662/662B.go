package main

import (
   "bufio"
   "fmt"
   "os"
)

type Edge struct {
   to    int
   color byte
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   G := make([][]Edge, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       var c string
       fmt.Fscan(reader, &u, &v, &c)
       col := c[0]
       G[u] = append(G[u], Edge{v, col})
       G[v] = append(G[v], Edge{u, col})
   }

   INF := int(1e9)
   resR, resB := 0, 0
   wrongR, wrongB := false, false
   Rcol := make([]int, n+1)
   Bcol := make([]int, n+1)
   visR := make([]bool, n+1)
   visB := make([]bool, n+1)
   var resultR, resultB []int

   // Process R-coloring
   for i := 1; i <= n; i++ {
       if !visR[i] {
           // BFS for component
           queue := []int{i}
           visR[i] = true
           Rcol[i] = 0
           comp := []int{i}
           for qi := 0; qi < len(queue); qi++ {
               v := queue[qi]
               for _, e := range G[v] {
                   if !visR[e.to] {
                       if e.color == 'R' {
                           Rcol[e.to] = Rcol[v]
                       } else {
                           Rcol[e.to] = Rcol[v] ^ 1
                       }
                       visR[e.to] = true
                       queue = append(queue, e.to)
                       comp = append(comp, e.to)
                   }
               }
           }
           // split and choose minimal
           var part0, part1 []int
           for _, v := range comp {
               if Rcol[v] == 0 {
                   part0 = append(part0, v)
               } else {
                   part1 = append(part1, v)
               }
           }
           if len(part0) < len(part1) {
               resultR = append(resultR, part0...)
               resR += len(part0)
           } else {
               resultR = append(resultR, part1...)
               resR += len(part1)
           }
       }
   }
   // Process B-coloring
   for i := 1; i <= n; i++ {
       if !visB[i] {
           queue := []int{i}
           visB[i] = true
           Bcol[i] = 0
           comp := []int{i}
           for qi := 0; qi < len(queue); qi++ {
               v := queue[qi]
               for _, e := range G[v] {
                   if !visB[e.to] {
                       if e.color == 'B' {
                           Bcol[e.to] = Bcol[v]
                       } else {
                           Bcol[e.to] = Bcol[v] ^ 1
                       }
                       visB[e.to] = true
                       queue = append(queue, e.to)
                       comp = append(comp, e.to)
                   }
               }
           }
           var part0, part1 []int
           for _, v := range comp {
               if Bcol[v] == 0 {
                   part0 = append(part0, v)
               } else {
                   part1 = append(part1, v)
               }
           }
           if len(part0) < len(part1) {
               resultB = append(resultB, part0...)
               resB += len(part0)
           } else {
               resultB = append(resultB, part1...)
               resB += len(part1)
           }
       }
   }
   // validate
   for u := 1; u <= n; u++ {
       for _, e := range G[u] {
           v := e.to
           // R constraints: R edge same color, B edge diff color
           if e.color == 'R' && Rcol[u] != Rcol[v] {
               wrongR = true
           }
           if e.color == 'B' && Rcol[u] == Rcol[v] {
               wrongR = true
           }
           // B constraints: B edge same color, R edge diff color
           if e.color == 'B' && Bcol[u] != Bcol[v] {
               wrongB = true
           }
           if e.color == 'R' && Bcol[u] == Bcol[v] {
               wrongB = true
           }
       }
   }
   // output
   if wrongR && wrongB {
       fmt.Fprint(writer, -1)
       return
   }
   if wrongB {
       resB = INF
   }
   if wrongR {
       resR = INF
   }
   if resR < resB {
       fmt.Fprintln(writer, resR)
       for i, v := range resultR {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
   } else {
       fmt.Fprintln(writer, resB)
       for i, v := range resultB {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
   }
}
