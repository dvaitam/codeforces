#include<iostream>
#include<cstdio>
#include<cstring>
using namespace std;
int vis[105],a[105],b[105],c[105],d[105];
int main()
{
    int n;
    scanf("%d",&n);
    for(int i=1;i<=n;++i)
        scanf("%d%d%d%d",a+i,b+i,c+i,d+i);
    for(int i=1;i<=n;++i)
        for(int j=1;j<=n;++j)
        if(a[i]<a[j]&&b[i]<b[j]&&c[i]<c[j])vis[i]=1;
    int ans=1005,num;
    for(int i=1;i<=n;++i)
        if(!vis[i]&&ans>d[i])ans=d[i],num=i;;
    printf("%d\n",num);
}