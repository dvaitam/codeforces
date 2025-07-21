package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var k int64
   fmt.Fscan(reader, &n, &k)
   li := make([]int64, n)
   ri := make([]int64, n)
   Lmin := int64(1<<62)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &li[i], &ri[i])
       if ri[i]-li[i] < Lmin {
           Lmin = ri[i] - li[i]
       }
   }
   // generate lucky numbers up to 1e18
   luckies := genLuckies()
   sort.Slice(luckies, func(i, j int) bool { return luckies[i] < luckies[j] })

   // sort segment endpoints
   lis := make([]int64, n)
   ris := make([]int64, n)
   copy(lis, li)
   copy(ris, ri)
   sort.Slice(lis, func(i, j int) bool { return lis[i] < lis[j] })
   sort.Slice(ris, func(i, j int) bool { return ris[i] < ris[j] })

   // two pointers for segments
   p1, p2 := 0, 0
   c1, c2 := 0, 0
   sumRi := big.NewInt(0)
   sumLi := big.NewInt(0)
   // initial X
   X := luckies[0]
   // initial p2 and sumLi for li > X
   for p2 < n && lis[p2] <= X {
       p2++
   }
   c2 = n - p2
   for i := p2; i < n; i++ {
       sumLi.Add(sumLi, big.NewInt(lis[i]))
   }

   best := 0
   l := 0
   bigX := big.NewInt(X)
   bigY := big.NewInt(0)
   cost1 := big.NewInt(0)
   cost2 := big.NewInt(0)
   tmp := big.NewInt(0)
   sumCost := big.NewInt(0)
   bigK := big.NewInt(k)

   for r := 0; r < len(luckies); r++ {
       Y := luckies[r]
       bigY.SetInt64(Y)
       // include new S1 segments with ri < Y
       for p1 < n && ris[p1] < Y {
           sumRi.Add(sumRi, big.NewInt(ris[p1]))
           p1++
           c1++
       }
       // adjust l to satisfy window length and cost
       for l <= r {
           // current X and Y
           X = luckies[l]
           // window length check
           if Y-X > Lmin {
               // move l
           } else {
               // compute cost1 = c1*Y - sumRi
               tmp.SetInt64(int64(c1))
               tmp.Mul(tmp, bigY)
               cost1.Sub(tmp, sumRi)
               // compute cost2 = sumLi - c2*X
               tmp.SetInt64(int64(c2))
               tmp.Mul(tmp, bigX)
               cost2.Sub(sumLi, tmp)
               // total cost
               sumCost.Add(cost1, cost2)
               if sumCost.Cmp(bigK) <= 0 {
                   break
               }
           }
           // increment l
           oldX := X
           l++
           if l >= len(luckies) {
               break
           }
           X = luckies[l]
           bigX.SetInt64(X)
           // remove from S2 those with li <= X
           for p2 < n && lis[p2] <= X {
               sumLi.Sub(sumLi, big.NewInt(lis[p2]))
               p2++
               c2--
           }
       }
       // update best
       if l <= r {
           cnt := r - l + 1
           if cnt > best {
               best = cnt
           }
       }
   }
   fmt.Fprint(writer, best)
}

// generate all lucky numbers up to 1e18
func genLuckies() []int64 {
   var res []int64
   // BFS
   queue := []int64{4, 7}
   for i := 0; i < len(queue); i++ {
       v := queue[i]
       res = append(res, v)
       if v <= 1000000000000000000/10 {
           queue = append(queue, v*10+4, v*10+7)
       }
   }
   return res
}
