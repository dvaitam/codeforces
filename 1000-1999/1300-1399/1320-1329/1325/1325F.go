package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func readInt() int {
   var c byte
   var err error
   // skip non-digits
   c, err = reader.ReadByte()
   for err == nil && (c < '0' || c > '9') {
       c, err = reader.ReadByte()
   }
   if err != nil {
       return 0
   }
   n := 0
   for err == nil && c >= '0' && c <= '9' {
       n = n*10 + int(c-'0')
       c, err = reader.ReadByte()
   }
   return n
}

func main() {
   defer writer.Flush()
   n := readInt()
   m := readInt()
   b := int(math.Sqrt(float64(n-1))) + 1
   g := make([][]int, n+1)
   for i := 0; i < m; i++ {
       u := readInt()
       v := readInt()
       g[u] = append(g[u], v)
       g[v] = append(g[v], u)
   }
   dep := make([]int, n+1)
   vis := make([]bool, n+1)
   ans := make([]int, 0, b)
   st := make([]int, 0, n)
   type frame struct{ u, parent, idx int }
   stack := make([]frame, 0, n)
   // start DFS at 1
   stack = append(stack, frame{1, 0, -1})
   for len(stack) > 0 {
       fr := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       u, p, idx := fr.u, fr.parent, fr.idx
       if idx == -1 {
           dep[u] = dep[p] + 1
           st = append(st, u)
           // prepare to iterate children
           stack = append(stack, frame{u, p, 0})
       } else if idx < len(g[u]) {
           v := g[u][idx]
           // next child later
           stack = append(stack, frame{u, p, idx + 1})
           if dep[v] > 0 {
               if dep[v] <= dep[u]-b+1 {
                   // found long cycle
                   c := dep[u] - dep[v] + 1
                   fmt.Fprintln(writer, 2)
                   fmt.Fprintln(writer, c)
                   // output cycle nodes
                   // nodes from u back to just after v
                   for j := len(st) - 1; j >= dep[v]; j-- {
                       fmt.Fprint(writer, st[j], " ")
                   }
                   // add v
                   fmt.Fprintln(writer, st[dep[v]-1])
                   writer.Flush()
                   os.Exit(0)
               }
           } else {
               // visit child
               stack = append(stack, frame{v, u, -1})
           }
       } else {
           // post visit u
           if !vis[u] {
               ans = append(ans, u)
               for _, v := range g[u] {
                   vis[v] = true
               }
           }
           // pop from stack path
           if len(st) > 0 {
               st = st[:len(st)-1]
           }
       }
   }
   // no cycle: output independent set
   fmt.Fprintln(writer, 1)
   for i := 0; i < b && i < len(ans); i++ {
       fmt.Fprint(writer, ans[i], " ")
   }
   fmt.Fprintln(writer)
}
