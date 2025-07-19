package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 998244353
const N = 300010

var two [N]int64

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   two[0] = 1
   for i := 1; i < N; i++ {
       two[i] = two[i-1] * 2 % mod
   }

   t := int(readInt(in))
   for ; t > 0; t-- {
       n := int(readInt(in))
       m := int(readInt(in))
       g := make([][]int, n+1)
       for i := 0; i < m; i++ {
           u := int(readInt(in))
           v := int(readInt(in))
           g[u] = append(g[u], v)
           g[v] = append(g[v], u)
       }
       col := make([]int8, n+1)
       var ans int64 = 1
       flag := false
       for i := 1; i <= n; i++ {
           if col[i] == 0 {
               // BFS for component
               var b, w int
               queue := []int{i}
               col[i] = 1
               b = 1
               for head := 0; head < len(queue); head++ {
                   u := queue[head]
                   for _, v := range g[u] {
                       if col[v] == 0 {
                           col[v] = -col[u]
                           if col[v] == 1 {
                               b++
                           } else {
                               w++
                           }
                           queue = append(queue, v)
                       } else if col[v] == col[u] {
                           flag = true
                           break
                       }
                   }
                   if flag {
                       break
                   }
               }
               if flag {
                   break
               }
               ans = ans * ((two[w] + two[b]) % mod) % mod
           }
       }
       if flag {
           ans = 0
       }
       out.WriteString(fmt.Sprintf("%d\n", ans))
   }
}

// readInt reads an integer from bufio.Reader
func readInt(in *bufio.Reader) int64 {
   var x int64
   var c byte
   var neg bool
   c, _ = in.ReadByte()
   for (c < '0' || c > '9') && c != '-' {
       c, _ = in.ReadByte()
   }
   if c == '-' {
       neg = true
       c, _ = in.ReadByte()
   }
   for c >= '0' && c <= '9' {
       x = x*10 + int64(c-'0')
       c, _ = in.ReadByte()
   }
   if neg {
       x = -x
   }
   return x
}
