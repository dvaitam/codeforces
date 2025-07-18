#include <cstdio>
#include <algorithm>
#include <cstring>
using namespace std;
char a[300005];
int main()
{
    int t;
    scanf("%d", &t);
    while(t--)
    {
        scanf("%s", a);
        int len = strlen(a);
        bool judge = 1;
        for (int i = 0; i < len;i++)
            if(a[i]!='a')
            {
                judge = 0;
                break;
            }
        if(judge)
        {
            printf("NO\n");
            continue;
        }
        bool ans = 0;
        int left = 0, right = len - 1;
        while(left<right)
        {
            if(a[left]==a[right] && a[left]=='a')
            {
                left++;
                right--;
            }
            else
            {
                if(a[left]=='a')
                    ans = 0;
                else if(a[right]=='a')
                    ans = 1;
                else
                    ans = 0;
                break;
            }
        }
        //printf("%d %d\n", left, right);
        if(left>right)
        {
            printf("NO\n");
            continue;
        }
        printf("YES\n");
        if(ans==0)
            printf("a%s\n", a);
        else
            printf("%sa\n", a);
    }
    return 0;
}