package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const N = 100500
const INF int64 = 100000000000000

type rat struct {
   en, de int
}

func gcd(a, b int) int {
   for b != 0 {
       a %= b
       a, b = b, a
   }
   return a
}

func newRat(en, de int) rat {
   g := gcd(en, de)
   return rat{en / g, de / g}
}

func ratLess(a, b rat) bool {
   if a.de != b.de {
       return a.de < b.de
   }
   return a.en < b.en
}

func rev(x int) int {
   if x == 0 {
       return 0
   }
   r := 0
   for x > 0 {
       r = r*10 + x%10
       x /= 10
   }
   return r
}

type pair struct {
   r rat
   x int
}

var S []pair

func count() {
   S = make([]pair, N)
   for i := 1; i <= N; i++ {
       en := i
       de := rev(i)
       S[i-1] = pair{newRat(en, de), i}
   }
   sort.Slice(S, func(i, j int) bool {
       if S[i].r.de != S[j].r.de {
           return S[i].r.de < S[j].r.de
       }
       if S[i].r.en != S[j].r.en {
           return S[i].r.en < S[j].r.en
       }
       return S[i].x < S[j].x
   })
}

func col(x, mx int) int {
   if x == 0 {
       return 0
   }
   rv := rev(x)
   r := newRat(rv, x)
   n := len(S)
   lo := sort.Search(n, func(i int) bool {
       return !ratLess(S[i].r, r)
   })
   hi := sort.Search(n, func(i int) bool {
       if ratLess(r, S[i].r) {
           return true
       }
       if ratLess(S[i].r, r) {
           return false
       }
       return S[i].x > mx
   })
   return hi - lo
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var mxx, mxy, w int
   if _, err := fmt.Fscan(in, &mxx, &mxy, &w); err != nil {
       return
   }
   count()
   x := mxx + 1
   y := 0
   cur := int64(0)
   bestArea := INF
   bestX, bestY := -1, -1
   for x > 0 {
       cur -= int64(col(x, y))
       x--
       for cur < int64(w) && y <= mxy {
           y++
           cur += int64(col(y, x))
       }
       if y > mxy {
           break
       }
       area := int64(x) * int64(y)
       if area < bestArea {
           bestArea = area
           bestX = x
           bestY = y
       }
   }
   if bestArea >= INF {
       fmt.Fprint(out, -1)
   } else {
       fmt.Fprintf(out, "%d %d", bestX, bestY)
   }
}
