#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#define maxn 100010
int ans[maxn][2],p,q;
int main()
{
    int n,x,p=0,q=0,s=0,i;
    scanf("%d",&n);
    for(i=1;i<=n+1;i++){
        if(i<=n) scanf("%d",&x);
        else x=0;//直到所有都出光
        while(s<x)ans[p++][0]=i,s++;//要繼續增加`,,qu
        //如果s==x,可以用同樣的edit做完他
        while(s>x)ans[q++][1]=i-1,s--;//要在前一個完掉
    }
    printf("%d\n",p);//個數是入的序號
    for(i=0;i<p;i++)printf("%d %d\n",ans[i][0],ans[i][1]);
    return 0;
}