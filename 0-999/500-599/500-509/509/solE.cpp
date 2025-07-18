/*
    N.K.A
*/
#include <bits/stdc++.h>

#include <iostream>
#include <fstream>
#include <string>
#include <cmath>
#include <iomanip>
#include <vector>
#include <algorithm>
#include <set>
#include <map>
#include <deque>

using namespace std;

#define pb push_back
#define fi first
#define se second
#define in insert
int S[ 500009 ];
double ans[ 500009 ];


bool vowel (char p){
    if (p == 'I')return true;
    if (p == 'E')return true;
    if (p == 'A')return true;
    if (p == 'O')return true;
    if (p == 'Y')return true;
    if (p == 'U')return true;
    return false;
}

int main(){
    ios::sync_with_stdio(false);
    string s;
    cin >> s;
    int N = (int) s.length();

    for (int i = 1; i <= N ; i++)
    {
        if ( vowel(s[i - 1]) ) S[i] = 1;
        else S[i] = 0;
        S[i] += S[i - 1];
    }

    ans[1] = S[N];
    double ret = ans[1];

    for (int len = 2; len <= N ; len++)
    {

       ans[len] = (ans[len - 1] + S[N] - S[len - 1] -  (S[N] - S[N - len + 1]) );
       ret += ans[len] / ( (double) len);


    }
    cout << fixed << setprecision ( 8 ) ;

    cout << ret << endl;


    return 0;
}