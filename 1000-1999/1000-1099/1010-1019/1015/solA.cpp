#include<iostream>
#include<cstdio>
using namespace std;

int points[105];

int main()
{
    int n,m,l,r,o;
    scanf("%d%d",&n,&m);
    for (int i = 0; i < n; i++)
    {
        scanf("%d%d",&l,&r);
        for (int j = l; j <= r; j++)
        {
            points[j] = 1;
        }
    }
    o = 0;
    for (int i = 1; i <= m; i++)
    {
        o += 1 - points[i];
    }
    printf("%d\n",o);
    o = 0;
    for (int i = 1; i <= m; i++)
    {
        if (!points[i])
        {
            printf("%s%d",o ? " " : "",i);
            o++;
        }
    }
    printf("\n");
    return 0;
}