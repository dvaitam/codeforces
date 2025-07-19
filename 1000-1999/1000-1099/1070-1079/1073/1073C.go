package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var N int
   fmt.Fscan(reader, &N)
   S := make([]byte, N)
   fmt.Fscan(reader, &S)
   var TX, TY int
   fmt.Fscan(reader, &TX, &TY)

   // early check
   need := abs(TX) + abs(TY)
   if N < need || (N-need)%2 != 0 {
       fmt.Fprintln(writer, -1)
       return
   }

   // prefix sums
   X := make([]int, N+1)
   Y := make([]int, N+1)
   for i := 0; i < N; i++ {
       X[i+1] = X[i]
       Y[i+1] = Y[i]
       switch S[i] {
       case 'U':
           Y[i+1]++
       case 'D':
           Y[i+1]--
       case 'R':
           X[i+1]++
       case 'L':
           X[i+1]--
       }
   }

   // check if between i and r segment of length L can cover remaining
   ok := func(i, r int) bool {
       x := X[i] + (X[N] - X[r])
       y := Y[i] + (Y[N] - Y[r])
       dist := abs(TX-x) + abs(TY-y)
       return dist <= r-i
   }

   ans := N + 1
   r := 0
   for i := 0; i <= N; i++ {
       if r < i {
           r = i
       }
       for r < N && !ok(i, r) {
           r++
       }
       if ok(i, r) {
           ans = min(ans, r-i)
       }
   }
   if ans > N {
       ans = -1
   }
   fmt.Fprintln(writer, ans)
}
