#include<stdio.h>
#include<string.h>
#define MAXN 1005
int vis[MAXN];
int main()
{
	int M,n,a,b,i;
	while(scanf("%d%d",&M,&n)!=EOF)
	{
		memset(vis,0,sizeof(vis));
		for(int i=0;i<n;i++) 
		{
			scanf("%d%d",&a,&b);
			vis[a]++,vis[b]++;
		}
		for(i=1;i<=M;i++)
			if(vis[i]==0)
				break;
		printf("%d\n",M-1);
		for(int  j=1;j<=M;j++)
		{
			if(j!=i)
				printf("%d %d\n",i,j);
		}
	}
	return 0;
}