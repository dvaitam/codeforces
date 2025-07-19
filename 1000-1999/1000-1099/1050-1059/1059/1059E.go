package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var L, S int64
   if _, err := fmt.Fscan(reader, &n, &L, &S); err != nil {
       return
   }
   w := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &w[i])
       if w[i] > S {
           fmt.Fprintln(writer, -1)
           return
       }
   }
   fa := make([]int, n+1)
   son := make([]int, n+1)
   for i := 2; i <= n; i++ {
       fmt.Fscan(reader, &fa[i])
       son[fa[i]]++
   }
   dep := make([]int, n+1)
   sum := make([]int64, n+1)
   dep[1] = 1
   sum[1] = w[1]
   for i := 2; i <= n; i++ {
       dep[i] = dep[fa[i]] + 1
       sum[i] = sum[fa[i]] + w[i]
   }
   // DSU parent
   id := make([]int, n+1)
   for i := 1; i <= n; i++ {
       id[i] = i
   }
   find := func(x int) int {
       for id[x] != x {
           id[x] = id[id[x]]
           x = id[x]
       }
       return x
   }
   // queue of leaves
   q := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if son[i] == 0 {
           q = append(q, i)
       }
   }
   v := make([]bool, n+1)
   var ans int
   head := 0
   for head < len(q) {
       x := q[head]
       head++
       if v[x] {
           continue
       }
       ans++
       var tmp int64
       var dist int64
       cur := x
       for cur != 0 && tmp <= S {
           xx := find(cur)
           dist += int64(dep[cur] - dep[xx])
           tmp += sum[cur] - sum[xx]
           if dist+1 > L || tmp+w[cur] > S {
               break
           }
           // cover xx
           id[xx] = find(fa[xx])
           dist++
           tmp += w[xx]
           v[xx] = true
           son[fa[xx]]--
           if fa[xx] != 0 && son[fa[xx]] == 0 && !v[fa[xx]] {
               q = append(q, fa[xx])
           }
           cur = fa[xx]
       }
   }
   fmt.Fprintln(writer, ans)
}
