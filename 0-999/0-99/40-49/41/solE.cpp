#include<cstdio>
int main()
{
    int n,i,j;
    scanf("%d",&n);
    printf("%d\n",(n>>1)*(n-(n>>1)));
    for(i=0;i<(n>>1);i++)
        for(j=(n>>1);j<n;j++)
            printf("%d %d\n",i+1,j+1);
}