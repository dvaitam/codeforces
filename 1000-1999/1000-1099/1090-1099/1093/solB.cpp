#include<cstdio>
#include<cstring>
#include<algorithm>
using namespace std;
const int N=1005;
char str[N];
inline int ispald(char *str,int len)
{
    int res=-2;
    for(int i=0;i<=len/2;i++)
    {
        if(str[i]!=str[len-i-1])
            return -1;
        if(i&&str[i]!=str[i-1])
            res=i;
    }
    return res;
}
int main()
{
    int t;
    scanf("%d",&t);
    while(t--)
    {
        scanf("%s",str);
        int len=strlen(str);
        int x=ispald(str,len);
        if(x==-1)
            puts(str);
        else if(x!=-2)
            swap(str[x],str[x-1]),puts(str);
        else
            puts("-1");
    }
    return 0;
}