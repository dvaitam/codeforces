#include <algorithm>
#include <iostream>
#include <string.h>
using namespace std;
int main()
{
	int a[110],b[110];
	int n,k;
	scanf("%d%d",&n,&k);
	int i;
	for(i=0;i<n;i++)
	{
		scanf("%d",&a[i]);
	}
	int ans=0,sum=0,j;
	for(i=0;i<k;i++)
	{
		ans=0;
		for(j=0;j<n;j++)
		{
			b[j]=a[j];
		}
		for(j=i;j<n;j=j+k)
		{
			b[j]=0;
		}
		for(j=0;j<n;j++)
		{
		ans=ans+b[j];
		}
		if(ans<0)ans=ans*-1;
		if(ans>sum)
		{
			sum=ans;
		}
	}
	printf("%d\n",sum);
	
}