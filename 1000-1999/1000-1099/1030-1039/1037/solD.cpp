#include<bits/stdc++.h>
#define maxn (500000+5)
#define rep(i,l,r) for(int i=l;i<=r;i++)
#define mod (998244353)
using namespace std;
inline int read()
{
    int x=0,f=1;char ch=getchar();
    while(ch<'0'||ch>'9'){if(ch=='-')f=-1;ch=getchar();}
    while(ch>='0'&&ch<='9'){x=10*x+ch-'0';ch=getchar();}
    return x*f;
}
long long gcd(long long x,long long y){return y?gcd(y,x%y):x;}
long long power(long x,long y)
{
	long long t=1;
	for(;y;y>>=1,x=x*x%mod)
		if(y&1)t=t*x%mod;
	return t;
}
struct edge
{
	int go,next;
}e[maxn];
int n,head[maxn],tot,dep[maxn],a[maxn],mark[maxn],fa[maxn];
bool v[maxn];
void insert(int x,int y)
{
	e[++tot].next=head[x];head[x]=tot;e[tot].go=y;
	e[++tot].next=head[y];head[y]=tot;e[tot].go=x;
}
void dfs(int x)
{
	v[x]=1;
	for(int i=head[x],y=e[i].go;i;i=e[i].next,y=e[i].go)if(!v[y])
	{
		dep[y]=dep[x]+1;
		fa[y]=x;
		dfs(y);
	}
}
bool check()
{
	if(a[1]!=1)return 0; 
	for(int l=1,r=1;l<=n;l=r+1)
	{
		while(r<n&&dep[a[r+1]]==dep[a[l]])r++;
		if(r<n&&dep[a[r+1]]<dep[a[l]])return 0;
		for(int i=l;i<=r;i++)
		{
			mark[a[i]]=i; 
			//cout<<i<<' '<<a[i]<<' '<<mark[a[i]]<<endl;
			if(i>l&&mark[fa[a[i]]]<mark[fa[a[i-1]]])return 0;
		}
	}
	return 1;
}
int main()
{
	n=read();
	rep(i,1,n-1)insert(read(),read());
	rep(i,1,n)a[i]=read();
	dfs(a[1]);
	//rep(i,1,n)cout<<i<<' '<<a[i]<<' '<<dep[a[i]]<<' '<<mark[a[i]]<<endl;
	cout<<(check()?"Yes":"No")<<endl;
	return 0;
}
/*
6
1 2
1 3
1 4
2 5
3 6
1 2 3 4 5 6
*/