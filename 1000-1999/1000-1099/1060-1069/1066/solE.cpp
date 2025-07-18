#include<cstdio>
#include<cstring>
#include<iostream>
#define N 200010
#define p 998244353
using namespace std;
char a[N], b[N];
int  s[N];
int main()
{
	int n, m, mul = 1, ans = 0;
	scanf("%d%d", &n, &m);
	scanf("%s%s", a + 1, b + 1);
	for(int i=1; i<=m; ++i) s[i] = s[i-1] + b[i] - '0';
	for(int i=n; i; --i)
	{
		if(a[i] == '1' && n - i < m)
		{
			int res = (long long)mul * s[m - (n - i)] % p;
			ans += res;
			if(ans >= p) ans -= p;
		}
		mul += mul;
		if(mul >= p) mul -= p;
	}
	cout << ans;
}