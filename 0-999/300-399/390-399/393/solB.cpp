#include<stdio.h>
int main()
{
	int a[170][170];
	int n,i,j;

	scanf("%d",&n);
	for(i=0;i<n;i++)
		for(j=0;j<n;j++)
			scanf("%d",&a[i][j]);
	for(i=0;i<n;i++)
	{
		for(j=0;j<n-1;j++)
			printf("%.8f ",(a[i][j]+a[j][i])/2.);
		printf("%.8f\n",(a[i][n-1]+a[n-1][i])/2.);
	}
	for(i=0;i<n;i++)
	{
		for(j=0;j<n-1;j++)
			printf("%.8f ",(a[i][j]-a[j][i])/2.);
		printf("%.8f\n",(a[i][n-1]-a[n-1][i])/2.);
	}
	return 0;
}