#include<cstdio>
#include<algorithm>
#include<iostream>
#include<cstring>
#include<queue>
#include<deque>
#include<set>
#include<map>
#include<cstdlib>
#include<ctime>
#define LL long long
#define inf 0x7ffffff
using namespace std;
inline LL read()
{
    LL x=0,f=1;char ch=getchar();
    while(ch<'0'||ch>'9'){if(ch=='-')f=-1;ch=getchar();}
    while(ch>='0'&&ch<='9'){x=x*10+ch-'0';ch=getchar();}
    return x*f;
}
inline void write(LL a)
{
    if (a<0){printf("-");a=-a;}
    if (a>=10)write(a/10);
    putchar(a%10+'0');
}
inline void writeln(LL a){write(a);printf("\n");}
int n,cnt,root,ans1,ans2;
int fa[1000010];
LL v[1000010];
LL sum[1000010];
struct edge{int to,next;}e[1000010];
int head[1000010];
LL tot;
inline void ins(int u,int v)
{
	e[++cnt].to=v;
	e[cnt].next=head[u];
	head[u]=cnt;
}
inline void dfs(int x)
{
	sum[x]=v[x];
	for (int i=head[x];i;i=e[i].next)
	{
		dfs(e[i].to);
		if (ans1&&ans2)return;
		if(sum[e[i].to]==tot)
		{
			if (ans1&&ans2)return;
			if (ans1)ans2=e[i].to;
			else ans1=e[i].to;
			continue;
		}
		sum[x]+=sum[e[i].to];
	}
}
int main()
{
	n=read();
	for (int i=1;i<=n;i++)
	{
		fa[i]=read();v[i]=read();tot+=v[i];
		if (fa[i])ins(fa[i],i);else root=i;
	}
	bool mmm=0;
	if(tot<0)tot=-tot,mmm=1;
	if (tot%3!=0){puts("-1");return 0;}
	tot/=3;if (mmm)tot=-tot;
	dfs(root);
	if (!ans2||!ans1)puts("-1");
	else printf("%d %d\n",ans1,ans2);
}