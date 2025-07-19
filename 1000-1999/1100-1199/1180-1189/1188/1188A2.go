package main

import (
   "bufio"
   "fmt"
   "os"
)

type Edge struct {
   u, v, val int
}
type Op struct {
   u, v, val int
}

var (
   n     int
   to    [][]int
   de    []int
   es    []Edge
   ops   []Op
)

func findLev(u, pre int) int {
   if de[u] == 1 {
       return u
   }
   for _, v := range to[u] {
       if v == pre {
           continue
       }
       return findLev(v, u)
   }
   return -1
}

func work(u, v int) []int {
   if de[u] == 1 {
       return []int{u}
   }
   first, second := -1, -1
   for _, x := range to[u] {
       if x == v {
           continue
       }
       if first == -1 {
           first = findLev(x, u)
       } else if second == -1 {
           second = findLev(x, u)
           break
       }
   }
   return []int{first, second}
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &n)
   to = make([][]int, n+1)
   de = make([]int, n+1)
   es = make([]Edge, 0, n-1)
   ops = make([]Op, 0, 4*(n-1))
   for i := 0; i < n-1; i++ {
       var u, v, val int
       fmt.Fscan(reader, &u, &v, &val)
       to[u] = append(to[u], v)
       to[v] = append(to[v], u)
       de[u]++
       de[v]++
       es = append(es, Edge{u, v, val})
   }
   for i := 1; i <= n; i++ {
       if de[i] == 2 {
           fmt.Fprintln(writer, "NO")
           return
       }
   }
   fmt.Fprintln(writer, "YES")
   tot := 0
   // process edges
   for _, e := range es {
       u, v, val := e.u, e.v, e.val
       ndsU := work(u, v)
       ndsV := work(v, u)
       if len(ndsU) == 2 {
           ops = append(ops, Op{ndsU[0], ndsU[1], -val / 2})
           if len(ndsV) == 2 {
               ops = append(ops, Op{ndsV[0], ndsV[1], -val / 2})
               ops = append(ops, Op{ndsU[0], ndsV[1], val / 2})
               ops = append(ops, Op{ndsU[1], ndsV[0], val / 2})
           } else {
               ops = append(ops, Op{ndsU[0], ndsV[0], val / 2})
               ops = append(ops, Op{ndsU[1], ndsV[0], val / 2})
           }
       } else {
           if len(ndsV) == 2 {
               ops = append(ops, Op{ndsV[0], ndsV[1], -val / 2})
               ops = append(ops, Op{ndsU[0], ndsV[1], val / 2})
               ops = append(ops, Op{ndsU[0], ndsV[0], val / 2})
           } else {
               ops = append(ops, Op{ndsU[0], ndsV[0], val})
           }
       }
   }
   tot = len(ops)
   fmt.Fprintln(writer, tot)
   for _, o := range ops {
       fmt.Fprintf(writer, "%d %d %d\n", o.u, o.v, o.val)
   }
}
