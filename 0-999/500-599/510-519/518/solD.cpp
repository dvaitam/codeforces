#include<bits/stdc++.h>

using namespace std;



typedef  long long ll;

const double eps=1e-7;



double dp[2009];

int main()

{

    int n,t;

    double p;

    scanf("%d%lf%d",&n,&p,&t);

    dp[0]=1;

    for(int time=1;time<=t;time++)

    {

        dp[n]+=p*dp[n-1];

        for(int num=n-1;num>0;num--)

            dp[num]=p*dp[num-1]+(1-p)*dp[num];

        dp[0]*=(1-p);



    }





    double ans=0;

    for(int num=0;num<=n;num++)ans+=num*dp[num];

    printf("%.10f\n",ans);



}