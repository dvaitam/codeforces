#include <bits/stdc++.h>
using namespace std;
 
using ll = long long;
 
const int MAXN = 200005;
const ll MOD = 998244353;
 
struct Query {
    int l, r, id;
    Query(int _l = 0, int _r = 0, int _id = 0) : l(_l), r(_r), id(_id) {}
};
 
int n, m, qCount;
int shelf[MAXN], blockIdx[MAXN], cnt[MAXN], totalCount[MAXN];
ll invArr[MAXN], pArr[MAXN], currentProd, answer[MAXN];
vector<Query> queriesByK[100005];
Query tempQueries[MAXN];
 
ll modExp(ll base, ll exp) {
    ll res = 1;
    while(exp) {
        if(exp & 1)
            res = (res * base) % MOD;
        base = (base * base) % MOD;
        exp >>= 1;
    }
    return res;
}
 
int BLOCK_SIZE;
 
bool moComparator(const Query &a, const Query &b) {
    if(blockIdx[a.l] != blockIdx[b.l])
        return blockIdx[a.l] < blockIdx[b.l];
    return a.r < b.r;
}
 
void addFilm(int x, int currentK) {
    currentProd = (currentProd * (totalCount[shelf[x]] + currentK - cnt[shelf[x]])) % MOD;
    cnt[shelf[x]]++;
}
 
void removeFilm(int x, int currentK) {
    cnt[shelf[x]]--;
    currentProd = (currentProd * invArr[ totalCount[shelf[x]] + currentK - cnt[shelf[x]] ]) % MOD;
}
 
void processQueriesForK(int kVal) {
    for (int i = 1; i <= n; i++)
        cnt[shelf[i]] = 0;
    int currentK = kVal;
    currentProd = 1;
    ll zz = (1LL * m * currentK) % MOD;
    pArr[n] = 1;
    for (int i = n - 1; i >= 1; i--)
        pArr[i] = (pArr[i+1] * (zz + n - i)) % MOD;
    int numQueries = 0;
    for (auto &qItem : queriesByK[kVal])
        tempQueries[++numQueries] = qItem;
    sort(tempQueries + 1, tempQueries + numQueries + 1, moComparator);
    int currentL = 1, currentR = 0;
    for (int i = 1; i <= numQueries; i++){
        int L = tempQueries[i].l, R = tempQueries[i].r;
        while (currentL > L)
            addFilm(--currentL, currentK);
        while (currentR < R)
            addFilm(++currentR, currentK);
        while (currentL < L)
            removeFilm(currentL++, currentK);
        while (currentR > R)
            removeFilm(currentR--, currentK);
        answer[tempQueries[i].id] = (currentProd * pArr[currentR - currentL + 1]) % MOD;
    }
}
 
int main(){
    ios::sync_with_stdio(false);
    cin.tie(nullptr);
    
    cin >> n >> m >> qCount;
    for (int i = 1; i <= n; i++){
        cin >> shelf[i];
        totalCount[shelf[i]]++;
    }
    
    invArr[1] = 1;
    for (int i = 2; i < MAXN; i++)
        invArr[i] = (invArr[MOD % i] * (MOD - MOD / i)) % MOD;
    
    BLOCK_SIZE = 1000;
    for (int i = 1; i <= n; i++)
        blockIdx[i] = (i - 1) / BLOCK_SIZE + 1;
    
    for (int i = 1; i <= qCount; i++){
        int L, R, kVal;
        cin >> L >> R >> kVal;
        queriesByK[kVal].push_back(Query(L, R, i));
    }
    
    for (int kVal = 0; kVal <= 100000; kVal++){
        if (!queriesByK[kVal].empty())
            processQueriesForK(kVal);
    }
    
    for (int i = 1; i <= qCount; i++)
        cout << answer[i] << "\n";
    
    return 0;
}