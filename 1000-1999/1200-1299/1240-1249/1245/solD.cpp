#include <bits/stdc++.h>
using namespace std;

typedef long long ll;
typedef pair<int, int> pii;

const int N = 2005;
const ll MOD = 1e9 + 7;
const ll INF = 2e18;

ll dis[N];
bool vis[N];
int pre[N];
pii node[N];
ll k[N], c[N];

int main() {
    int n;
    scanf("%d", &n);
    for(int i = 0; i < n; i++) scanf("%d%d", &node[i].first, &node[i].second);
    for(int i = 0; i < n; i++) { scanf("%d", &c[i]); dis[i] = c[i]; }
    for(int i = 0; i < n; i++) scanf("%d", &k[i]);
    ll ans = 0;
    vector<int> ans1;
    vector<pii> ans2;
    for(int i = 0; i < n; i++) {
        int v = -1;
        bool cho;
        for(int j = 0; j < n; j++) {
            if(!vis[j] && (v == -1 || dis[j] < dis[v])) {
                v = j;
                if(dis[v] == c[v]) cho = true;
                else cho = false;
            }
        }
        ans += dis[v];
        vis[v] = true;
        if(cho) ans1.push_back(v + 1);
        else ans2.push_back(pii(pre[v] + 1, v + 1));
        for(int j = 0; j < n; j++) {
            ll cost = (k[v] + k[j]) * (abs(node[v].first - node[j].first) + abs(node[v].second - node[j].second));
            if(cost < dis[j]) { dis[j] = cost; pre[j] = v; }
        }
    }
    printf("%lld\n", ans);
    printf("%d\n", ans1.size());
    for(int i = 0; i < ans1.size(); i++) {
        printf("%d%c", ans1[i], i == ans1.size() - 1 ? '\n' : ' ');
    }
    printf("%d\n", ans2.size());
    for(int i = 0; i < ans2.size(); i++) {
        printf("%d %d\n", ans2[i].first, ans2[i].second);
    }
    return 0;
}