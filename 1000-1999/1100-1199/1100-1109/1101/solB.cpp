#include<cstdio>
#include<cstring>
using namespace std;
const int maxn=500010;
char str[maxn];
int main()
{
    scanf("%s",str);
    int len=strlen(str);
    int flag=0,flag1=0;
    int st1,st2,en1,en2;
    for(int i=0;i<len;i++)
    {
        if(str[i]=='['&&flag==0)
        {
            st1=i;
            flag=1;
        }
        if(str[i]==':'&&flag1==0&&flag)
        {
            st2=i;
            flag1=1;
        }
                if(flag&&flag1)
            break;
    }
    if(flag==0||flag1==0)
    {
        printf("-1\n");
    }
    else{
    flag=0,flag1=0;
    for(int i=len-1;i>st2;i--)
    {
        if(str[i]==']'&&flag==0)
        {
            en1=i;
            flag=1;
        }
        if(str[i]==':'&&flag1==0&&flag)
        {
            en2=i;
            flag1=1;
        }
        if(flag&&flag1)
            break;
    }
    if(flag==0||flag1==0)
        {
            printf("-1\n");
        }
        else
        {
            int cnt=4;
            for(int i=st2+1;i<en2;i++)
            {
                if(str[i]=='|')
                    cnt++;
            }
            printf("%d\n",cnt);
        }
    }
    return 0;
}