package main

import (
   "bufio"
   "fmt"
   "os"
)

type pair struct{ first, second int64 }

func f(n, o int64) []pair {
   if n <= 0 {
       return nil
   }
   if n == 1 {
       var res []pair
       for n != o {
           res = append(res, pair{0, n})
           res = append(res, pair{n, n})
           n *= 2
       }
       return res
   }
   if n == 2 && o == 2 {
       return []pair{{0, 1}, {1, 1}}
   }
   if (n == 3 && o == 4) || (n == 4 && o == 4) {
       return []pair{{1, 3}, {2, 2}, {0, 4}}
   }
   if n == 5 && o == 8 {
       return []pair{{3, 5}, {2, 2}, {4, 4}, {0, 1}, {1, 1}, {0, 2}, {2, 2}, {0, 4}, {4, 4}, {0, 8}, {0, 8}}
   }
   if n == 6 && o == 8 {
       return []pair{{3, 5}, {2, 6}, {4, 4}, {0, 1}, {1, 1}, {2, 2}, {0, 4}, {4, 4}, {0, 8}, {0, 8}}
   }
   if n == o {
       return f(n-1, o)
   }
   if o >= 2*n {
       res := f(n, o/2)
       nn := n
       for nn >= 4 {
           res = append(res, pair{o / 2, o / 2})
           res = append(res, pair{0, o})
           nn -= 2
       }
       if nn == 3 {
           res = append(res, pair{o / 2, o / 2})
           res = append(res, pair{0, o / 2})
           res = append(res, pair{o / 2, o / 2})
           res = append(res, pair{0, o})
       } else if nn == 2 {
           res = append(res, pair{o / 2, o / 2})
           res = append(res, pair{0, o})
       }
       return res
   }
   // o < 2*n
   var res []pair
   for i := o - n; i < o/2; i++ {
       res = append(res, pair{i, o - i})
   }
   res1 := f(o-n-1, o)
   res2 := f(n-o/2, o/2)
   for i := range res2 {
       res2[i].first *= 2
       res2[i].second *= 2
   }
   if n-o/2 > 2 {
       // drop last
       res2 = res2[:len(res2)-1]
       res = append(res, res2...)
       res = append(res, res1...)
       res = append(res, pair{0, o / 2}, pair{o / 2, o / 2}, pair{0, o})
       return res
   } else if o-n-1 > 2 {
       res1 = res1[:len(res1)-1]
       res = append(res, res1...)
       res = append(res, res2...)
       res = append(res, pair{0, o / 2}, pair{o / 2, o / 2}, pair{0, o})
       return res
   }
   panic("unreachable")
}

func nextPow2(n int64) int64 {
   if n&(n-1) == 0 {
       return n
   }
   o := int64(1)
   for o < n {
       o <<= 1
   }
   return o
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   fmt.Fscan(in, &t)
   for ; t > 0; t-- {
       var n int64
       fmt.Fscan(in, &n)
       if n == 2 {
           fmt.Fprintln(out, -1)
           continue
       }
       o := nextPow2(n)
       res := f(n, o)
       fmt.Fprintln(out, len(res))
       for _, p := range res {
           fmt.Fprintln(out, p.first, p.second)
       }
   }
}
