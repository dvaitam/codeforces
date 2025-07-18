#include <cmath>

#include <bitset>

#include <cstdio>

#include <cstring>

#include <algorithm>



using namespace std;



bitset<1010> a[1010],b;

int ans1[1010][2],ans2[1010][2];



int main()

{

	#ifndef ONLINE_JUDGE

		freopen("input.txt","r",stdin);

		freopen("output.txt","w",stdout);

	#endif

	int n;scanf("%d",&n);

	for (int i=1;i<=n;i++)

		a[i][i%n+1]=a[i%n+1][i]=1;

	for (int i=0;i<n-3;i++)

	{

		int x,y;scanf("%d%d",&x,&y);

		a[x][y]=a[y][x]=1;

	}

	int m1=0,m2=0;

	while (1)

	{

		bool ok=0;

		for (int i=3;i<n;i++)

			if (!a[1][i])

			{

				int k=i;while (!a[1][k]) k++;

				b=a[i-1]&a[k];int x;

				for (int j=2;j<=n;j++) if (b[j]) {x=j;break;}

				ans1[++m1][0]=i-1;ans1[m1][1]=k;

				a[1][x]=a[x][1]=1;

				a[i-1][k]=a[k][i-1]=0;

				ok=1;break;

			}

		if (!ok) break;

	}

	for (int i=1;i<=n;i++) a[i].reset();

	for (int i=1;i<=n;i++)

		a[i][i%n+1]=a[i%n+1][i]=1;

	for (int i=0;i<n-3;i++)

	{

		int x,y;scanf("%d%d",&x,&y);

		a[x][y]=a[y][x]=1;

	}

	while (1)

	{

		bool ok=0;

		for (int i=3;i<n;i++)

			if (!a[1][i])

			{

				int k=i;while (!a[1][k]) k++;

				b=a[i-1]&a[k];int x;

				for (int j=2;j<=n;j++) if (b[j]) {x=j;break;}

				ans2[++m2][0]=1;ans2[m2][1]=x;

				a[1][x]=a[x][1]=1;

				a[i-1][k]=a[k][i-1]=0;

				ok=1;break;

			}

		if (!ok) break;

	}

	printf("%d\n",m1+m2);

	for (int i=1;i<=m1;i++) printf("%d %d\n",ans1[i][0],ans1[i][1]);

	for (int i=m2;i;i--) printf("%d %d\n",ans2[i][0],ans2[i][1]);

	return 0;

}