#include <iostream>
#include <vector>
#include <cmath>
#include <cstdio>
#include <algorithm>
#include <vector>
#include <map>
#include <bits/stdc++.h>
using namespace std;


int main() {

    int games, wallets;
    int gamesarr[1000];
    cin >> games >> wallets;
    vector<int> wal;
    for (int i=0;i<games;i++) {
        cin >> gamesarr[i];


    }
    int x;
    for (int i=0;i<wallets;i++){
        cin >>x;
        wal.push_back(x);
    }
    int count=0;
    for (int i=0;i<games;i++){
        if (!wal.empty()){
            if (gamesarr[i]<=wal.front()){
                count++;
                wal.erase(wal.begin());
            }

        }
        else{
            break;
        }

    }
    cout << count;

    return 0;
}