#include<stdio.h>
#include<string.h>
#include<algorithm>
using namespace std;
int a[1111];
int find(int x)
{
	if(a[x]!=x)
		a[x]=find(a[x]);
	return a[x];
}
int main()
{
	int n;
	while(scanf("%d",&n)!=-1)
	{
		int i,x,y,x0,y0,t=0,k=0,b[1111][3];
		for(i=1;i<=n;i++)
			a[i]=i;
		for(i=0;i<n-1;i++)
		{
			scanf("%d%d",&x,&y);
			x0=find(x);
			y0=find(y);
			if(x0==y0)
			{
				b[t][0]=x;
				b[t][1]=y;
				t++;
			}
			else
			{
				a[y0]=x0;
			}
		}
		for(i=1;i<=n;i++)
			if(a[i]==i)
				b[k++][2]=i;
			printf("%d\n",t);
			for(i=0;i<t;i++)
				printf("%d %d %d %d\n",b[i][0],b[i][1],b[i][2],b[i+1][2]);
	}
	return 0;
}