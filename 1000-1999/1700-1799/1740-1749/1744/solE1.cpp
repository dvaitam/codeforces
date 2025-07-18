#include<iostream>

#include<bits/stdc++.h>

#include <cstddef>

#include<cstdint>

#include <map>

#include <algorithm>

#include<queue>

#include<unordered_map>

#include<unordered_set>

#define ll long long

#define ld long double

#define Omar_Said ios_base::sync_with_stdio(false); cin.tie(NULL); cout.tie(NULL);

using namespace std;

const int N = 2e5 + 5;



template<typename t>

t lcm(t a, t b) {

    return (a * b / __gcd(a, b));

}





/*

int longestCommonSubsequence(string str1,string str2,int m,int n,vector<vector<int> >& mumo)

{

    if(mumo[m][n]!=-1)

        return  mumo[m][n];

if(m==0||n==0)

    mumo[m][n]= 0;

else {

    if (str1[m - 1] == str2[n - 1])

        mumo[m][n] = 1 + longestCommonSubsequence(str1, str2, m - 1, n - 1, mumo);

    else

        mumo[m][n] = max(longestCommonSubsequence(str1, str2, m - 1, n, mumo),

                         longestCommonSubsequence(str1, str2, m, n - 1, mumo));

}

return mumo[m][n];

}

*/



void solve()

{

    ll a,b,c,d;

    cin>>a>>b>>c>>d;

    map<ll,ll>afacts;

    map<ll,ll>bfacts;

    for(ll i=1;i<=sqrt(a);i++)

    {

        if(a%i==0)

        {

            afacts[i]=a/i;

            afacts[a/i]=i;

        }

    }

    for(ll i=1;i<=sqrt(b);i++)

    {

        if(b%i==0)

        {

            bfacts[i]=b/i;

            bfacts[b/i]=i;

        }

    }

    for(auto i:afacts)

    {

        for(auto j:bfacts)

        {

            ll x=i.first*j.first;

            ll y=i.second*j.second;

            ll k1=a/x+1;

            ll k2 =b/y+1;

            if(k1*x<=c&&k2*y<=d)

            {

                cout<<k1*x<<' '<<k2*y<<endl;

                return;

            }

        }

    }

    cout<<-1<<' '<<-1<<endl;

}

int main() {

    Omar_Said;



    int t;

    cin>>t;

    while(t--)

    {

        solve();

    }



    return 0;

}

// arr[i]=arr[i-1]-> to transmit number of all elements less than index

// arr[x-1]-> all indices less than x and already including elements less than index

//ans+=arr[x-1]-> adds all of them

//arr[i]++-> at that index there exist an element less than index

//string tones[12]={"C","C#","D","D#","E", "F", "F#", "G", "G#", "A", "B", "H"};

/*for(int i=0;i<12;i++) {

    x = 11;

    for (int j = 11; j >= 0; j--) {

        arr[i][j] = x - i;

        x--;

        if (x - i < 0)

            x = 11 + i;

    }

}

*/