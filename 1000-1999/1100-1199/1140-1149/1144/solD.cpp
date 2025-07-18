#include<bits/stdc++.h>
#define maxn 200005
#define FOR(a, b, c) for(int a=b; a<=c; a++)
#define inf 19260817
#define ll long long
using namespace std;

int n, a[maxn], cnt, top, x, rua, f[maxn], loc;

inline int read()
{
    char c=getchar();long long x=0,f=1;
    while(c<'0'||c>'9'){if(c=='-')f=-1;c=getchar();}
    while(c>='0'&&c<='9'){x=x*10+c-'0';c=getchar();}
    return x*f;
}

int main()
{
    n = read();
    FOR(i, 1, n)
    {
        x = read();
        f[i] = x;
        a[x]++;
        if(a[x] > cnt)
        {
            cnt = a[x];
            top = x;
            loc = i;
        }
    }
    rua = n - cnt;
    printf("%d\n", rua);
    if(!rua) return 0;
    FOR(i, loc, n)
    {
        if(f[i] < top)
            printf("1 %d %d\n", i, i-1);
        if(f[i] > top)
            printf("2 %d %d\n", i, i-1);
    }
    for(int i=loc; i>=1; i--)
    {
        if(f[i]==top) continue;
        if(f[i] < top)
            printf("1 %d %d\n", i, i+1);
        if(f[i] > top)
            printf("2 %d %d\n", i, i+1);
    }
    return 0;
}