package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func power(a, b, mod int) int {
   res := 1 % mod
   a %= mod
   for b > 0 {
       if b&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       b >>= 1
   }
   return res
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   var n, k int
   fmt.Fscan(reader, &n, &k)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       a[i]--
   }
   vis := make([]bool, n)
   ans := make([]int, n)
   f := make([][][]int, n+1)
   for i := 0; i < n; i++ {
       if vis[i] {
           continue
       }
       j := i
       var t []int
       for !vis[j] {
           t = append(t, j)
           vis[j] = true
           j = a[j]
       }
       sz := len(t)
       f[sz] = append(f[sz], t)
   }
   // validate
   for i := 1; i <= n; i++ {
       if i%2 == 0 && len(f[i]) > 0 {
           if k > 20 || len(f[i])%(1<<k) != 0 {
               fmt.Fprintln(writer, "NO")
               return
           }
       }
   }
   fmt.Fprintln(writer, "YES")
   // build answer
   for i := 1; i <= n; i++ {
       for len(f[i]) > 0 {
           cnt := len(f[i])
           t := min(k, bits.Len(uint(cnt))-1)
           pw := power(2, k-t, i)
           take := 1 << t
           tmp := f[i][cnt-take:]
           f[i] = f[i][:cnt-take]
           // rotate individual cycles
           for idx := range tmp {
               g := tmp[idx]
               ng := make([]int, i)
               for j := 0; j < i; j++ {
                   ng[(pw*j)%i] = g[j]
               }
               tmp[idx] = ng
           }
           all := i << t
           // link cycles
           for j := 0; j < all; j++ {
               c1 := j & (take - 1)
               p1 := j / take
               nj := (j + 1) % all
               c2 := nj & (take - 1)
               p2 := nj / take
               ans[tmp[c1][p1]] = tmp[c2][p2]
           }
       }
   }
   // output
   for i := 0; i < n; i++ {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprint(ans[i] + 1))
   }
   writer.WriteByte('\n')
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var T int
   fmt.Fscan(reader, &T)
   for T > 0 {
       solve(reader, writer)
       T--
   }
}
