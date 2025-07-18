#include<cstdio>
#include<cstdlib>

int T;
int b[1006],a[1005];
int ans[3],t,n;
bool bo;
int kano()
{
    char ch=getchar();
    int w=1,ans=0;
    for(;ch<'0'||ch>'9';ch=getchar())if(ch=='-')w=-1;
    for(;ch>='0'&&ch<='9';ch=getchar())ans=ans*10+ch-'0';
    return w*ans;
}

int main()
{
    T=kano();
    while(T--)
    {

        t=0;
        bo=0;

        n=kano();
        for(int i=1;i<=n;i++)a[i]=kano();
        for(int i=1;i<=n;i++)
        {
            while(t>0&&a[b[t]]>a[i])
            {
                if(t>1)
                {
                    bo=1;
                    ans[0]=b[t-1];
                    ans[1]=b[t];
                    ans[2]=i;
                }
                t--;
            }
            b[++t]=i;
        }
        if(bo)
        {
            printf("YES\n%d %d %d",ans[0],ans[1],ans[2]);
        }
        else
        {
            printf("NO");
        }
        if(T!=0){
            putchar('\n');
        }
    }
}