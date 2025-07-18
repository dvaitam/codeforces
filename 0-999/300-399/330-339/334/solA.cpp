#include <cstdio>
using namespace std;

int main()
{
    int b,candy;
    scanf("%d",&b);
    for(int i = 1; i <=b; i++)
    {
        for(int k = 0; k < b; k++)
        {
            if(k < b/2) printf("%d ",k*b+i);
            else printf("%d ",k*b+b-(i-1));
        }
        printf("\n");
    }
    return 0;
}