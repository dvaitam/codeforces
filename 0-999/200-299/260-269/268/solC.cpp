#include<stdio.h>
#include<algorithm>
using namespace std;
int main()
{
    int n,m;
    scanf("%d %d",&n,&m);
    int min;

    if(n<m) min=n;
    else min = m;

    int l=min;
    printf("%d\n",min+1);
    int i;
    for(i=0;i<=min;i++,l--)
    {
        if(i==min) printf("%d %d",i,l);
        else printf("%d %d\n",i,l);
    }
    return 0;
}