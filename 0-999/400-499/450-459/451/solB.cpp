#include <algorithm>
#include <bitset>
#include <cassert>
#include <climits>
#include <cmath>
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <iostream>
#include <iterator>
#include <limits>
#include <map>
#include <numeric>
#include <queue>
#include <set>
#include <sstream>
#include <stack>
#include <string>
#include <utility>
#include <vector>

using namespace std;

int main(){
  
    ios :: sync_with_stdio(false);
    int i,n,min=0,max;
    cin>>n;
    int a[n];
    max = n-1;
    for (int i = 0; i < n; ++i){
        cin>>a[i];
    }
    for (int i = 0; i < n-1; ++i)
    {
        if(a[i]<=a[i+1]){
            min = i;
        }else{
            min = i;
            break;
        }
    }
    //cout<<"min = "<<min<<endl;
    for (i = n-1; i >=min; --i)
    {
        if(a[i]>=a[i-1]){
            max = i;
        }
        else{
            max = i;
            break;
        }
        if(a[i]<a[min])
            break;
    }
    //cout<<"max = "<<max<<endl;
    for(i=min;i<max;i++){
        if(a[i]<a[i+1]){
            //cout<<a[i]<<" "<<a[i+1]<<endl;
            break;
        }
    }
    if(i == max){
        if(min>0 && a[i]>=a[min-1] || min == 0)
            cout<<"yes"<<endl<<min+1<<" "<<max+1;
        else
            cout<<"no";
    }
    else cout<<"no";
    return 0;
}