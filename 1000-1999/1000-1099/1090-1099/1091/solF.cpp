#include <bits/stdc++.h>

using namespace std;

#define LL long long
#define MOD 998244353
#define N 100007

LL a[N];
char s[N];

int main()
{
	int n;
	scanf("%d",&n);
    a[0] = 1LL;
	for(int i=0;i<n;i++)
    {
        scanf("%lld", &a[i]);
        a[i]*=2;
    }
    scanf(" %s", s);
    LL ans = 0LL;
    LL now = 6LL;
    LL nowW = 0LL;
    LL nowG = 0LL;
    for(int i=0;i<n;i++)
    {
        if(s[i]=='W')
        {
            now = 4LL;
            nowW += a[i]/2;
            ans += a[i]*2;
        }
        if(s[i]=='G')
        {
            LL tmp = min(nowW, a[i]/2);
            nowW -= tmp;
            a[i] -= tmp*2;
            ans += tmp*4;
            nowG += tmp*2;
            nowG += a[i]/2;
            ans += a[i]*3;
        }
        if(s[i]=='L')
        {
            LL tmp = min(nowW, a[i]/2);
            nowW -= tmp;
            a[i] -= tmp*2;
            ans += tmp*4;
            tmp = min(nowG, a[i]/2);
            nowG -= tmp;
            a[i] -= tmp*2;
            ans += tmp*6;
            ans += a[i] * now;
        }
    }
    printf("%lld\n", ans/2);
	return 0;
}