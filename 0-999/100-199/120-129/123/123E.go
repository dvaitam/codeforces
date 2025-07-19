package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func readInt() int {
   c, err := reader.ReadByte()
   if err != nil {
       return 0
   }
   for c <= ' ' {
       c, _ = reader.ReadByte()
   }
   sign := 1
   if c == '-' {
       sign = -1
       c, _ = reader.ReadByte()
   }
   val := 0
   for c >= '0' && c <= '9' {
       val = val*10 + int(c-'0')
       c, _ = reader.ReadByte()
   }
   return val * sign
}

func main() {
   defer writer.Flush()
   n := readInt()
   adj := make([][]int, n+1)
   for i := 1; i < n; i++ {
       u := readInt()
       v := readInt()
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   st := make([]int64, n+1)
   en := make([]int64, n+1)
   var sst, sen int64
   for i := 1; i <= n; i++ {
       st[i] = int64(readInt())
       en[i] = int64(readInt())
       sst += st[i]
       sen += en[i]
   }
   parent := make([]int, n+1)
   order := make([]int, 0, n)
   stack := make([]int, 0, n)
   stack = append(stack, 1)
   parent[1] = 0
   for len(stack) > 0 {
       x := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       order = append(order, x)
       for _, y := range adj[x] {
           if y != parent[x] {
               parent[y] = x
               stack = append(stack, y)
           }
       }
   }
   size := make([]int64, n+1)
   stSum := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       size[i] = 1
       stSum[i] = st[i]
   }
   var ans float64
   for i := len(order) - 1; i >= 0; i-- {
       x := order[i]
       p := parent[x]
       if p != 0 {
           size[p] += size[x]
           stSum[p] += stSum[x]
           ans += float64(stSum[x]) * float64(size[x]) * float64(en[p])
       }
       ans += float64(sst-stSum[x]) * float64(n-size[x]) * float64(en[x])
   }
   res := ans / float64(sst) / float64(sen)
   fmt.Fprintf(writer, "%.11f\n", res)
}
