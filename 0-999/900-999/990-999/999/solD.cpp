#include<iostream> 
#include<cstdio>
#include<cmath>
#include<cstdlib>
#include<cstring>
#include<algorithm>
using namespace std;
int read()
{
	int x=0,f=1;char c=getchar();
	while (c<'0'||c>'9') {if (c=='-') f=-1;c=getchar();}
	while (c>='0'&&c<='9') x=(x<<1)+(x<<3)+(c^48),c=getchar();
	return x*f;
}
#define N 200010
int n,m,a[N],d,p[N],nxt[N],cnt[N],move[N],pre[N];
long long ans=0;
int dis(int x,int y)
{
	if (x<y) return y-x;
	else return m-x+y;
}
int find(int x){return pre[x]==x?x:pre[x]=find(pre[x]);}
int main()
{
	n=read(),m=read();
	for (int i=1;i<=n;i++)
	{
		a[i]=read();
		cnt[a[i]%m]++;
		nxt[i]=p[a[i]%m],p[a[i]%m]=i;
	}
	d=n/m;
	int tmp;for (tmp=0;tmp<m;tmp++) if (cnt[tmp]>d) break;
	pre[tmp]=tmp;for (int i=tmp+1;i<m;i++) if (cnt[i]>d) pre[i]=i;else pre[i]=pre[i-1];
	if (tmp) 
	{
		if (cnt[0]>d) pre[0]=0;else pre[0]=pre[m-1];
		for (int i=1;i<tmp;i++) if (cnt[i]>d) pre[i]=i;else pre[i]=pre[i-1];
	}
	for (int i=0;i<m;i++)
		for (int j=find(i);cnt[i]<d;j=find(j))
		{
			for (int k=p[j];cnt[j]>d&&cnt[i]<d;k=p[j]=nxt[k])
			cnt[j]--,cnt[i]++,a[k]+=dis(j,i),ans+=dis(j,i);
			if (cnt[j]==d) pre[find(j)]=find(j?j-1:m-1);
		}
	cout<<ans<<endl;
	for (int i=1;i<=n;i++) printf("%d ",a[i]);
	return 0;
}