package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const (
   MaxOut = 100000
)

var (
   rem int
   res []int64
   children [][]int
   dpLayers [][]int
   maxLayer int
   idx1 int
   vals []int64
)

func gen(vIdx int, depth int64) {
   if rem <= 0 {
       return
   }
   // handle v==1
   if vIdx == idx1 {
       res = append(res, 1)
       rem--
       return
   }
   // handle prime-like chain: children [1, v]
   ch := children[vIdx]
   if len(ch) == 2 && ch[0] == idx1 && ch[1] == vIdx {
       // sequence is [1 repeated depth times, v]
       // output first rem elements
       if depth >= int64(rem) {
           for i := 0; i < rem; i++ {
               res = append(res, 1)
           }
           rem = 0
       } else {
           // depth < rem
           for i := int64(0); i < depth; i++ {
               res = append(res, 1)
           }
           rem -= int(depth)
           if rem > 0 {
               res = append(res, vals[vIdx])
               rem--
           }
       }
       return
   }
   // normal case
   if depth == 0 {
       res = append(res, vals[vIdx])
       rem--
       return
   }
   // dive into children in order
   for _, u := range ch {
       if rem <= 0 {
           return
       }
       gen(u, depth-1)
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var X int64
   var k int64
   if _, err := fmt.Fscan(in, &X, &k); err != nil {
       return
   }
   // divisors of X
   vals = make([]int64, 0)
   for i := int64(1); i*i <= X; i++ {
       if X%i == 0 {
           vals = append(vals, i)
           if i*i != X {
               vals = append(vals, X/i)
           }
       }
   }
   sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
   n := len(vals)
   idx := make(map[int64]int, n)
   for i, v := range vals {
       idx[v] = i
   }
   idx1 = idx[1]
   // build children lists (divisors) in increasing order
   children = make([][]int, n)
   for i, v := range vals {
       for j, u := range vals {
           if u > v {
               break
           }
           if v%u == 0 {
               children[i] = append(children[i], j)
           }
       }
   }
   // dp layers
   // dpLayers[d][i] = s(vals[i], d) capped
   dpLayers = [][]int{}
   prev := make([]int, n)
   for i := range prev {
       prev[i] = 1
   }
   dpLayers = append(dpLayers, prev)
   INF := MaxOut + 5
   maxLayer = 0
   for d := 1; d <= int(k); d++ {
       cur := make([]int, n)
       for i := 0; i < n; i++ {
           sum := 0
           for _, j := range children[i] {
               sum += prev[j]
               if sum > INF {
                   sum = INF
                   break
               }
           }
           cur[i] = sum
       }
       dpLayers = append(dpLayers, cur)
       maxLayer = d
       if cur[idx[X]] >= INF {
           break
       }
   }
   if int64(maxLayer) > k {
       maxLayer = int(k)
   }
   // generate
   rem = MaxOut
   res = make([]int64, 0, MaxOut)
   root := idx[X]
   gen(root, k)
   // output
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i, v := range res {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprintf(out, "%d", v)
   }
   out.WriteByte('\n')
}
