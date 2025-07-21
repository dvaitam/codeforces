package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var d int
   if _, err := fmt.Fscan(reader, &d); err != nil {
       return
   }
   const N = 10000000
   isPrime := make([]bool, N+1)
   for i := 2; i <= N; i++ {
       isPrime[i] = true
   }
   for i := 2; i*i <= N; i++ {
       if isPrime[i] {
           for j := i * i; j <= N; j += i {
               isPrime[j] = false
           }
       }
   }
   count := 0
   for i := 2; i <= N; i++ {
       if !isPrime[i] {
           continue
       }
       ri := reverse(i)
       if ri != i && ri <= N && isPrime[ri] {
           count++
           if count == d {
               fmt.Println(i)
               return
           }
       }
   }
}

func reverse(x int) int {
   rev := 0
   for x > 0 {
       rev = rev*10 + x%10
       x /= 10
   }
   return rev
}
