package main

import (
   "bufio"
   "fmt"
   "os"
)

type Item struct {
   w, s, v int
}

var (
   S int
   seg [][]Item
   f   []int
   changes [][2]int
   ans int
)

func add(u, L, R, l, r int, it Item) {
   if l >= R || r <= L {
       return
   }
   if l <= L && R <= r {
       seg[u] = append(seg[u], it)
       return
   }
   m := (L + R) >> 1
   add(u<<1, L, m, l, r, it)
   add(u<<1|1, m, R, l, r, it)
}

func dfs(u, L, R int) {
   base := len(changes)
   // apply items
   for _, it := range seg[u] {
       w := it.w
       s := it.s
       v := it.v
       for cw := S - w; cw >= 0; cw-- {
           if f[cw] >= 0 && cw <= s {
               nw := cw + w
               nv := f[cw] + v
               if nv > f[nw] {
                   changes = append(changes, [2]int{nw, f[nw]})
                   f[nw] = nv
               }
           }
       }
   }
   if R - L == 1 {
       // leaf: update ans
       for cw := 0; cw <= S; cw++ {
           if f[cw] > ans {
               ans = f[cw]
           }
       }
   } else {
       m := (L + R) >> 1
       dfs(u<<1, L, m)
       dfs(u<<1|1, m, R)
   }
   // rollback
   for len(changes) > base {
       c := changes[len(changes)-1]
       changes = changes[:len(changes)-1]
       f[c[0]] = c[1]
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   // read n, S
   fmt.Fscan(reader, &n, &S)
   items := make([]struct{l, r, w, s, v int}, n)
   T := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &items[i].l, &items[i].r, &items[i].w, &items[i].s, &items[i].v)
       if items[i].r > T {
           T = items[i].r
       }
   }
   // build segment tree
   size := 1
   for size < T {
       size <<= 1
   }
   seg = make([][]Item, size<<1)
   for i := 0; i < n; i++ {
       it := Item{items[i].w, items[i].s, items[i].v}
       add(1, 0, size, items[i].l, items[i].r, it)
   }
   // init dp
   f = make([]int, S+1)
   for i := range f {
       f[i] = -1
   }
   f[0] = 0
   // dfs
   dfs(1, 0, size)
   fmt.Fprintln(writer, ans)
}
