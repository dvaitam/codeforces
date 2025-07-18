#include<iostream>
#include<cstdio>
#include<cstdlib>
#include<cstring>
#include<cmath>
#include<algorithm>
#define maxn 5010
#define each(e,x) for(__typeof(x.begin()) e=x.begin();e!=x.end();++e)
using namespace std;
typedef long long LL;
double dp[2][maxn],pw[maxn];
void read()
{
    int n,T;
    cin>>n>>T;
    int cur=0;
    static int x[maxn],y[maxn];
    for(int i=1;i<=n;++i)
        cin>>x[i]>>y[i];
    for(int i=n;i>=1;--i)
    {
        double p=x[i]/100.0,sum=0,P=1;
        cur^=1;
        for(int j=1;j<y[i];++j)
            P*=1-p;
        for(int j=1;j<=T;++j)
        {
            sum*=1-p;
            sum+=(dp[cur^1][j-1]+1)*p;
            if(j>=y[i])
                sum+=(dp[cur^1][j-y[i]]+1)*P*(1-p);
            if(j>y[i])
                sum-=(dp[cur^1][j-y[i]-1]+1)*P*(1-p);
            dp[cur][j]=sum;
        }
    }
    printf("%.10f\n",dp[cur][T]);
}
int main()
{
    read();
    return 0;
}