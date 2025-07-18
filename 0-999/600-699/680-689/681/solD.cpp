#include<cmath>
#include<ctime>
#include<cstdio>
#include<cstdlib>
#include<cstring>
#include<complex>
#include<map>
#include<set>
#include<stack>
#include<queue>
#include<vector>
#include<iostream>
#include<algorithm>
#define X first
#define Y second
typedef long long ll;
typedef unsigned long long ull;
typedef std::pair<int,int> pii;
template<typename T>inline T abs(T a){ return a<0?-a:a; }
template<typename T>inline T min(T a,T b){ return a<b?a:b; }
template<typename T>inline T max(T a,T b){ return a>b?a:b; }
template<typename T>inline void Swap(T&a,T&b){ T t=a;a=b,b=t; }
template<typename T>inline bool umin(T&mn,T u){ return u<mn?mn=u,true:false; }
template<typename T>inline bool umax(T&mx,T u){ return u>mx?mx=u,true:false; }
template<typename T>inline void read(T&r)
{
	char c;r=0;bool flag=false;
	do c=getchar();while(c!='-'&&(c<'0'||c>'9'));
	if(c=='-')flag=true,c=getchar();
	do r=r*10+c-'0',
	   c=getchar();while(c>='0'&&c<='9');
	if(flag)r*=-1;
}
template<typename T1,typename T2>
inline void read(T1&r,T2&s){ read(r),read(s); }
template<typename T1,typename T2,typename T3>
inline void read(T1&r,T2&s,T3&t){ read(r),read(s),read(t); }
template<typename T1,typename T2,typename T3,typename T4>
inline void read(T1&r,T2&s,T3&t,T4&u){ read(r),read(s),read(t),read(u); }
template<typename T>inline void write(T w)
{
	if(w==0)putchar('0');else{
	if(w<0)putchar('-'),w*=-1;
	static int s[21],top;top=0;
	while(w)s[top++]=w%10,w/=10;
	while(top--)putchar(s[top]+'0');
	}
}
/**************************模板****************************
srt:2016 6 15 1:31
end:2016 6 15 1:47
nme:Gifts by the List
src:cf round 357 div2
agr:
smr:
**********************************************************/

const int N=100001;
struct ac{ int to,nx; }e[N<<1];
int first[N],L=1,fa[N];
inline void addedge(int u,int v)
{e[L].to=v,e[L].nx=first[u],first[u]=L++;}

bool list[N];
int q[N],head,tail;
int gift[N],st[N],top;
void Bfs(int u)
{
	q[tail++]=u;int j;
	while(head<tail)
	{
		u=q[head++];
		for(j=first[u];j;j=e[j].nx)
			q[tail++]=e[j].to;
	}
}
int main()
{
	int n,m,u,v,i;read(n,m);
	while(m--)read(u,v),addedge(u,v),fa[v]=u;
	for(i=1;i<=n;i++)
	{
		read(gift[i]);
		list[gift[i]]=true;
		if(!fa[i])Bfs(i);
	}
	while(tail--)
	{
		i=q[tail];
		if(gift[i]!=i&&gift[fa[i]]!=gift[i])
			return puts("-1"),0;
		if(list[i])st[top++]=i;
	}
	printf("%d\n",top);
	for(i=0;i<top;i++)
		printf("%d\n",st[i]);
}