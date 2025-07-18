#include <stdio.h>
#include <bits/stdc++.h>

using namespace std;
    int n, a, b, c, t;

    int gcd(int a, int b) {
   if (b == 0)
   return a;
   return gcd(b, a % b);
}
int main()
{
    scanf("%d", &n);
    for(int i=0; i<n; i++)
    {
        scanf("%d%d", &a, &b);
        if(gcd(a, b)!=1) printf("infinite\n");
        else printf("finite\n");
    }
    return 0;
}