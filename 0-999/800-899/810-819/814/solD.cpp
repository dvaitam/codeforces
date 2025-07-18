#include<iostream>
#include<cstdio>
#include<cstring>
#include<vector>
#include<map>
#include<queue>
#include<cmath>
#include<algorithm>
#define lson l, mid, i<<1
#define rson mid+1, r, i<<1|1
#define PII pair<int, int>
using namespace std;
typedef long long LL;
//const int MOD = 1e9+7;
const auto INF  = 0x3f3f3f3f;
const int N = 1e3+5;
const double PI = acos(-1.0);
struct node {
    double x, y, r;
}maps[N];

bool cmp(const struct node &u, const struct node &v) {
    return u.r > v.r;
}

int ans[N];
bool visit[N];

bool ok(int i, int j) {
    double temp = (maps[i].x-maps[j].x)*(maps[i].x-maps[j].x) +
        (maps[i].y-maps[j].y)*(maps[i].y-maps[j].y);
    temp = sqrt(temp);
    if(temp < maps[i].r+maps[j].r)
        return true;
    return false;
}

int main() {
    int n;
    scanf("%d", &n);
    for(int i=1; i<=n; i++) {
        scanf("%lf%lf%lf", &maps[i].x, &maps[i].y, &maps[i].r);
    }
    sort(maps+1, maps+1+n, cmp);
    
//    double answer = 0;
//    for(int i=1; i<=n; i++) {
//        if(i&1)
//            answer += maps[i].r*maps[i].r;
//        else
//            answer -= maps[i].r*maps[i].r;
//    }
//    printf("answer=%.9lf\n", answer * PI);
//    
//    return 0;
//    
    
    ans[1] = 0;
    visit[1] = true;
    double sum = 0;
    sum += maps[1].r*maps[1].r;
    for(int i=2; i<=n; i++) {
        bool flag = false;
        for(int j=i-1; j>=1; j--) {
            if(ok(i, j)) {
                flag = true;
                if(!visit[j]) {
                    sum += maps[i].r*maps[i].r;
                    visit[i] = true;
                    ans[i] = ans[j];
                    break;
                }
                ans[i] = !ans[j];
                flag = false;
                for(int k=i-1; k>=1; k--) {
                    if(ans[k] == ans[i]) {
                        if(ok(i, k)) {
                            flag = true;
                            if(visit[k])
                                sum -= maps[i].r*maps[i].r;
                            else {
                                sum += maps[i].r*maps[i].r;
                                visit[i] = true;
                            }
                            break;
                        }
                    }
                }
                if(!flag) {
                    visit[i] = true;
                    sum += maps[i].r*maps[i].r;
                }
                flag = true;
                break;
            }
        }
        if(!flag) {
            ans[i] = 0;
            sum += maps[i].r*maps[i].r;
            visit[i] = true;
        }
    }
    printf("%.9lf\n", sum*PI);
}