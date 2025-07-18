#include <bits/stdc++.h>
#include<cstdio>
#include<vector>
#include<queue>
#include<ctime>
#include<algorithm>
#include<cstdlib>
#include<stack>
#include<cstring>
#include<cmath>
using namespace std;

typedef long long LL;

const int INF = 1000000000;

LL l,r;

inline LL getint()
{
    LL ret = 0,f = 1;
    char c = getchar();
    while (c < '0' || c > '9')
    {
        if (c == '-') f = -1;
        c = getchar();
    }
    while (c >= '0' && c <= '9') ret = ret * 10 + c - '0',c = getchar();
    return ret * f;
}

int main()
{
    #ifdef AMC
        freopen("AMC1.txt","r",stdin);
        freopen("AMC2.txt","w",stdout);
    #endif
    puts("YES");
    l = getint(); r = getint();
    for (LL i = l; i <= r; i += 2)
        printf("%lld %lld\n",i,i + 1);
    return 0;
}