#include <stdio.h>
#include <stdlib.h>

int main()
{
    int a;
    scanf("%d",&a);
    if(a>36)
    printf("-1\n");
    else
    {
       if(a%2==0)
       {a=a/2;
       while(a--)
       printf("8");
       printf("\n");
       }
    else
    {
       a=a/2;
       while(a--)
       printf("8");
       printf("4\n");
    }
    }
    return 0;
}