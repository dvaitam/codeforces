/*jai mata di
 let's rock*/
#include<bits/stdc++.h>
using namespace std;
const int N=1004;
char s[N][N];
int main()
{
    int n,a,b;
    scanf("%d %d %d",&n,&a,&b);
    if(n==1)
    {
        printf("YES\n");
        printf("0");
        return 0;
    }
    if(n==2)
    {
        if(a==1 && b==1)
        {
            printf("NO\n");
            return 0;
        }
    }
    if(n==3)
    {
        if(a==1 && b==1)
        {
            printf("NO\n");
            return 0;
        }
    }
    int i,j;
    for(i=1;i<=n;i++)
        for(j=1;j<=n;j++)
           s[i][j]='0';
    if(a>1)
    {
        if(b!=1)
        {
            printf("NO\n");
            return 0;
        }
        else
        {
            printf("YES\n");
            int val=n-a+1;
            for(i=1;i<val;i++)
                {
                    s[i][i+1]='1';
                    s[i+1][i]='1';
                }
            for(i=1;i<=n;i++)
               printf("%s\n",s[i]+1);
        }
    }
    if(a==1)
    {
        if(b==1)
        {
            printf("YES\n");
            for(i=1;i<n;i++)
            {
                s[i][i+1]='1';
                s[i+1][i]='1';
            }
            for(i=1;i<=n;i++)
               printf("%s\n",s[i]+1);
            return 0;
        }
        printf("YES\n");
        for(i=1;i<=n;i++)
            for(j=1;j<=n;j++)
                if(i!=j) s[i][j]='1';
        int val=n-b+1;
        for(i=1;i<=val;i++)
        {
            s[1][i]='0';
            s[i][1]='0';
        }
        for(i=1;i<=n;i++)
            printf("%s\n",s[i]+1);
    }
    return 0;
}