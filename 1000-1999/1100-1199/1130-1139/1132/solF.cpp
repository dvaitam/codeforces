/*ЗАПУСКАЕМ 
░ГУСЯ░▄▀▀▀▄░РАБОТЯГИ░░ 
▄███▀░◐░░░▌░░░░░░░ 
░░░░▌░░░░░▐░░░░░░░ 
░░░░▐░░░░░▐░░░░░░░ 
░░░░▌░░░░░▐▄▄░░░░░ 
░░░░▌░░░░▄▀▒▒▀▀▀▀▄ 
░░░▐░░░░▐▒▒▒▒▒▒▒▒▀▀▄ 
░░░▐░░░░▐▄▒▒▒▒▒▒▒▒▒▒▀▄ 
░░░░▀▄░░░░▀▄▒▒▒▒▒▒▒▒▒▒▀▄ 
░░░░░░▀▄▄▄▄▄█▄▄▄▄▄▄▄▄▄▄▄▀▄ 
░░░░░░░░░░░▌▌░▌▌░░░░░ 
░░░░░░░░░░░▌▌░▌▌░░░░░ 
░░░░░░░░░▄▄▌▌▄▌▌░░░░░
НАСТРОЙКА НА КРИТЫ ██████████████] 100% СОЧНОСТИ*/
#include <iostream>
#include <vector>
#include <algorithm>
#include <set>
#include <map>
#include <string>
#include <math.h>
#include <queue>
#include <bitset>
#include <iomanip>
#include <cstring>
#include <cstdio>
#include <chrono>
#include <ctime>
#include <unordered_set>
#include <random>
  
#define ep emplace_back
#define F first
#define S second
#define ALL(a) (a).begin(), (a).end()
#define max_a(a) (max_element(a.begin(), a.end()))
  
using namespace std;
  
typedef long long ll;
typedef long double ld;
mt19937 rd(chrono :: system_clock :: now().time_since_epoch().count());
  
#pragma GCC optimize("unroll-loops") // развернуть цикл
#pragma GCC optimize("Ofast") // скомпилировать с о3
/*
#pragma GCC optimize("no-stack-protector") // что-то со стеком
#pragma GCC target("sse,sse2,sse3,ssse3,popcnt,abm,mmx,tune=native") // оптимизации процессора
#pragma GCC optimize("fast-math") // оптимизации сопроцессора
*/

const int inf = 1e9;

int main(){
   ios::sync_with_stdio(0), cin.tie(0);
#ifdef LOCAL
    freopen("input.txt", "r", stdin);
    freopen("output.txt", "w", stdout);
#endif

    int n;
    cin >> n;
    string s;
    cin >> s;

    string tmp;
    for (int i = 0; i < n; i++){
        if (tmp.empty() || tmp.back() != s[i])
            tmp += s[i];
    }
    s = tmp;
    n = s.size();

    int dp[n][n], i, j, k;
    for (i = 0; i < n; i++)
        for (j = 0; j < n; j++)
            dp[i][j] = inf;
    for (i = 0; i < n; i++)
        dp[i][i] = 1;

    for (i = n - 1; i >= 0; i--){
        for (j = i + 1; j < n; j++){
            dp[i][j] = dp[i + 1][j] + 1;
            for (k = i + 1; k <= j; k++){
                if (s[i] == s[k])
                    dp[i][j] = min(dp[i][j], dp[i + 1][k - 1] + dp[k][j]);
            }
        }
    }    
    
    cout << dp[0][n - 1];
}