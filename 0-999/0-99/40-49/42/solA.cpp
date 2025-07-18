#include <bits/stdc++.h>

using namespace std;

#define MAX 20



int _ratio[MAX+3];

int quantity[MAX+3];

int n,V;



double check(double x) {

 double sum = 0;

 for(int i = 1; i <= n; i++) {

    if((double)_ratio[i]*x <= (double)quantity[i]) {

        sum += (double)(_ratio[i])*x;

    }

    else return -1;

    if((double)sum > V) return -1;

 }

 return sum;

}



int main(void) {

  int v;

  cin>>n>>V;

  for(int i = 1; i <=n; i++){

    cin>>_ratio[i];

  }

  double low = 0, high = 0, mid;

  for(int i = 1; i <= n; i++) {

        cin>>quantity[i];

        high = max(high,(double)quantity[i]);

  }

  int counter = 0;

  double res = 0;

  while(counter<= 50) {

    mid = (low+high)/2.0;

    double ans = check(mid);

    if(ans == -1) {

        high = mid;

    }

    else {

        low = mid;

        res = max(res,ans);

    }

    counter++;

  }

  cout<<res<<endl;

  return 0;

}