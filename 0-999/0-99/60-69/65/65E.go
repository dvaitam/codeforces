package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

var (
   n, m int
   lnk, nxt []int
   son []int
   Free [][]int
   res [][]int
   root []int
   flag []bool
   now, c1, c2, masterkill int
   rv []int
   cnt [][2]int
   dnt []int
   reader *bufio.Reader
   writer *bufio.Writer
)

func readInt() int {
   var x int
   fmt.Fscan(reader, &x)
   return x
}

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func dfs(u int) {
   flag[u] = true
   for v := son[u]; v != 0; v = nxt[v] {
       if lnk[v] > 0 {
           lnk[v^1] = -lnk[v^1]
           if !flag[lnk[v]] {
               dfs(lnk[v])
               c2++
           } else {
               lnk[v], lnk[v^1] = 0, 0
               Free[u] = append(Free[u], v/2)
               root[len(root)-1] = u
               c1++
           }
       }
   }
}

func efs(u int) {
   dnt[u] = 1
   for v := son[u]; v != 0; v = nxt[v] {
       if lnk[v] > 0 {
           efs(lnk[v])
           dnt[u] += dnt[lnk[v]]
           if dnt[lnk[v]] == 1 {
               masterkill = u
           }
       }
   }
}

func goDFS(u int) {
   for v := son[u]; v != 0; v = nxt[v] {
       if lnk[v] != 0 {
           if lnk[v] < 0 {
               lnk[v] = -lnk[v]
           }
           lnk[v^1] = 0

           rv = append(rv, lnk[v])
           goDFS(lnk[v])
           rv = append(rv, u)
           if now < len(root) {
               res = append(res, copySlice(rv))
               tmp := make([]int, 3)
               tmp[0] = v/2
               tmp[1] = u
               tmp[2] = root[now]
               res = append(res, tmp)
               rv = rv[:1]
               rv[0] = root[now]
               goDFS(root[now])
               rv = append(rv, u)
           }
       }
   }
   for i := 0; i < len(Free[u]); i++ {
       if now < len(root) {
           res = append(res, copySlice(rv))
           tmp := make([]int, 3)
           tmp[0] = Free[u][i]
           tmp[1] = u
           tmp[2] = root[now]
           res = append(res, tmp)
           rv = rv[:1]
           rv[0] = root[now]
           goDFS(root[now])
           rv = append(rv, u)
       } else {
           break
       }
   }
}

func copySlice(a []int) []int {
   b := make([]int, len(a))
   copy(b, a)
   return b
}

func main() {
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &n, &m)
   size := 2*m + 5
   lnk = make([]int, size)
   nxt = make([]int, size)
   son = make([]int, n+2)
   Free = make([][]int, n+2)
   cnt = make([][2]int, n+2)
   dnt = make([]int, n+2)
   flag = make([]bool, n+2)

   root = make([]int, 0, n)
   res = make([][]int, 0)
   rv = make([]int, 0, n)

   j := 1
   for i := 0; i < m; i++ {
       u := readInt()
       v := readInt()
       j++
       lnk[j] = v
       nxt[j] = son[u]
       son[u] = j
       j++
       lnk[j] = u
       nxt[j] = son[v]
       son[v] = j
   }
   for i := 1; i <= n; i++ {
       if !flag[i] {
           c1, c2 = 0, 0
           root = append(root, i)
           dfs(i)
           cnt[i][0], cnt[i][1] = c1, c2
       }
   }
   // ensure start at 1
   if len(root) > 0 {
       root[0] = 1
   }
   sort.Slice(root, func(i, j int) bool {
       if i == 0 {
           return true
       }
       if j == 0 {
           return false
       }
       ai, aj := root[i], root[j]
       if cnt[ai][0] != cnt[aj][0] {
           return cnt[ai][0] > cnt[aj][0]
       }
       return cnt[ai][1] > cnt[aj][1]
   })
   for i := 1; i < len(root); i++ {
       masterkill = root[i]
       efs(root[i])
       if len(Free[root[i]]) == 0 {
           root[i] = masterkill
       }
   }
   now = 0
   for now < len(root) {
       if now > 0 {
           k := -1
           u := root[now]
           if len(Free[u]) > 0 {
               k = Free[u][len(Free[u])-1]
               Free[u] = Free[u][:len(Free[u])-1]
           } else {
               for e := son[u]; e != 0; e = nxt[e] {
                   v2 := lnk[e]
                   if v2 > 0 && dnt[v2] == 1 {
                       k = e/2
                       root = append(root, abs(lnk[e]))
                       lnk[e], lnk[e^1] = 0, 0
                       break
                   }
               }
           }
           if k == -1 {
               break
           }
           tmp := make([]int, 3)
           tmp[0], tmp[1], tmp[2] = k, 1, root[now]
           res = append(res, tmp)
           rv = rv[:0]
       }
       rv = append(rv, root[now])
       now++
       u := root[now-1]
       goDFS(u)
       if rv[len(rv)-1] != 1 {
           rv = append(rv, 1)
       }
       res = append(res, copySlice(rv))
   }
   if now < len(root) {
       fmt.Fprintln(writer, "NO")
       return
   }
   fmt.Fprintln(writer, "YES")
   fmt.Fprintln(writer, len(res)/2)
   for i := 0; i < len(res); i++ {
       if i%2 == 0 {
           fmt.Fprint(writer, len(res[i]))
           for _, v := range res[i] {
               fmt.Fprint(writer, " ", v)
           }
           fmt.Fprintln(writer)
       } else {
           a := res[i]
           fmt.Fprintln(writer, a[0], a[1], a[2])
       }
   }
}
