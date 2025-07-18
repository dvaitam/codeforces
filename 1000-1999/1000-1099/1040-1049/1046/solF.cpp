#include<cstdio>
#include<iostream>
#include<cstring>
#include<algorithm>

#define ll long long
#define ri register int

using namespace std;

inline void read(int &x)
{
	x = 0;
	bool f = 0;
	char ch = getchar();
	while(!isdigit(ch)) f = ch == '-', ch = getchar();
	while(isdigit(ch)) x = (x<<3)+(x<<1)+ch-48, ch = getchar();
	if(f) x = -x;
}

const int MAXN = 200050;
int n, a[MAXN], x, k;
ll ans;

int main()
{
	read(n);
	for(ri i = 1; i <= n; i++)
		read(a[i]);
	read(x), read(k);
	for(ri i = 1; i <= n; i++)
	{
		ans += k*(a[i]/(x+k));
		a[i] %= (x+k);
		if(a[i] > x) ans += k;
	}
		
	printf("%I64d", ans);
	return 0;
}