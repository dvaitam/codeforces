package main

import (
   "bufio"
   "os"
   "strconv"
)

const M = 15000016

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   n := readInt(reader)
   a := make([]int, n)
   gcd := 0
   for i := 0; i < n; i++ {
       a[i] = readInt(reader)
       if i == 0 {
           gcd = a[i]
       } else {
           gcd = gcdFunc(gcd, a[i])
       }
   }
   for i := 0; i < n; i++ {
       a[i] /= gcd
   }

   // compute smallest prime factor (spf) via linear sieve
   spf := make([]int32, M)
   primes := make([]int32, 0, 1000000)
   for i := 2; i < M; i++ {
       if spf[i] == 0 {
           spf[i] = int32(i)
           primes = append(primes, int32(i))
       }
       for _, p := range primes {
           ip := int64(i) * int64(p)
           if ip >= M || p > spf[i] {
               break
           }
           spf[ip] = p
       }
   }

   cnt := make([]int32, M)
   ans := 0
   for _, v := range a {
       t := v
       for t > 1 {
           p := int(spf[t])
           for t%p == 0 {
               t /= p
           }
           cnt[p]++
           if int(cnt[p]) > ans {
               ans = int(cnt[p])
           }
       }
   }
   if ans == 0 {
       writer.WriteString("-1")
   } else {
       writer.WriteString(strconv.Itoa(n - ans))
   }
}

func gcdFunc(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func readInt(r *bufio.Reader) int {
   x := 0
   neg := false
   var ch byte
   for {
       b, err := r.ReadByte()
       if err != nil {
           return x
       }
       ch = b
       if (ch >= '0' && ch <= '9') || ch == '-' {
           break
       }
   }
   if ch == '-' {
       neg = true
   } else {
       x = int(ch - '0')
   }
   for {
       b, err := r.ReadByte()
       if err != nil {
           break
       }
       ch = b
       if ch < '0' || ch > '9' {
           break
       }
       x = x*10 + int(ch-'0')
   }
   if neg {
       x = -x
   }
   return x
}
