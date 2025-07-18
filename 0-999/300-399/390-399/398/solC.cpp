/*

            الحمد لله الذي لم يتخذ ولدا و لم يكن له شريك في الملك و لم يكن له ولي من الذل و كبره تكبيرا

                                                صل على رسول الله


                            من مر صدفة فليدع لي بالرضا و بحسن الخاتمة فضلا
 */
/*
 ***Please read the statement carefully
 ***Don't rush into the coding as long as you got a minor idea on the problem....
 ***Carefully analyze the problem, especially those number theory,dp, or greedy problems. Those ones hide vast worlds behind them
 ***There is no shame in bringing out a piece of paper and trying to derive a pattern or an equation for a problem even if it's an easy problem.
    That's surely much better than getting a wrong submission streak because of the rush.
 ***The cp world is just an alternate domain where you enjoy thinking. Don't make it a nightmare!
 ***Some problems are much easier than what they look. Before trying to look for an optimization, check the constraints, check the constraints,
    check the constraints!! ex. https://mirror.codeforces.com/contest/1807/problem/F
 ***Couldn't solve the problem? No worries, the problem is the one that lost the chance of you solving it :)
 ***Finally, cp is nothing without having a wife by your side, so pray for me to get one :)
 */

#include "bits/stdc++.h"

#define Alhamdulilah ios::sync_with_stdio(false);
#define OSSAMAs_domain cin.tie(0);cout.tie(0);
#define file_i freopen("input.txt", "r", stdin);
#define file_o freopen("output.txt", "w", stdout);

using namespace std;
using ll = long long;
void solve(ll tc) {
    //why would someone actually create such a problem?

    int n;cin>>n;
    if(n==5)return void(cout<<"1 2 6\n1 3 6\n2 4 5\n4 5 1\n3 4\n3 5\n");
    for(int i=1;i<=n>>1;i++)cout<<i<<' '<<i+n/2<<' '<<1<<'\n';
    for(int i=1;i+n/2<n;i++)cout<<i+n/2<<' '<<i+n/2+1<<' '<<2*i-1<<'\n';
    for(int i=1;i<n>>1;i++)cout<<i<<' '<<i+1<<'\n';
    cout<<"1 3\n";
}

signed main() {

    Alhamdulilah
    OSSAMAs_domain
    //  file_i
    //file_o
    int t = 1;
  //     cin >> t;

    for (int i = 1; i <= t; i++) {

        solve(i);
        fflush(stdout);
    }

    cerr << "Time elapsed: " << 1.0 * clock() / CLOCKS_PER_SEC << " s.\n";


    return 0;
}