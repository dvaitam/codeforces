#include<cstdio>
using namespace std;
int n,a[100020];
int main()
{
    int flag=-1;
    scanf("%d",&n);
    for(int i=1; i<=n; i++)
        scanf("%d",&a[i]);
    if(a[n]==0)
    {
        for(int i=1; i<n; i++)
            if(a[i]==0)
            {
                flag=i;
                break;
            }

        if(flag==n-1)
        {
            printf("NO");
            return 0;
        }

        printf("YES\n");
        if(flag>0)
        {
            if(flag>1)
            {
                printf("(");
                for(int i=1; i<flag-1; i++)
                    printf("%d->",a[i]);
                printf("%d)->",a[flag-1]);
            }
            printf("(0)->");
            printf("(");
            for(int i=flag+1; i<n-1; i++)
                printf("%d->",a[i]);


            printf("%d)",a[n-1]);

            printf("->0");
        }
        else
        {
            for(int i=1; i<n; i++)
                printf("%d->",a[i]);
            printf("0");
        }
    }
    else
        printf("NO\n");
    return 0;
}