#include<cstdio>
#include<cstring>
#include<algorithm>
using namespace std;
const int maxn=150;
int ans[maxn][maxn];
int a[maxn],b[maxn],c[maxn];
int main()
{
	int n,m;
	while(~scanf("%d%d",&n,&m))
	{
		int tot=0;
		for(int i=1;i<=n;i++)
		{
			scanf("%d",&a[i]);
			tot^=a[i];
		}
		for(int i=1;i<=m;i++)
		{
			scanf("%d",&b[i]);
			tot^=b[i];
		}
		if(tot)
		printf("NO\n");
		else
		{
			printf("YES\n");
			memset(c,0,sizeof(c)); 
			for(int i=1;i<n;i++)
			{
				tot=0;
				for(int j=1;j<m;j++)
				{
					ans[i][j]=1;
					tot^=1;
					c[j]^=1;
				}
				ans[i][m]=tot^a[i];
				c[m]^=ans[i][m];
			}
			for(int i=1;i<=m;i++)
			ans[n][i]=b[i]^c[i];
			for(int i=1;i<=n;i++)
			{
				for(int j=1;j<m;j++)
				printf("%d ",ans[i][j]);
				printf("%d\n",ans[i][m]); 
			}
		}
	}
}