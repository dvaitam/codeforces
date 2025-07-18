#include <cstdio>
#include <cstring>
#include <algorithm>
#include <ctime>
#include <cmath>
#include <iostream>
#define maxn 500050
#define maxm 2000400
#define rep(i,n) for(int i=1;i<=n;i++)
using namespace std;
int num[maxn],bel[maxn],wz[maxn];
int a[maxn],dl[maxn];
int n,m,top,now;
int fir[maxn],en[maxm],nex[maxm],tot;
int get()
{
	char t=getchar();
	while(t<'0'||t>'9')t=getchar();
	int x=0;
	while(t>='0'&&t<='9')
	{
		x*=10;
		x+=t-'0';
		t=getchar();
	}
	return x;
}
void tjb(int x,int y)
{
	en[++tot]=y;
	nex[tot]=fir[x];
	fir[x]=tot;
}
bool cmp(int a,int b)
{
	return ((bel[a]<bel[b])||((bel[a]==bel[b])&&(a<b)));
}
int main()
{
//	freopen("E.in","r",stdin);
//	freopen("E.out","w",stdout);
	
	n=get();
	m=get();
	rep(i,m)
	{
		int x=get();
		int y=get();
		tjb(x,y);
		tjb(y,x);
	}
	rep(i,n)dl[i]=wz[i]=i;
	rep(i,n)
	{
		if(now<i)
		{
			now=i;
			++top;
		}
		int x=dl[i];
		bel[x]=top;
		++num[top];
		int td=n;
		for(int k=fir[x];k;k=nex[k])
		{
			int j=en[k];
			if(wz[j]<=now)continue;
			if(wz[j]>td)continue;
			int y=dl[td];
			swap(wz[y],wz[j]);
			swap(dl[wz[y]],dl[wz[j]]);
			--td;
		}
		now=td;
	}
	rep(i,n)a[i]=i;
	sort(a+1,a+1+n,cmp);
	printf("%d\n",top);
	now=0;
	rep(i,top)
	{
		printf("%d ",num[i]);
		rep(j,num[i]-1)
		{
			++now;
			printf("%d ",a[now]);
		}
		++now;
		printf("%d\n",a[now]);
	}
	return 0;
}