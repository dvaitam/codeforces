#include<cstdio>

#define MAXX 100111

int n,i,j,k;
int a[MAXX];
bool as[MAXX];
long long sum;

int main()
{
    scanf("%d",&n);
    for(i=0;i<n;++i)
        scanf("%d",a+i);
    for(i=n-1;i>=0;--i)
        if(sum>0)
        {
            sum-=a[i];
            as[i]=true;
        }
        else
            sum+=a[i];
    if(sum>0)
        for(i=0;i<n;++i)
            as[i]^=true;
    for(i=0;i<n;++i)
        putchar(as[i]?'+':'-');
    puts("");
}