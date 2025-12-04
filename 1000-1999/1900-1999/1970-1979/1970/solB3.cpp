#include <bits/stdc++.h>
#define pii pair<int,int>
#define s1 first
#define s2 second
#define N 200011
using namespace std;
int n,x[N],y[N],to[N];
pii a[N];
int main()
{
	scanf("%d",&n);
	for(int i=1;i<=n;++i)scanf("%d",&a[i].s1),a[i].s2=i;
	sort(a+1,a+1+n);
	bool flg=0;
	if(a[1].s1)
	{
		flg=1;
		for(int i=2;i<=n;++i)if(a[i].s1==a[i-1].s1){swap(a[1],a[i-1]);swap(a[2],a[i]);flg=0;break;}
	}
	x[a[1].s2]=1;y[a[1].s2]=1;
	for(int i=2;i<=n;++i)
	{
		int v=a[i].s1,u=a[i].s2;
		if(!v)x[u]=i,y[u]=1,to[u]=a[i].s2;
		else if(v<=i-1)
		{
			x[u]=i,y[u]=y[a[i-v].s2];
			o[u]=a[i-v].s2;
		}
		else
		{
			x[u]=i;y[u]=1+v-(i-1);
			o[u]=a[1].s2;
		}
	}
	if(flg)
	{
		if(n==2)printf("NO\n"),exit(0);
		y[a[2].s2]=1;
		o[a[1].s2]=a[2].s2;to[a[2].s2]=a[3].s2;to[a[3].s2]=a[1].s2;
	}
	else if(!a[1].s1)to[a[1].s2]=a[1].s2;
	else to[a[1].s2]=a[2].s2;
	printf("YES\n");
	for(int i=1;i<=n;++i)printf("%d %d\n",x[i],y[i]);
	for(int i=1;i<=n;++i)printf("%d ",to[i]);putchar(10);
	fclose(stdin);fclose(stdout);return 0;
}