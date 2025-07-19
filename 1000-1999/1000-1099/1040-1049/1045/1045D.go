package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   reader = bufio.NewReader(os.Stdin)
)

// readInt reads an integer from input
func readInt() int {
   var x int
   var c byte
   // skip non-digits
   for {
       b, err := reader.ReadByte()
       if err != nil {
           return x
       }
       c = b
       if c >= '0' && c <= '9' {
           break
       }
   }
   // read digits
   for {
       if c < '0' || c > '9' {
           break
       }
       x = x*10 + int(c-'0')
       b, err := reader.ReadByte()
       if err != nil {
           break
       }
       c = b
   }
   return x
}

// readFloat reads a float with decimals, mimicking C++ style
func readFloat() float64 {
   var x float64
   var c byte
   // skip non-digits
   for {
       b, err := reader.ReadByte()
       if err != nil {
           return x
       }
       c = b
       if c >= '0' && c <= '9' {
           break
       }
   }
   // integer part
   for {
       if c < '0' || c > '9' {
           break
       }
       x = x*10 + float64(c-'0')
       b, err := reader.ReadByte()
       if err != nil {
           return x
       }
       c = b
   }
   // fractional part
   if c == '.' {
       var y int
       for {
           b, err := reader.ReadByte()
           if err != nil {
               break
           }
           c = b
           if c < '0' || c > '9' {
               break
           }
           y = y*10 + int(c-'0')
       }
       x += float64(y) * 0.01
   }
   return x
}

func main() {
   n := readInt()
   qv := make([]float64, n+1)
   p := make([]float64, n+1)
   sum := make([]float64, n+1)
   par := make([]int, n+1)
   d := make([]int, n+1)
   adj := make([][]int, n+1)
   vis := make([]bool, n+1)

   for i := 1; i <= n; i++ {
       qv[i] = readFloat()
   }
   for i := 1; i < n; i++ {
       u := readInt() + 1
       v := readInt() + 1
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // iterative DFS to set parent and child count
   type item struct{ node, idx int }
   stack := []item{{1, 0}}
   vis[1] = true
   par[1] = 0
   for len(stack) > 0 {
       cur := &stack[len(stack)-1]
       u := cur.node
       if cur.idx < len(adj[u]) {
           v := adj[u][cur.idx]
           cur.idx++
           if vis[v] {
               continue
           }
           vis[v] = true
           par[v] = u
           d[u]++
           stack = append(stack, item{v, 0})
       } else {
           stack = stack[:len(stack)-1]
       }
   }
   // initial calculations
   p[0] = 1.0
   ans := 1.0
   for i := 1; i <= n; i++ {
       sum[par[i]] += qv[i]
       ans -= qv[i] * (p[par[i]] - float64(d[i]) + sum[i])
       p[i] = qv[i]
   }
   // queries
   qc := readInt()
   for i := 0; i < qc; i++ {
       pos := readInt() + 1
       r := readFloat()
       s := r - p[pos]
       sum[par[pos]] += s
       ans -= s * (p[par[pos]] - float64(d[pos]) + sum[pos])
       p[pos] = r
       fmt.Printf("%.5f\n", ans)
   }
}
