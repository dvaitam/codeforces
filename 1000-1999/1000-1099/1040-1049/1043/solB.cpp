#include<bits/stdc++.h>
using namespace std;
const int N = 2e3+5;
int T, n, m;
int arr[N];
int dif[N];
bitset<N> per;
int main()
{
scanf("%d", &n);
for(int i = 0; i < n; ++i)
{
    scanf("%d", arr + i);
}
dif[0] = arr[0];
for(int i = 1; i < n; ++i)
{
    dif[i] = arr[i] - arr[i - 1];
}
int ans = n;
for(int i = 1; i <= n; ++i)
{
    per[i] = 1;
    for(int x = 0; x < n; ++x)
    {
        if(dif[x] != dif[x % i])
        {
            per[i] = 0;
            --ans;
            break;
        }
    }
}
printf("%d\n", ans);
for(int i = 1; i <= n; ++i)
{
    if(per[i])
        printf("%d ", i);
}
puts("");
}