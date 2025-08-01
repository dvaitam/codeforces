#include <cstdio>
#include <cstring>
#include <algorithm>

using namespace std;

char s1[2500],s2[2500];
int ans[2500][2];

int main()
{
    #ifndef ONLINE_JUDGE
        freopen("input.txt","r",stdin);
        freopen("output.txt","w",stdout);
    #endif
    scanf("%s",s1+1);scanf("%s",s2+1);
    int l1=strlen(s1+1),l2=strlen(s2+1);
    int cnt=0,i=1;
    while (i<=l2)
    {
        int tl,tr,len=0;
        for (int j=1;j<=l1;j++)
            if (s1[j]==s2[i])
            {
                int k=j+1;
                while (k<=l1&&k-j+i<=l2&&s1[k]==s2[k-j+i]) k++;
                if (k-j>len) {tl=j;tr=k-1;len=k-j;}
            }
        for (int j=1;j<=l1;j++)
            if (s1[j]==s2[i])
            {
                int k=j-1;
                while (k&&j-k+i<=l2&&s1[k]==s2[j-k+i]) k--;
                if (j-k>len) {tl=j;tr=k+1;len=j-k;}
            }
        if (len==0) {puts("-1");return 0;}
        ans[++cnt][0]=tl;ans[cnt][1]=tr;
        i+=len;
    }
    printf("%d\n",cnt);
    for (int i=1;i<=cnt;i++) printf("%d %d\n",ans[i][0],ans[i][1]);
    return 0;
}