#include<cstdio>
#include<cstring>
#include<algorithm>
#include<cmath>
#include<iostream>
using namespace std;
int n,m,p,t;
int gcd(int a,int b){return b?gcd(b,a%b):a;}
void imp(){
	printf("Impossible");
	
}
int main()
{
	t=0;
	scanf("%d%d",&n,&m);
	
	if(m<(n-1)) 
	{
		imp();
		return 0;
	}
	
	if(n<1000)
	{
		for(int i=1;i<=n;i++)
			for(int j=i+1;j<=n;j++)
			{
				p=gcd(i,j);
				if (p==1) t++;
			}
		if(t<m) 
		{
			imp();
			return 0;
		}
	}
	printf("Possible\n");
	for(int i=1;i<=n;i++)
		for(int j=i+1;j<=n;j++)
		{
			p=gcd(i,j);
			if(p==1)
			{
				printf("%d %d\n",i,j);
				m--;
				if(!m) return 0;
			}
		}
		
	return 0;
}