#include <stdio.h>
#include <vector>
#include <algorithm>
using namespace std;
int read()
{
    char ch = getchar();
    int f = 1;
    int x = 0;
    while(ch < '0' || ch > '9'){if(ch == '-')f = 0;ch = getchar();}
    while(ch >= '0' && ch <= '9'){x = x * 10 + ch - '0';ch = getchar();}
    return f?x:x*-1;
}
int n;
int a[100001];
int b[100001];
bool cmp(int a,int b)
{
    return a > b;
}
int main()
{
    n = read();
    long long sum1 = 0,sum2 = 0;
    
    for(int i = 0;i < n;i ++)
    {
        a[i] = read();
    }
    for(int i = 0;i < n;i ++)
    {
        b[i] = read();
    }
    sort(a,a + n,cmp);
    sort(b,b + n,cmp);

    int i = 0,j = 0,flag = 0;
    
    while(i != n || j != n)
    {
        if(flag == 0)
        {
            if(j == n || a[i] > b[j])
            {
                sum1 += a[i];
                i ++;
            }
            else
            {
                j ++;
            }
            flag = 1;
        }
        else
        {
            if(i == n || a[i] < b[j])
            {
                sum2 += b[j];
                j ++;
            }
            else
            {
                i ++;
            }
            flag = 0;
        }
    }
    printf("%I64d ",sum1 - sum2);

    return 0;
}