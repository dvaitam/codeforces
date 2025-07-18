//#include<bits/stdc++.h>
#include<iostream>
#include<queue>
#include<iomanip>
#include<cctype>
#include<cstdio>
#include<deque>
#include<utility>
#include<cmath>
#include<ctime>
#include<cstring>
#include<string>
#include<cstdlib>
#include<vector>
#include<algorithm>
#include<stack>
#include<map>
#include<set>
#include<bitset>
#define INF 1000000000
#define ll long long
#define db double
#define pb push_back
#define un unsigned
using namespace std;
char *fs,*ft,buf[1<<15];
inline char getc()
{
	return (fs==ft&&(ft=(fs=buf)+fread(buf,1,1<<15,stdin),fs==ft))?0:*fs++;
}
inline int read()
{
	int x=0,f=1;char ch=getc();
	while(ch<'0'||ch>'9'){if(ch=='-')f=-1;ch=getc();}
	while(ch>='0'&&ch<='9'){x=x*10+ch-'0';ch=getc();}
	return x*f;
}
const int MAXN=200010;
int n,k;
int c[MAXN],flag,minn=INF,top,maxx,ans[MAXN];
struct wy
{
	int x,y;
	int id;
	int friend operator <(wy a,wy b)
	{
		return a.y<b.y;
	}
}t[MAXN];
priority_queue<wy> q;
inline int cmp(wy a,wy b){return a.x<b.x;}
int main()
{
	//freopen("1.in","r",stdin);
	n=read();k=read();
	for(int i=1;i<=n;++i)
	{
		int x,y;
		x=read();y=read();
		t[i]=(wy){x,y,i};
		++c[x];--c[y+1];
		minn=min(minn,x);
		maxx=max(maxx,y);
	}
	sort(t+1,t+1+n,cmp);
	for(int i=minn;i<=maxx;++i)
	{
		c[i]+=c[i-1];
		while(t[flag].x<=i&&flag<=n){q.push(t[flag]);++flag;}
		while(c[i]>k)
		{
			wy w=q.top();
			q.pop();
			ans[++top]=w.id;
			--c[i];++c[w.y+1];
		}
	}
	printf("%d\n",top);
	sort(ans+1,ans+1+top);
	for(int i=1;i<=top;++i)printf("%d ",ans[i]);
	return 0;
}