package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = 1000000000

type Pair struct {
   d, node int
}

func bfs(s int, sm [][]int, ds []int) {
   n := len(sm)
   for i := 0; i < n; i++ {
       ds[i] = INF
   }
   ds[s] = 0
   q := make([]int, 0, n)
   q = append(q, s)
   for head := 0; head < len(q); head++ {
       v := q[head]
       for _, to := range sm[v] {
           if ds[to] > ds[v]+1 {
               ds[to] = ds[v] + 1
               q = append(q, to)
           }
       }
   }
}

func upd(md *[3]Pair, t Pair) {
   if t.d > md[0].d {
       md[2] = md[1]
       md[1] = md[0]
       md[0] = t
   } else if t.d > md[1].d {
       md[2] = md[1]
       md[1] = t
   } else if t.d > md[2].d {
       md[2] = t
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   sm := make([][]int, n)
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       a--
       b--
       if a == b {
           continue
       }
       sm[a] = append(sm[a], b)
   }
   for i := 0; i < n; i++ {
       // remove duplicates
       mapi := make(map[int]struct{}, len(sm[i]))
       uniq := sm[i][:0]
       for _, v := range sm[i] {
           if _, ok := mapi[v]; !ok {
               mapi[v] = struct{}{}
               uniq = append(uniq, v)
           }
       }
       sm[i] = uniq
   }

   dout := make([][]int, n)
   for i := 0; i < n; i++ {
       dout[i] = make([]int, n)
   }
   mdout := make([][3]Pair, n)
   mdin := make([][3]Pair, n)
   for i := 0; i < n; i++ {
       for j := 0; j < 3; j++ {
           mdout[i][j] = Pair{-1, -1}
           mdin[i][j] = Pair{-1, -1}
       }
   }

   ds := make([]int, n)
   for i := 0; i < n; i++ {
       bfs(i, sm, ds)
       copy(dout[i], ds)
       for j := 0; j < n; j++ {
           if j != i && dout[i][j] < INF {
               upd(&mdout[i], Pair{dout[i][j], j})
               upd(&mdin[j], Pair{dout[i][j], i})
           }
       }
   }

   ans := -1
   res := [4]int{0, 0, 0, 0}
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if i == j || dout[i][j] >= INF {
               continue
           }
           for p0 := 0; p0 < 3; p0++ {
               in := mdin[i][p0]
               if in.d < 0 || in.node == j {
                   continue
               }
               for p1 := 0; p1 < 3; p1++ {
                   out := mdout[j][p1]
                   if out.d < 0 || out.node == in.node || out.node == i {
                       continue
                   }
                   cur := dout[i][j] + in.d + out.d
                   if cur > ans {
                       ans = cur
                       res[0] = in.node
                       res[1] = i
                       res[2] = j
                       res[3] = out.node
                   }
               }
           }
       }
   }
   // output 1-based
   fmt.Fprintf(writer, "%d %d %d %d\n", res[0]+1, res[1]+1, res[2]+1, res[3]+1)
}
