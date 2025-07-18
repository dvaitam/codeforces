#include<bits/stdc++.h>
using namespace std;
void read(int &num)
{
	char ch; bool flag = 0;
	while(!isdigit(ch=getchar()))if(ch=='-')flag=!flag;
	for(num=0;isdigit(ch);num=num*10+ch-'0',ch=getchar());
	if(flag) num=-num;
}

const int MAXN = 200005;

int n, fir[MAXN], to[MAXN*2], nxt[MAXN*2], cnt;

inline void add(int u, int v)
{
	to[++cnt] = v; nxt[cnt] = fir[u]; fir[u] = cnt;
}

int f[MAXN][3];

void dp(int u, int fa)
{
	int v, mn = 0x7f7f7f7f;
	f[u][0] = f[u][1] = f[u][2] = 0;
	for(int i = fir[u]; i; i = nxt[i]) if((v=to[i]) != fa)
	{
		dp(v, u);
		f[u][0] += min(min(f[v][0], f[v][1]), f[v][2]);
		f[u][1] += min(f[v][0], f[v][2]);
		f[u][2] += min(f[v][0], f[v][2]);
		mn = min(mn, f[v][0] - min(f[v][0], f[v][2]));
	}
	if(fa != 1) f[u][0]++;
	f[u][2] += mn;
}

int main()
{
	int x, y;
	read(n);
	for(int i = 1; i < n; i++)
		read(x), read(y), add(x, y), add(y, x);
	int Ans = 0;
	for(int i = fir[1]; i; i = nxt[i])
	{
		dp(to[i], 1);
		for(int j = fir[to[i]]; j; j = nxt[j]) if(to[j] != to[i])
			Ans += min(min(f[to[j]][0], f[to[j]][1]), f[to[j]][2]);
	}
	printf("%d\n", Ans);
}