package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   v   [][]int
   res []int
   used []bool
   par  []int
   dep  []int
   cc   []int
   vv   []int
   cycle bool
)

func dfs1(x, p int) {
   used[x] = true
   par[x] = p
   for _, y := range v[x] {
       if cycle {
           return
       }
       if used[y] && y != p {
           cycle = true
           cx := x
           for cx != y {
               cc = append(cc, cx)
               cx = par[cx]
           }
           cc = append(cc, y)
           return
       }
       if !used[y] {
           dfs1(y, x)
       }
       if cycle {
           return
       }
   }
}

func dfs2(x, p, d int) {
   used[x] = true
   par[x] = p
   dep[x] = d
   if len(v[x]) >= 3 {
       vv = append(vv, x)
   }
   for _, y := range v[x] {
       if y != p {
           dfs2(y, x, d+1)
       }
   }
}

func goVisit(x, p int) {
   vv = append(vv, x)
   for _, y := range v[x] {
       if y != p {
           goVisit(y, x)
       }
   }
}

func out(n int, w *bufio.Writer) {
   w.WriteString("YES\n")
   for i := 0; i < n; i++ {
       w.WriteString(fmt.Sprintf("%d", res[i]))
       if i+1 < n {
           w.WriteByte(' ')
       } else {
           w.WriteByte('\n')
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var tn int
   fmt.Fscan(reader, &tn)
   for tt := 0; tt < tn; tt++ {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       v = make([][]int, n)
       res = make([]int, n)
       used = make([]bool, n)
       par = make([]int, n)
       dep = make([]int, n)
       cc = cc[:0]
       cycle = false
       for i := 0; i < m; i++ {
           var x, y int
           fmt.Fscan(reader, &x, &y)
           x--
           y--
           v[x] = append(v[x], y)
           v[y] = append(v[y], x)
       }
       // find cycle
       for i := 0; i < n; i++ {
           if !used[i] && !cycle {
               dfs1(i, -1)
           }
       }
       if cycle {
           for _, x := range cc {
               res[x] = 1
           }
           out(n, writer)
           continue
       }
       // degree >=4
       found := false
       for x := 0; x < n; x++ {
           if len(v[x]) >= 4 {
               res[x] = 2
               for i := 0; i < 4; i++ {
                   res[v[x][i]] = 1
               }
               out(n, writer)
               found = true
               break
           }
       }
       if found {
           continue
       }
       // path case
       for i := 0; i < n; i++ {
           used[i] = false
       }
       for i := 0; i < n; i++ {
           if used[i] {
               continue
           }
           vv = vv[:0]
           dfs2(i, -1, 0)
           if len(vv) >= 2 {
               x := vv[0]
               y := vv[1]
               xx, yy := x, y
               for xx != yy {
                   if dep[xx] > dep[yy] {
                       res[xx] = 2
                       xx = par[xx]
                   } else if dep[xx] < dep[yy] {
                       res[yy] = 2
                       yy = par[yy]
                   } else {
                       res[xx] = 2
                       res[yy] = 2
                       xx = par[xx]
                       yy = par[yy]
                   }
               }
               res[xx] = 2
               cnt := 0
               for _, z := range v[x] {
                   if res[z] == 0 {
                       res[z] = 1
                       cnt++
                       if cnt == 2 {
                           break
                       }
                   }
               }
               cnt = 0
               for _, z := range v[y] {
                   if res[z] == 0 {
                       res[z] = 1
                       cnt++
                       if cnt == 2 {
                           break
                       }
                   }
               }
               out(n, writer)
               found = true
               break
           }
       }
       if found {
           continue
       }
       // degree ==3 case
       for x := 0; x < n; x++ {
           if len(v[x]) == 3 {
               vvv := make([][]int, 0, 3)
               for _, nei := range v[x] {
                   vv = vv[:0]
                   goVisit(nei, x)
                   // copy vv
                   tmp := make([]int, len(vv))
                   copy(tmp, vv)
                   vvv = append(vvv, tmp)
               }
               // sort by size: simple selection
               if len(vvv[0]) > len(vvv[1]) { vvv[0], vvv[1] = vvv[1], vvv[0] }
               if len(vvv[0]) > len(vvv[2]) { vvv[0], vvv[2] = vvv[2], vvv[0] }
               if len(vvv[1]) > len(vvv[2]) { vvv[1], vvv[2] = vvv[2], vvv[1] }
               // case a
               if len(vvv[0]) >= 2 && len(vvv[1]) >= 2 && len(vvv[2]) >= 2 {
                   res[x] = 3
                   for i := 0; i < 2; i++ {
                       res[vvv[0][i]] = 2 - i
                       res[vvv[1][i]] = 2 - i
                       res[vvv[2][i]] = 2 - i
                   }
                   out(n, writer)
                   found = true
                   break
               }
               // case b
               if len(vvv[1]) >= 3 && len(vvv[2]) >= 3 {
                   res[x] = 4
                   res[vvv[0][0]] = 2
                   for i := 0; i < 3; i++ {
                       res[vvv[1][i]] = 3 - i
                       res[vvv[2][i]] = 3 - i
                   }
                   out(n, writer)
                   found = true
                   break
               }
               // case c
               if len(vvv[1]) >= 2 && len(vvv[2]) >= 5 {
                   res[x] = 6
                   res[vvv[0][0]] = 3
                   res[vvv[1][0]] = 4
                   res[vvv[1][1]] = 2
                   for i := 0; i < 5; i++ {
                       res[vvv[2][i]] = 5 - i
                   }
                   out(n, writer)
                   found = true
                   break
               }
           }
       }
       if found {
           continue
       }
       // no case
       writer.WriteString("NO\n")
   }
}
