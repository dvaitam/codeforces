#include <stdio.h>
#include <time.h>
#include <algorithm>
using namespace std;


double dp[2][310]={{0}};


int main()
{
	int i1,i2;
	int x,k,p;
	int r;
	int yn;
	double ans=0,ans2;
	double pd,pp;
	int next=1;
	srand(time(NULL));
	scanf("%d %d %d",&x,&k,&p);
	
	
	pd=(double)p/100;
	pp=(double)(100-p)/100;
	dp[0][0]=1;
	for(i1=0; i1<k; i1++)
	{
		dp[next][0]=0;
		for(i2=0; i2<300; i2++)
		{
			dp[next][i2+1]=dp[1-next][i2]*pp;
		}
		for(i2=0; i2<300; i2+=2)
		{
			ans+=dp[1-next][i2]*pd;
			dp[next][i2/2]+=dp[1-next][i2]*pd;
		}
		next=1-next;
	}
	for(i1=0; i1<300; i1++)
	{
		r=i1+x;
		while(r%2==0 )
		{
			r/=2;
			ans+=dp[1-next][i1];
		}
	}
	printf("%lf",ans);
	 
	return 0;
}