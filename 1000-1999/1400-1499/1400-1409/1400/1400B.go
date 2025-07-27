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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var p, f int64
       var cntS, cntW int64
       var s, w int64
       fmt.Fscan(reader, &p, &f)
       fmt.Fscan(reader, &cntS, &cntW)
       fmt.Fscan(reader, &s, &w)
       // ensure s is light
       if s > w {
           s, w = w, s
           cntS, cntW = cntW, cntS
       }
       // brute swords taken by main
       var best int64 = 0
       maxTakeS := cntS
       if p/s < maxTakeS {
           maxTakeS = p / s
       }
       for takeS := int64(0); takeS <= maxTakeS; takeS++ {
           // swords taken by main
           remP := p - takeS*s
           // axes taken by main
           takeW := remP / w
           if takeW > cntW {
               takeW = cntW
           }
           // remaining items
           leftS := cntS - takeS
           leftW := cntW - takeW
           // follower takes
           // swords
           takeSf := leftS
           if f/s < takeSf {
               takeSf = f / s
           }
           remF := f - takeSf*s
           // axes
           takeWf := remF / w
           if takeWf > leftW {
               takeWf = leftW
           }
           total := takeS + takeW + takeSf + takeWf
           if total > best {
               best = total
           }
       }
       fmt.Fprintln(writer, best)
   }
}
