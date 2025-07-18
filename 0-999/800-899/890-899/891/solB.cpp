#include<bits/stdc++.h>
#define int long long
#define dd second
#define mp make_pair
#define pb push_back
#define ff first
#define dd second
#define pp pair<int,int>
using namespace std;
int tab[23];
int n;

int sorted[23];

main()
{

    ios_base::sync_with_stdio(false);
    cin >> n;
    for(int i = 0; i < n; ++i) cin >> tab[i];
    vector<int> k;
    for(int i = 0; i < n; ++i) k.pb(i);
//    do
//    {
//        random_shuffle(k.begin(),k.end());
//        //for(int i = 0; i < k.size(); ++i) cout << k[i] << " "; cout << endl;
//        bool czy = false;
//        for(int mask =1; mask < ((1 << n)-1); ++mask)
//        {
//            int a = 0; int b = 0;
//            for(int bit = 0; bit < n; ++bit)
//            {
//                if((mask >> bit) & 1)
//                {
//                    a += tab[bit]; b += tab[k[bit]];
//                }
//            }
//            if(a == b)
//            {
//                czy = true; break;
//            }
//        }
//        if(czy) continue;
//        for(int i = 0; i < n; ++i)
//        {
//            cout << tab[k[i]] << " ";
//        }
//        cout << endl;
//        exit(0);
//
//    }while(2);
    for(int i = 0; i < n; ++i)
    {
        sorted[i] = tab[i];
    }
    sort(sorted, sorted+n);
    for(int i = 0; i < n; ++i)
    {
        for(int j = 0; j < n; ++j)
        {
            if(tab[i] == sorted[j]) cout << sorted[(j+1) % n] << " ";
        }
    }
    cout << endl;
}