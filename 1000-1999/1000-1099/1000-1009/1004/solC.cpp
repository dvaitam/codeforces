#include<algorithm>
#include<iostream>
#include<cstring>
#include<cstdio>
#include<cmath>
#include<queue>
#define ri register int

using namespace std;

inline char gch()
{
	static char buf[100010], *h = buf, *t = buf;
	return h == t && (t = (h = buf) + fread(buf, 1, 100000, stdin), h == t) ? EOF : *h ++;
}

inline void re(int &x)
{
	x = 0;
	char a; bool b = 0;
	while(!isdigit(a = getchar()))
		b = a == '-';
	while(isdigit(a))
		x = (x << 1) + (x << 3) + a - '0', a = getchar();
	if(b == 1)
		x *= -1;
}

int n, a[100010], num[100010];

long long ans;

bool flag[100010], fl[100010];

int main()
{
	re(n);
	for(ri i = 1; i <= n; i ++)
		re(a[i]);
	for(ri i = n; i >= 1; i --)
		num[i] = ((fl[a[i]] == 0) ? 1 : 0) + num[i + 1], fl[a[i]] = 1;
	memset(fl, 0, sizeof(fl));
	for(ri i = 1; i < n; i ++)
		if(fl[a[i]] == 0)
			ans += num[i + 1], fl[a[i]] = 1;
	printf("%I64d", ans);
}