//AshrafSustS19


#include<bits/stdc++.h>

using namespace std;
using ll = long long int;

// #define LOCAL 1
#ifdef LOCAL
#include "C:\Users\Zadeed\Documents\AshrafSustS19\app_debug.cpp"
#else
#define debug(...)
#define debugArr(arr, n)
#endif


int n, m, np;
vector<int> adj;
vector<int> MEMadj;

vector<vector<int>> MEM;
vector<int> lnxt;


int dp(int u, int bm){
    if (MEM[u][bm] != -1){
        return MEM[u][bm];
    }
    int res = 0;
    for (int i = 0; i < n; i++){
        int v = (1 << i);
        if ((v & bm) && (v & adj[u])){
            if ((v ^ bm) == 0){
                res = 1;
                lnxt[u] = i;
                return MEM[u][bm] = res;
            }
            if (bm % v != 0){
                
                int p = dp(i, bm ^ v);;
                if (p == 1){
                    lnxt[u] = i;
                    return MEM[u][bm] = p;

                }
            }
        }
    }
    return MEM[u][bm] = res;
    
}



int main(){
    ios_base::sync_with_stdio(false);
    cin.tie(0);
    cin >> n >> m;
    np = 1 << n;
    adj.resize(n, 0);
    MEMadj.resize(np, 0);
    MEM.resize(n, vector<int> (np, -1));

    for (int i = 0; i < m; i++){
        int a, b;
        cin >> a >> b;
        adj[--a] |= (1 << (--b));
        adj[b] |= 1 << a;
    }
    
    for (int i = 0; i < np; i++){
        for (int j = 0; j < n; j++){
            if ((1 << j) & i){
                MEMadj[i] |= adj[j];
            }
        }
    }
    bool ispos = false;
    lnxt.resize(n, -1);
    for (int i = 0; i < np && !ispos; i++){
        if (MEMadj[i] == np - 1){
            for (int j = 0; j < n; j++){
                if ((1 << j) & i){
                    
                    int res = dp(j, i);
                    if (res == 1){
                        if (lnxt[j] == j){
                            for (int k = 0; k < n; k++){
                                if (k == j) continue;
                                if ((1 << k) & adj[j]){
                                    lnxt[j] = k;
                                }
                                break;
                            }
                        }
                        ispos = true;
                        break;
                    }
                    break;
                }
            }
        }
    }
    if (!ispos){
        cout << "NO\n";
        return 0;
    }
    for (int i = 0; i < n; i++){
        if (lnxt[i] == -1){
            for (int j = 0; j < n; j++){
                if (((1 << j) & adj[i]) && lnxt[j] != -1){
                    lnxt[i] = j;
                    break;
                }
            }
        }
    }
    cout << "YES\n";
    for (int i = 0; i < n; i++){
        cout << lnxt[i] + 1 << " ";
    }
    cout << "\n";
}