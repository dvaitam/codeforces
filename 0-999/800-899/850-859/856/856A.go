package main

import (
   "bufio"
   "fmt"
   "os"
)

const N = 1 << 20

var isPrime [N]bool
var primes []int

func initPrimes() {
   for i := 2; i < N; i++ {
       isPrime[i] = true
   }
   for p := 2; p < 10000; p++ {
       if isPrime[p] {
           primes = append(primes, p)
           for x := p * p; x < N; x += p {
               isPrime[x] = false
           }
       }
   }
}

func main() {
   initPrimes()
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var tc int
   fmt.Fscan(reader, &tc)
   diff := make([]bool, N)
   for tc > 0 {
       tc--
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       // reset diff
       for i := range diff {
           diff[i] = false
       }
       // compute differences
       for i := 0; i < n; i++ {
           for j := 0; j < i; j++ {
               d := a[i] - a[j]
               if d < 0 {
                   d = -d
               }
               diff[d] = true
           }
       }
       // find smallest prime step
       step := 0
       for _, p := range primes {
           ok := true
           for x := p; x < N; x += p {
               if diff[x] {
                   ok = false
                   break
               }
           }
           if ok {
               step = p
               break
           }
       }
       if step > 0 {
           writer.WriteString("YES\n")
           for i := 0; i < n; i++ {
               if i > 0 {
                   writer.WriteByte(' ')
               }
               fmt.Fprintf(writer, "%d", i*step+1)
           }
           writer.WriteByte('\n')
       } else {
           writer.WriteString("NO\n")
       }
   }
}
