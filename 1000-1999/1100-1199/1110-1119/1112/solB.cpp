#include <bits/stdc++.h>
#define openfiles ifstream cin("input.txt");ofstream cout("output.txt");
#define faster ios_base::sync_with_stdio(false);cin.tie(NULL);cout.tie(NULL);
#define endl '\n'
#define ft first
#define sd second
#define ls (p<<1)
#define rs ((p<<1)|1)
using namespace std;
typedef long long ll;
typedef long double ld;
typedef pair <int, int> pii;
typedef pair <ll, ll> pll;
const int inf = 1'000'000'000;
const ll infll = 2'000'000'000;
const ld eps = 0.000'000'000'1;

struct kebab{
    int start{0};
    int finish{0};
    int i{-1};
};

int main(){
    //openfiles;
    faster;
    int n, k;
    cin >> n >> k;
    vector <int> a(n);
    for (int i = 0; i < n; i++){
        cin >> a[i];
    }
    vector <kebab> t(k);
    vector <int> ans(n, 0);
    int m = 0, time = 0, next_time = 0, mn, v;
    for (int i = 0; i < n; i++){
        mn = inf;
        v = -1;
        for (int j = 0; j < k; j++){
            if (t[j].finish < mn){
                mn = t[j].finish;
                v = j;
            }
        }
        ///
        if (t[v].finish) m++;
        time = t[v].finish;
        t[v].start = time;
        t[v].finish = time+a[i];
        t[v].i = i;
        ///
        mn = inf;
        v = -1;
        for (int j = 0; j < k; j++){
            if (t[j].finish < mn){
                mn = t[j].finish;
                v = j;
            }
        }
        next_time = t[v].finish;
        ///
        int u = int((100.*m/n)+0.5);
        for (int j = 0; j < k; j++){
            if (t[j].i == -1) continue;
            if (u+t[j].start-1 >= time && u+t[j].start <= next_time){
                ans[t[j].i] = 1;
            }
        }
    }
    while (true){
        mn = inf;
        v = -1;
        for (int j = 0; j < k; j++){
            if (t[j].i == -1) continue;
            if (t[j].finish < mn){
                mn = t[j].finish;
                v = j;
            }
        }
        //
        if (v == -1) break;
        //
        m++;
        time = t[v].finish;
        t[v].start = 0;
        t[v].finish = 0;
        t[v].i = -1;
        ///
        mn = inf;
        v = -1;
        for (int j = 0; j < k; j++){
            if (t[j].i == -1) continue;
            if (t[j].finish < mn){
                mn = t[j].finish;
                v = j;
            }
        }
        next_time = t[v].finish;
        ///
        int u = int((100.*m/n)+0.5);
        for (int j = 0; j < k; j++){
            if (t[j].i == -1) continue;
            if (u+t[j].start-1 >= time && u+t[j].start <= next_time){
                ans[t[j].i] = 1;
            }
        }
    }
    int sum = 0;
    for (int i = 0; i < n; i++) sum += ans[i];
    cout << sum;
    return 0;
}