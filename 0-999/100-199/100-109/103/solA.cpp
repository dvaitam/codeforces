#include <bits/stdc++.h>

using namespace std;

#include <stdio.h>
#include <stdlib.h>

const int N = 105;

int n;
int64_t a[N];

int main() {    
    int i;
    int64_t ans = 0;
    scanf("%d", &n);
    for (i = 1; i <= n; i++)
        scanf("%lld", a + i);
    ans = n;
    for (i = 1; i <= n; i++)
        ans += (a[i] - 1) * i;
    printf("%lld\n", ans);
    return 0;
}