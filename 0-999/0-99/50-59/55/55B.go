package main

import (
   "fmt"
)

// dfs recursively applies operations to minimize result.
func dfs(nums []int64, ops []string, idx int) int64 {
   if len(nums) == 1 {
       return nums[0]
   }
   minRes := int64(1<<63 - 1)
   op := ops[idx]
   n := len(nums)
   for i := 0; i < n; i++ {
       for j := i + 1; j < n; j++ {
           var t int64
           if op == "*" {
               t = nums[i] * nums[j]
           } else {
               t = nums[i] + nums[j]
           }
           next := make([]int64, 0, n-1)
           for k := 0; k < n; k++ {
               if k == i || k == j {
                   continue
               }
               next = append(next, nums[k])
           }
           next = append(next, t)
           res := dfs(next, ops, idx+1)
           if res < minRes {
               minRes = res
           }
       }
   }
   return minRes
}

func main() {
   var a, b, c, d int64
   var op1, op2, op3 string
   if _, err := fmt.Scan(&a, &b, &c, &d); err != nil {
       return
   }
   fmt.Scan(&op1, &op2, &op3)
   nums := []int64{a, b, c, d}
   ops := []string{op1, op2, op3}
   result := dfs(nums, ops, 0)
   fmt.Println(result)
}
