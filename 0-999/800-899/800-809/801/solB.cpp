#include <stdio.h>
#include <string.h>
char x[102];
char y[102];
char z[102];
main()
{
    scanf("%s",x);
    scanf("%s",y);
    int bomb=0,i=0;
    for (;i<strlen(x);i++)
    {
        if (x[i]>=y[i])
        {
            z[i]=y[i];
        }
        else if (x[i]<y[i])
        {
            printf("-1");
            return 0;
        }
    }
    z[i]='\0';
    printf("%s",z);
}