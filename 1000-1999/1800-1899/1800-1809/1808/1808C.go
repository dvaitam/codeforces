package main

import (
   "bufio"
   "fmt"
   "os"
)

func solve(a, b int64, out *bufio.Writer) {
   // extract digits
   var digita, digitb []int
   ta, tb := a, b
   for ta > 0 {
       digita = append(digita, int(ta%10))
       ta /= 10
   }
   for tb > 0 {
       digitb = append(digitb, int(tb%10))
       tb /= 10
   }
   na, nb := len(digita), len(digitb)
   if na < nb {
       var res int64
       for i := 0; i < na; i++ {
           res = res*10 + 9
       }
       fmt.Fprintln(out, res)
       return
   }
   // reverse to get most significant first
   for i := 0; i < na/2; i++ {
       digita[i], digita[na-1-i] = digita[na-1-i], digita[i]
   }
   for i := 0; i < nb/2; i++ {
       digitb[i], digitb[nb-1-i] = digitb[nb-1-i], digitb[i]
   }
   // try all possible [l,r] ranges of digits
   for diff := 0; diff <= 9; diff++ {
       for l := 0; l+diff <= 9; l++ {
           r := l + diff
           var front int64
           // match prefix
           for i := 0; i < na; i++ {
               if digita[i] == digitb[i] {
                   if digita[i] < l || digita[i] > r {
                       break
                   }
                   front = front*10 + int64(digita[i])
                   if i == na-1 {
                       fmt.Fprintln(out, front)
                       return
                   }
               } else {
                   // try pick a digit from a
                   if digita[i] >= l && digita[i] <= r {
                       res := front*10 + int64(digita[i])
                       for j := i + 1; j < na; j++ {
                           res = res*10 + int64(r)
                       }
                       if res >= a {
                           fmt.Fprintln(out, res)
                           return
                       }
                   }
                   // try pick a digit from b
                   if digitb[i] >= l && digitb[i] <= r {
                       res := front*10 + int64(digitb[i])
                       for j := i + 1; j < na; j++ {
                           res = res*10 + int64(l)
                       }
                       if res <= b {
                           fmt.Fprintln(out, res)
                           return
                       }
                   }
                   // try pick intermediate
                   for k := digita[i] + 1; k <= digitb[i]-1; k++ {
                       if k >= l && k <= r {
                           res := front*10 + int64(k)
                           for j := i + 1; j < na; j++ {
                               res = res*10 + int64(l)
                           }
                           fmt.Fprintln(out, res)
                           return
                       }
                   }
                   break
               }
           }
       }
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   fmt.Fscan(in, &t)
   for ; t > 0; t-- {
       var a, b int64
       fmt.Fscan(in, &a, &b)
       solve(a, b, out)
   }
}
