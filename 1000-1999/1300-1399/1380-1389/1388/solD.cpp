#include<bits/stdc++.h>
using namespace std;
template<typename T>
inline void read(T &x)
{
	x=0;
	int f=1;
	char ch=getchar();
	while(ch<'0'||ch>'9')
	{
		if(ch=='-') f=-1;
		ch=getchar();
}
	while(ch>='0'&&ch<='9') x=x*10+(ch^48),ch=getchar();
	x*=f;
	return;
}	
template<typename T>
void write(T x)
{
	if(x<0) putchar('-'),x=-x;
	if(x>=10) write(x/10);
	putchar(x%10+'0');
	return;
}
const int MAXN=200010;
int n;
int a[MAXN],b[MAXN];
int ind[MAXN];
queue<int> q;
int tot=1;
int edge[MAXN],nxt[MAXN],hd[MAXN];
long long f[MAXN];
int d[MAXN];
inline void add_edge(int u,int v)
{
	edge[tot]=v,nxt[tot]=hd[u],hd[u]=tot++;
	++d[v];
}
int main()
{
	read(n);
	for(int i=1;i<=n;++i) read(a[i]),f[i]=a[i];
	for(int i=1;i<=n;++i) read(b[i]);
	for(int i=1;i<=n;++i)
	if(b[i]!=-1) ++ind[b[i]];
	for(int i=1;i<=n;++i)
	if(ind[i]==0) q.push(i);
	while(!q.empty())
	{
		int x=q.front();q.pop();
		if(b[x]!=-1)
		{
			if(f[x]<0) add_edge(b[x],x);
			else f[b[x]]+=f[x],add_edge(x,b[x]);
			--ind[b[x]];
			if(ind[b[x]]==0) q.push(b[x]);
		}
	}
	long long ans=0;
	for(int i=1;i<=n;++i) ans+=f[i];
	write(ans),putchar('\n');
	for(int i=1;i<=n;++i)
	if(d[i]==0) q.push(i);
	while(!q.empty())
	{
		int x=q.front();q.pop();
		write(x),putchar(' ');
		for(int i=hd[x];i;i=nxt[i])
		{
			--d[edge[i]];
			if(d[edge[i]]==0) q.push(edge[i]);
		}
	}
	putchar('\n');
	return 0;
}