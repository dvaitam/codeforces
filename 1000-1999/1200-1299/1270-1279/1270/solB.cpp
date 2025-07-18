#pragma GCC optimize("O3")
#include <bits/stdc++.h>
#include <tr1/unordered_map>
using namespace std;
#define ll long long
#define fi first
#define se second
#define re register
#define pb push_back
const int N=1e6+10;
const int mod=1e9+7;
void read(int &a)
{
    a=0;int d=1;char ch;
    while(ch=getchar(),ch>'9'||ch<'0')
        if(ch=='-')
            d=-1;
    a=ch^48;
    while(ch=getchar(),ch>='0'&&ch<='9')
        a=(a<<3)+(a<<1)+(ch^48);
    a*=d;
}
int a[N];
int main()
{
    int T;
    read(T);
    while(T--)
    {
        int n;
        read(n);
        for(re int i=1;i<=n;i++) read(a[i]);
        bool f=0;
        for(re int i=2;i<=n;i++)
        {
            if(abs(a[i]-a[i-1])>=2)
            {
                puts("YES");
                printf("%d %d\n",i-1,i);
                f=1;
                break;
            }
        }
        if(!f) puts("NO");
    }
    return 0;
}