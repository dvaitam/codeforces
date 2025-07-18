#include <iostream>
#include <string>
#include <cmath>
#include <algorithm>
#include <stack>
#include <vector>
#include <queue>
#include <map>
#include <iomanip>
#define ll long long
#define loop(n) for (int i = 0; i < n; i++)
#define rloop(n) for (int i = 1; i < n; i++)
#define iloop(n) for (int i = n - 1; i >= 0; i--)
#define nloop(n) for (int j = 0; j < n; j++)
#define jloop(n) for (int j = n - 1; j >= 0; j--)
#define srt(v) sort(v.begin(), v.end())
#define srtg(v) sort(v.begin(), v.end(), greater<ll>());
using namespace std;
bool isprim(ll n)
{
    for (ll i = 2; i * i <= n; i++)
        if (n % i == 0)
            return false;
    return true;
}
ll factorial(ll n)
{
    if (n == 1 || n == 2)
        return n;
    else
        return factorial(n - 1) * n;
}
ll pow(int a, int b)
{
    if (b == 0)
        return 1;
    int sq = pow(a, b / 2);
    sq *= sq;
    if (b % 2 != 0)
        sq *= a;
    return sq;
}
ll lcm(int hp,int tp){
        int n1=hp; int n2=tp; int lcm;
        while (hp != tp) {
            if (hp > tp)
                hp -= tp;
            else
                tp -= hp;
        }
     return   lcm = (n1 * n2) / hp;
}
int main(){
    ios::sync_with_stdio(false);
    cin.tie(0);
    cout.tie(0);
    int t;
    cin>>t;
    while(t--){
        int l,r;
        cin>>l>>r;
        if (r%l==0) cout<<l<<" "<<r<<"\n";
        else if (r%2==0&&r/2>=l) cout<<r/2<<" "<<r<<"\n";
        else if ((r-1)%2==0&&(r-1)/2>=l) cout<<(r-1)/2<<" "<<r-1<<"\n";
        else cout<<"-1 -1\n";
        } 
 }