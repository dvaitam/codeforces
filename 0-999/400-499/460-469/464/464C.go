package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
   "strconv"
)

const MOD = 1000000007
const PHI = MOD - 1

// modPow computes a^e % MOD
func modPow(a, e int64) int64 {
   res := int64(1)
   base := a % MOD
   for e > 0 {
       if e&1 == 1 {
           res = (res * base) % MOD
       }
       base = (base * base) % MOD
       e >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read initial string s
   line, _ := reader.ReadString('\n')
   s := strings.TrimSpace(line)
   // read number of queries
   line, _ = reader.ReadString('\n')
   n, _ := strconv.Atoi(strings.TrimSpace(line))

   // read all queries
   ds := make([]int, n)
   ts := make([]string, n)
   for i := 0; i < n; i++ {
       line, _ = reader.ReadString('\n')
       line = strings.TrimSpace(line)
       parts := strings.SplitN(line, "->", 2)
       d := parts[0]
       ds[i] = int(d[0] - '0')
       ts[i] = parts[1]
   }

   // initialize value and length for digits 0..9
   val := make([]int64, 10)
   length := make([]int64, 10)
   for d := 0; d < 10; d++ {
       val[d] = int64(d)
       length[d] = 1
   }

   // process queries in reverse
   for i := n - 1; i >= 0; i-- {
       d := ds[i]
       t := ts[i]
       var newVal int64 = 0
       var newLen int64 = 0
       // build for digit d by concatenating each character in t
       for _, ch := range t {
           cd := int(ch - '0')
           // append current mapping of cd
           // newVal = newVal * 10^{length[cd]} + val[cd]
           pow := modPow(10, length[cd])
           newVal = (newVal*pow + val[cd]) % MOD
           newLen = (newLen + length[cd]) % PHI
       }
       val[d] = newVal
       length[d] = newLen
   }

   // compute result for initial string
   var result int64 = 0
   for _, ch := range s {
       d := int(ch - '0')
       pow := modPow(10, length[d])
       result = (result*pow + val[d]) % MOD
   }
   fmt.Println(result)
}
