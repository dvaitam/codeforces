#include<bits/stdc++.h>
#define ll long long
using namespace std;
int main()
{
	int n;
	scanf("%d",&n);
	if(n%2==0)
	{
		if(n==2)
		{
			printf("1\n1 1");
		}
		else
		{
			if(n%4==0)
			printf("0\n");
			else
			printf("1\n");
			printf("%d ",n/2);
			int flag=1;
			for(int i=1;i<=n;i+=2)
			{
				if(flag)
					printf("%d ",i);
				else
					printf("%d ",i+1);
				flag^=1;
			}
		}
	}
	else
	{
		if(n==3)
		{
			printf("0\n");
			printf("1 3");
		}
		else
		{
			if(n%4==1)
			printf("1\n%d ",n/2);
			else
			printf("0\n%d 1 ",n/2+1);
			int flag=1;
			for(int i=2;i<=n;i+=2)
			{
				if(flag)
					printf("%d ",i);
				else
					printf("%d ",i+1);
				flag^=1;
			}
		}
	}
}