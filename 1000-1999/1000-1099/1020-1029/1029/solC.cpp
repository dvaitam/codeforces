#pragma G++ optimize (2)
#include <iostream>
#include <cstdio>
#include <cmath>
#include <cstring>
#include <cstdlib>
#include <algorithm>
#include <string>
#define INF 0x3f3f3f3f
#define NO 300005
#define MO 100005
typedef long long ll;
//by Oliver
using namespace std;
ll read()
{
	char ch = ' ', last;
	ll ans = 0;
	while (ch < '0' || ch > '9')
		last = ch, ch = getchar();
	while (ch >= '0' && ch <= '9')
		ans = ans * 10 + int(ch - '0'), ch = getchar();
	if (last == '-')
		return -ans;
	return ans;
}
void write(ll x)
{
	if (x >= 10)
		write(x / 10);
	putchar(x % 10 + '0');
}
//head

int n, l[NO], r[NO];
//variable

int ints(int l, int r)
{
	if (r - l <= 0)
		return 0;
	return r - l;
}
void init()
{
	n = read();
	for (int i = 1; i <= n; i++)
		l[i] = read(), r[i] = read();
}
//functions

int main()
{
	init();
	int L = -INF, R = INF, posL, posR;
	for (int i = 1; i <= n; i++)
	{
		if (l[i] > L)
			L = l[i], posL = i;
		if (r[i] < R)
			R = r[i], posR = i;
	}
	int ans = ints(L, R);
	L = -INF, R = INF;
	for (int i = 1; i <= n; i++)
	{
		if (i == posL)
			continue;
		if (l[i] > L)
			L = l[i];
		if (r[i] < R)
			R = r[i];
	}
	ans = max(ans, ints(L, R));
	L = -INF, R = INF;
	for(int i = 1; i <= n; i++)
	{
		if (i == posR)
			continue;
		if (l[i] > L)
			L = l[i];
		if (r[i] < R)
			R = r[i];
	}
	ans = max(ans, ints(L, R));
	cout << ans << endl;
	return 0;
}
//main