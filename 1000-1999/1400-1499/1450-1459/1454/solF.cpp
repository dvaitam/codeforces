#include <bits/stdc++.h>
#define fo(i,m,n) for(int i = m; i <= n; ++i)
#define fod(i,m,n) for(int i = m; i >= n; --i)
#define fov(i,x) for(auto &i : x)
using namespace std;
typedef int64_t i64;
typedef pair<int,int> _ii;
const int N = 2e5+7;
int a[N], n, cnt, L[N], R[N], x,y,z;
int max_L[N], max_R[N];
stack <int> A;
void Pre()
{
    while(A.size()) A.pop();
    a[0] = 0;
    A.push(0);
    fo(i,1,n)
    {
        while(a[A.top()] >= a[i]) A.pop();
        L[i] = A.top()+1;
        A.push(i);
    }
    while(A.size()) A.pop();
    a[n+1] = 0;
    A.push(n+1);
    fod(i,n,1)
    {
        while(a[A.top()] >= a[i]) A.pop();
        R[i] = A.top()-1;
        A.push(i);
    }
    max_L[0] = 0;
    fo(i,1,n) max_L[i] = max(max_L[i-1],a[i]);
    max_R[n+1] = 0;
    fod(i,n,1) max_R[i] = max(max_R[i+1],a[i]);
}
bool Run()
{
    fo(i,1,n)
    {
        int u = max_L[L[i]-1];
        int v = max_R[R[i]+1];
        if(u < a[i] && L[i] != i && a[L[i]] == a[i]) {u = a[i]; L[i]++;}
        if(v < a[i] && R[i] != i && a[R[i]] == a[i]) {v = a[i]; R[i]--;}
        if(u == v && u == a[i])
        {
            x = L[i] - 1;
            y = R[i] - L[i] + 1;
            z = n-x-y;
            return true;
        }
    }
    return false;
}
int main()
{
    ios_base::sync_with_stdio(false);
    cin.tie(0); cout.tie(0);
    int t;
    cin >> t;
    while(t--)
    {
        cin >> n;
        fo(i,1,n) cin >> a[i];
        Pre();
        if(Run())
        {
            cout << "YES" << '\n';
            cout << x << ' ' << y << ' ' << z << '\n';
        }
        else cout << "NO" << '\n';
    }
    return 0;
}