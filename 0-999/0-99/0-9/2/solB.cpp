#include<cstdio>

#include<iostream>

using namespace std;

#define INF 0x3f3f3f3f

int in()

{

	char ch=getchar();

	while (ch<'0'||ch>'9')

		ch=getchar();

	int ret=0;

	while (ch>='0'&&ch<='9')

		ret=ret*10+ch-'0',ch=getchar();

	return ret;

}

int n,tmp,ans=INF,cnt;

int a[1010][1010][2],f[1010][1010][2],g[1010][1010];

char pl[2010];

bool flag;

void getans(int x,int y)

{

	if (x==1&&y==1)

		return;

	if (g[x][y]==0)

	{

		getans(x-1,y);

		pl[cnt++]='D';

	}

	else

	{

		getans(x,y-1);

		pl[cnt++]='R';

	}

}

int main()

{

	n=in();

	for (int i=1;i<=n;i++)

		for (int j=1;j<=n;j++)

		{

			tmp=in();

			if (!tmp)

			{

				if (!flag)

				{

					cnt=0;

					for (int k=2;k<=j;k++)

						pl[cnt++]='R';

					for (int k=2;k<=i;k++)

						pl[cnt++]='D';

					for (int k=j+1;k<=n;k++)

						pl[cnt++]='R';

					for (int k=i+1;k<=n;k++)

						pl[cnt++]='D';

					ans=1;

					flag=1;

				}

				a[i][j][0]=a[i][j][1]=10000;

			}

			else

			{

				while (tmp%2==0)

					tmp/=2,a[i][j][0]++;

				while (tmp%5==0)

					tmp/=5,a[i][j][1]++;

			}

		}

	for (int i=0;i<=n;i++)

		f[i][0][0]=f[i][0][1]=f[0][i][0]=f[0][i][1]=INF;

	f[1][0][0]=f[1][0][1]=0;	

	for (int i=1;i<=n;i++)

		for (int j=1;j<=n;j++)

			if (f[i-1][j][0]<f[i][j-1][0]||f[i-1][j][0]==f[i][j-1][0]&&f[i-1][j][1]<f[i][j-1][1])

			{

				f[i][j][0]=f[i-1][j][0]+a[i][j][0],f[i][j][1]=f[i-1][j][1]+a[i][j][1];

				g[i][j]=0;

			}

			else

			{

				f[i][j][0]=f[i][j-1][0]+a[i][j][0],f[i][j][1]=f[i][j-1][1]+a[i][j][1];

				g[i][j]=1;

			}

	if (min(f[n][n][0],f[n][n][1])<ans)

	{

		ans=min(f[n][n][0],f[n][n][1]);

		cnt=0;

		getans(n,n);

	}

	for (int i=0;i<=n;i++)

		f[i][0][0]=f[i][0][1]=f[0][i][0]=f[0][i][1]=INF;

	f[1][0][0]=f[1][0][1]=0;	

	for (int i=1;i<=n;i++)

		for (int j=1;j<=n;j++)

			if (f[i-1][j][1]<f[i][j-1][1]||f[i-1][j][1]==f[i][j-1][1]&&f[i-1][j][0]<f[i][j-1][0])

			{

				f[i][j][0]=f[i-1][j][0]+a[i][j][0],f[i][j][1]=f[i-1][j][1]+a[i][j][1];

				g[i][j]=0;

			}

			else

			{

				f[i][j][0]=f[i][j-1][0]+a[i][j][0],f[i][j][1]=f[i][j-1][1]+a[i][j][1];

				g[i][j]=1;

			}

	if (min(f[n][n][0],f[n][n][1])<ans)

	{

		ans=min(f[n][n][0],f[n][n][1]);

		cnt=0;

		getans(n,n);

	}

	printf("%d\n%s",ans,pl);

	return 0;

}