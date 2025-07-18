#include <cstdio>
#include <algorithm>
#include <cstring>
#include <stack>
#include <queue>
#include <vector>
#include <cmath>
#include <set>
using namespace std;
#define N 200001
#define INF 0x3f3f3f3f
struct node {
    int a;
    int b;
}c[N];
bool cmp(node a, node b) {
    if(a.b == b.b) {
        return a.a < b.a;
    }
    return a.b < b.b;
}
int main()
{
    int n , k, t;
    int ans=0;
    int maxi = 0;
    int flag = 0;
    scanf("%d", &n);
    for(int i = 1; i <= n; ++i) {
        scanf("%d", &c[i].a);
    }
    for(int i = 1; i <= n; ++i) {
        scanf("%d", &c[i].b);
        if(c[i].a != 0 && c[i].b != 0) {
            if(c[i].a < 0 && c[i].b < 0) {
                int gcd = __gcd(-c[i].a, -c[i].b);
                c[i].a = (-c[i].a)/gcd;
                c[i].b = (-c[i].b)/gcd;
            }else if(c[i].a < 0) {
                int gcd = __gcd(-c[i].a, c[i].b);
                c[i].a = (c[i].a/gcd);
                c[i].b = (c[i].b/gcd);
            }else if(c[i].b < 0) {
                int gcd = __gcd(c[i].a, -c[i].b);
                c[i].a = (-c[i].a/gcd);
                c[i].b = (-c[i].b/gcd);
            }else {
                int gcd = __gcd(c[i].a, c[i].b);
                c[i].a = (c[i].a/gcd);
                c[i].b = (c[i].b/gcd);
            }
        }else {
            if(c[i].a == 0 && c[i].b == 0) {
                ans ++;
                continue;
            }
            if(c[i].a == 0 && c[i].b != 0) {
                flag ++;
            }
            if(c[i].b == 0){
                maxi ++;
            }
        }
    }
    sort(c + 1, c + n + 1, cmp);
    if(flag != (n-ans)) {
        maxi = max(maxi, 1);
    }
    int sum = 0;
    for(int i = 2; i <= n; ++i) {
        if(c[i].a == 0 && c[i].b == 0) {
            continue;
        }
        if(c[i].b == 0)
            continue;
        if(c[i].a == c[i-1].a && c[i].b == c[i-1].b && c[i].a != 0) {
            sum ++;
        }else {
            if(sum)
                maxi = max(maxi, sum+1);
            sum = 0;
        }
    }
    if(sum)
        maxi = max(maxi, sum+1);
    printf("%d\n", min(maxi+ans,n));
    return 0;
}