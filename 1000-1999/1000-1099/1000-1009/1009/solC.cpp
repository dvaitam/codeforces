#include<bits/stdc++.h>
#define LL long long
using namespace std;

LL n,m,a,b,c,x,d;
long double sum,ans;

inline LL read()
{
    char ch=getchar();
    LL f=1,x=0;
    while (ch<'0' || ch>'9')
    {
        if (ch=='-') f=-1;
        ch=getchar();
    }
    while (ch>='0' && ch<='9')
    {
        x=x*10+ch-'0';
        ch=getchar();
    }
    return f*x;
}

int main()
{
    n=read(); m=read();
    sum=0;
    for (int i=1;i<=m;i++) 
    {
        x=read(); d=read();
        if (d>=0) sum+=x*n+((d*(n-1)*n)/2);
        else
        {
            sum+=x*n;
            LL mid;
            mid=n/2+1;
            //for (int j=1;j<=n;j++) sum+=d*abs(mid-j);
            LL front=mid-1,last=n-mid;
            sum+=d*((((1+front)*front)/2)+(((1+last)*last)/2));
        }
    }
    ans=sum/(long double) n;
    //printf("%.10llf",ans);
    cout<<fixed<<setprecision(15)<<ans;
    //system("pause");
    return 0;
}