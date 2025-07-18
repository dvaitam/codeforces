#include<cstdio>

#include<algorithm>

#include<iostream>

#include<queue>

#include<cstring>

#include<vector>

using namespace std;

int n,m,k,cnt=1;

int mx,sum;

int a[110];

int main()

{

	scanf("%d%d%d",&n,&m,&k);

	a[1]=1,sum=1,mx=1;

	if(m==0&&k==0){

		for(int i=1;i<=n;i++)printf("1 ");printf("\n");return 0;

	}

	if(m+k+1>=n&&k==0){

		printf("-1\n");return 0;

	}

	if(m+k+1<n)

	a[2]=1,sum=2,cnt=2;

	while(k--)a[++cnt]=sum+1,mx=max(mx,a[cnt]),sum=sum+a[cnt];

	while(m--)a[++cnt]=mx+1,mx=max(mx,a[cnt]);

	for(int i=1;i<=n;i++)if(a[i]>50000){printf("-1\n");return 0;}

	for(int i=1;i<=n;i++)if(a[i])printf("%d ",a[i]);else printf("1 ");

	return 0;

}