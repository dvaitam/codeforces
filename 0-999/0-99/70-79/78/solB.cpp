#include <iostream>
#include <stdio.h>
#include <string.h>
using namespace std;

int c[105];
int visit[10];
int n;
bool ans;

void dfs(int k)
{
    if (k>n)
    {
        memset(visit,0,sizeof(visit));
        for (int p=1;p<=n;p++)
            visit[c[p]]=1;
        for (int q=1;q<=7;q++)
            if (visit[q]==0)
                return;
        ans=true;
        return;
    }


    for (int i=1;i<=7;i++)
    {
        bool flag=true;
        if (k>=4)
        {
            c[k]=i;
            for (int j=k-4+1;j<=k-1;j++)
            {
                for (int l=j+1;l<=k;l++)
                {
                    if (c[j]==c[l])
                    {                                       
                        flag=false;
                        break;
                    }
                }
            }
        }
        if (flag==true)
        {
            c[k]=i;
            dfs(k+1);
        }
        else 
        {
            c[k]=0;
        }
        if (ans==true) return;
    }
}






int main()
{
    //red, orange, yellow, green, blue, indigo or violet
    while (scanf("%d",&n)!=EOF)
    {
        ans=false;
        dfs(1);
        for (int i=1;i<=n;i++)
        {
            if (c[i]==1) printf("R");
            if (c[i]==2) printf("O");
            if (c[i]==3) printf("Y");
            if (c[i]==4) printf("G");
            if (c[i]==5) printf("B");
            if (c[i]==6) printf("I");
            if (c[i]==7) printf("V");
        }
        printf("\n");
    }
    return 0;
}