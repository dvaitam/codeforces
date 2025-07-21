package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func countSameFirstLast(n uint64) uint64 {
   if n == 0 {
       return 0
   }
   s := strconv.FormatUint(n, 10)
   length := len(s)
   if length == 1 {
       return n
   }
   // precompute powers of 10 up to length-2
   pow10 := make([]uint64, length)
   pow10[0] = 1
   for i := 1; i < length; i++ {
       pow10[i] = pow10[i-1] * 10
   }
   var ans uint64
   // lengths shorter than current
   for k := 1; k < length; k++ {
       if k == 1 {
           ans += 9
       } else {
           ans += 9 * pow10[k-2]
       }
   }
   // numbers of same length
   firstDigit := uint64(s[0] - '0')
   // count those with first digit less than firstDigit
   if length == 1 {
       ans += firstDigit
   } else {
       for d := uint64(1); d < firstDigit; d++ {
           ans += pow10[length-2]
       }
       // for first digit equal to firstDigit
       // get middle part
       var mid uint64
       if length > 2 {
           m, _ := strconv.ParseUint(s[1:len(s)-1], 10, 64)
           mid = m
       }
       ans += mid
       lastDigit := uint64(s[len(s)-1] - '0')
       if lastDigit >= firstDigit {
           ans++
       }
   }
   return ans
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var l, r uint64
   if _, err := fmt.Fscan(reader, &l, &r); err != nil {
       return
   }
   var res uint64
   if l > 1 {
       res = countSameFirstLast(r) - countSameFirstLast(l-1)
   } else {
       res = countSameFirstLast(r)
   }
   fmt.Println(res)
}
