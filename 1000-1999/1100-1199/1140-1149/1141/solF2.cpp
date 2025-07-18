# include <bits/stdc++.h>

using namespace std;

const int N = 1505;

int n, a[N], dp[N * N], mx, b[N * N];
vector < pair <int, pair <int, int> > > v;
vector < pair <int, int> > ans;

int main(){
      cin >> n;

      for(int i = 1; i <= n; i ++){
            cin >> a[i];
      }

      for(int i = 1; i <= n; i ++){
            int sum = 0;
            for(int j = i; j <= n; j ++){
                  sum += a[j];
                  v.push_back({sum, {j, i}});
            }
      }

      sort(v.begin(), v.end());

      for(int i = 0; i < v.size(); i ++){
            int j = i;
            while(j + 1 < v.size() && v[j + 1].first == v[i].first)
                  j ++;
            for(int k = i; k <= j; k ++){
                  int f = v[k].second.second, s = v[k].second.first;
                  int lo = i, hi = k;
                  if(v[i].second.first >= f){
                        b[k] = 1;
                        dp[k] = 1;
                  } else {
                        while(hi - lo > 1){
                              int md = (lo + hi) >> 1;
                              if(v[md].second.first < f)
                                    lo = md;
                              else
                                    hi = md;
                        }
                        if(v[hi].second.first < f)
                              lo = hi;
                        dp[k] = dp[lo] + 1;
                        b[k] = dp[lo] + 1;
                  }
                  if(k != i)
                        dp[k] = max(dp[k], dp[k - 1]);
                  mx = max(mx, dp[k]);
            }
            i = j;
      }

      int mm = mx;

      for(int i = 0; i < v.size(); i ++){
            int j = i;
            while(j + 1 < v.size() && v[j + 1].first == v[i].first)
                  j ++;
            bool ok = 0;
            int mm = 2e9;
            for(int k = j; k >= i; k --){
                  if(b[k] == mx){
                        ok = 1;
                        int f = v[k].second.second, s = v[k].second.first;
                        if(s < mm){
                              ans.push_back({f, s});
                              mx --;
                              mm = f;
                        }
                  }
            }
            i = j;
            if(ok)
                  break;
      }


      cout << mm << endl;

      for(int i = 0; i < mm; i ++)
            cout << ans[i].first << " " << ans[i].second << endl;
}