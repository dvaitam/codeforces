#include <iostream>
#include <cmath>
#include <cstdio>
using namespace std;
int a[105],n;

int main()
{
    while(scanf("%d",&n) !=EOF)
    {
         int i,j;
         for(i=1;i<=n;i++)  scanf("%d",&a[i]);
         int ansi,ansj,temp;
         ansi = n;ansj = 1;
         temp = max(a[1] - a[n],a[n] -a[1]);
         for(i=2;i<=n;i++)
         {
               if(max(a[i-1] - a[i],a[i] - a[i-1]) < temp)
               {
                    ansi = i-1;ansj = i;
                    temp = max(a[i] - a[i-1],a[i-1] - a[i]);
               }
         }    
         printf("%d %d\n",ansi,ansj);            
    }
    return 0;
}