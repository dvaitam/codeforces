#include<cstdio>
#include<algorithm>
#include<cmath>

using namespace std;

#define FOR(i, a, b) for(int i = a; i <= b; i++)

const int MAX = 1010;

int main()
{
    int n, r;
    scanf("%d%d", &n, &r);
    int push[MAX] = {0};
    FOR(i, 0, n - 1)
        scanf("%d", &push[i]);
    int store[MAX] = {0};
    double ans[MAX] = {0};
    FOR(i, 0 , n - 1)
    {
        double mx = r;
        for(int j = 0; j < i; j++)
        {
            if(fabs(push[j] - push[i]) <= r + r)
            {
                double t = sqrt(4 * r * r - pow(fabs(push[j] - push[i]), 2)) + ans[j];
                mx = max(mx, t);
            }
        }
        ans[i] = mx;
    }
    FOR(i, 0, n - 1)
    {
        printf("%.10lf", ans[i]);
        if(i != n - 1)
            printf(" ");
    }
    return 0;
}