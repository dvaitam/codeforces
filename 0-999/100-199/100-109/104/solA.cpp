#include<cstdio>
int main()
{
    int m,n;

    scanf("%d",&n);
    m=n-10;
    if(m<1 || m>11)
        printf("0\n");
    else if(m==10)
        printf("15\n");
    else printf("4\n");
return 0;
}