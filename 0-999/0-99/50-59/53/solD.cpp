#include<cstdio>
#include<algorithm>
using namespace std;
int main()
{
	int a[310],b[310],swp[100000],cnt=0,n;
	scanf("%d",&n);
	for(int i=0;i<n;i++)
		scanf("%d",a+i);
	for(int i=0;i<n;i++)
		scanf("%d",b+i);
	for(int i=0;i<n;i++)
	{
		if(a[i]!=b[i])
		{
			int j=i;
			for(;j<n;j++)
				if(a[i]==b[j])
					break;
			for(int k=j;k>i;k--)
			{
				swp[cnt++]=k;
				swap(b[k],b[k-1]);
			}
		}
	}
	printf("%d\n",cnt);
	for(int i=0;i<cnt;i++)
		printf("%d %d\n",swp[i],swp[i]+1);
	return 0;
}