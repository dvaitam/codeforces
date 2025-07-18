//BISMILLAHIR RAHMANIR RAHIM
#include<stdio.h>
int main()
{
    int n,i,A[305];
    scanf("%d",&n);
    for(i=0;i<n;i++)
        scanf("%d",&A[i]);
    for(i=0;i<n;i++)
    {
        while(A[i])
        {
            putchar('P');
            A[i]--;
            if(A[i])
            {
                if(i==n-1)
                    putchar('L'),putchar('R');
                else
                    putchar('R'),putchar('L');
            }
        }
        if(i!=n-1)
            putchar('R');
    }
    printf("\n");
    return 0;
}