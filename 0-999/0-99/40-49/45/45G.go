package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   be := make([]int, n+1)
   cnt := 0
   sum := n * (n + 1) / 2

   var isPrime = func(x int) bool {
       if x < 2 {
           return false
       }
       for i := 2; i*i <= x; i++ {
           if x%i == 0 {
               return false
           }
       }
       return true
   }

   var work = func(x int) {
       cnt++
       for i := n; i >= 1; i-- {
           if x >= i && be[i] == 0 {
               be[i] = cnt
               x -= i
           }
       }
   }

   if isPrime(sum) {
       work(sum)
   } else if sum%2 == 1 && isPrime(sum-2) {
       work(2)
       work(sum - 2)
   } else {
       if sum%2 == 1 {
           work(3)
           sum -= 3
       }
       for i := 2; i <= sum; i++ {
           if isPrime(i) && isPrime(sum-i) {
               work(i)
               work(sum - i)
               break
           }
       }
   }

   for i := 1; i <= n; i++ {
       fmt.Fprintf(out, "%d", be[i])
       if i < n {
           out.WriteByte(' ')
       }
   }
}
