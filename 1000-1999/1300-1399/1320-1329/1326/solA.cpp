#include <bits/stdc++.h>
#define ll long long
#define ld long double
#define OK(x) cout<<(x?"YES\n":"NO\n")
#define Ok(x) cout<<(x?"Yes\n":"No\n")
#define ok(x) cout<<(x?"yes\n":"no\n")
#define ed '\n'
#define al(n) n.begin(), n.end()
#define all(n) n.rbegin(), n.rend()
#define YES cout << "YES" << endl
#define NO cout << "NO" << endl
#define speed ios_base::sync_with_stdio(0),cin.tie(0),cout.tie(0);
#define tc ll testcase;cin>>testcase;while(testcase--)

//                                " وَأَن لَّيْسَ لِلْإِنسَانِ إِلَّا مَا سَعَى "

template<typename T>
std::istream &operator>>(std::istream &input, std::vector<T> &data) {
    for (T &x: data)input >> x;
    return input;
}
template<typename T>
std::ostream &operator<<(std::ostream &output, const std::vector<T> &data) {
    for (const T &x: data)output << x << " ";
    return output;
}
using namespace std;
vector<ll> divisors(ll n)
{
    vector<ll>divs;
    for(int i=1;i*i<=n;i++){
        if (n % i == 0) {
            divs.push_back(i);
            if (i != n / i)divs.push_back(n / i);
        }
    }
    return divs;
}
void solve()
{
    int n;
    cin>>n;
    if(n>1)
    {
        cout<<2;
        while(n-->1)cout<<7;
    }else cout<<-1;
    cout<<ed;
}
int main(){

    speed

#ifndef ONLINE_JUDGE
    freopen("input.txt", "r", stdin);
    freopen("output.txt", "w", stdout);
#endif

    tc solve();

    return 0;
}