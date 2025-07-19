package main

import (
   "bufio"
   "bytes"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, K int
   fmt.Fscan(reader, &n, &K)
   K--
   a := make([]int, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   inf := int(1e9 + 1)
   a[0], a[n+1] = inf, inf
   pre := make([]int, n+2)
   nxt := make([]int, n+2)
   st := make([]int, 0, n+2)
   // previous greater-or-equal
   st = append(st, 0)
   for i := 1; i <= n; i++ {
       for len(st) > 0 && a[i] >= a[st[len(st)-1]] {
           st = st[:len(st)-1]
       }
       pre[i] = st[len(st)-1]
       st = append(st, i)
   }
   // next greater
   st = st[:0]
   st = append(st, n+1)
   for i := n; i >= 1; i-- {
       for len(st) > 0 && a[i] > a[st[len(st)-1]] {
           st = st[:len(st)-1]
       }
       nxt[i] = st[len(st)-1]
       st = append(st, i)
   }

   var ans int64
   for i := 1; i <= n; i++ {
       L := i - pre[i] - 1
       R := nxt[i] - i - 1
       ans = (ans + calc(L, R, K)*int64(a[i])) % mod
   }
   if ans < 0 {
       ans += mod
   }
   buf := &bytes.Buffer{}
   buf.WriteString(fmt.Sprintf("%d", ans))
   os.Stdout.Write(buf.Bytes())
}

func calc(L, R, k int) int64 {
   a := int64(L) / int64(k)
   b := int64(R) / int64(k)
   x := int64(L) % int64(k)
   y := int64(R) % int64(k)
   ans := (a + b) % mod
   ans = (a*int64(k)%mod*b + ans) % mod
   ans = (x*b + ans) % mod
   ans = (y*a + ans) % mod
   if x == 0 || y == 0 {
       return ans
   }
   extra := x + y - int64(k) + 1
   if extra > 0 {
       ans = (ans + extra) % mod
   }
   return ans
}
