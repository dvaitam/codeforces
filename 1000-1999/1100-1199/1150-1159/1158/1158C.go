package main

import (
   "bufio"
   "os"
   "strconv"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func readInt() int {
   sign := 1
   b, _ := reader.ReadByte()
   for (b < '0' || b > '9') && b != '-' {
       b, _ = reader.ReadByte()
   }
   if b == '-' {
       sign = -1
       b, _ = reader.ReadByte()
   }
   x := 0
   for b >= '0' && b <= '9' {
       x = x*10 + int(b-'0')
       b, _ = reader.ReadByte()
   }
   return x * sign
}

func main() {
   defer writer.Flush()
   cas := readInt()
   for cas > 0 {
       cas--
       n := readInt()
       a := make([]int, n+2)
       V := make([][]int, n+2)
       initDeg := make([]int, n+2)
       vis := make([]bool, n+2)
       stk := make([]int, n+2)
       Q := make([]int, n+2)
       ans := make([]int, n+2)
       for i := 1; i <= n; i++ {
           a[i] = readInt()
           if a[i] == -1 {
               a[i] = i + 1
           }
       }
       // initialization
       stk[0] = n + 1
       top := 0
       for i := 1; i <= n+1; i++ {
           vis[i] = false
           initDeg[i] = 0
       }
       vis[n+1] = true
       flg := false
       // build graph
       for i := n; i >= 1; i-- {
           if !vis[a[i]] {
               flg = true
               break
           }
           for stk[top] != a[i] {
               vis[stk[top]] = false
               V[i] = append(V[i], stk[top])
               initDeg[stk[top]]++
               top--
           }
           V[stk[top]] = append(V[stk[top]], i)
           initDeg[i]++
           top++
           stk[top] = i
           vis[i] = true
       }
       if flg {
           writer.WriteString("-1\n")
           continue
       }
       // topological sort
       H, T := 1, 0
       for i := 1; i <= n+1; i++ {
           if initDeg[i] == 0 {
               T++
               Q[T] = i
           }
       }
       for H <= T {
           u := Q[H]
           H++
           for _, v := range V[u] {
               initDeg[v]--
               if initDeg[v] == 0 {
                   T++
                   Q[T] = v
               }
           }
       }
       if T != n+1 {
           writer.WriteString("-1\n")
           continue
       }
       // assign answers
       for i := 2; i <= n+1; i++ {
           ans[Q[i]] = n - i + 2
       }
       // output
       for i := 1; i <= n; i++ {
           writer.WriteString(strconv.Itoa(ans[i]))
           if i < n {
               writer.WriteByte(' ')
           }
       }
       writer.WriteByte('\n')
   }
}
