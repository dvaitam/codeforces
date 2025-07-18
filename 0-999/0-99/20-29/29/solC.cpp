#define _CRT_SECURE_NO_WARNINGS
#include <stdio.h>
#include <string.h>
#include <string>
#include <math.h>
#include <algorithm>
#include <vector>
#include <queue>
using namespace std;
const int max_size=100003;
struct para{
    int a,b;
    bool operator < (para x) const{
        return a<x.a || a==x.a && b<x.b;
    }
};
vector <para> staff;
int n;

para find(para x){
    int l=0,r=(int)staff.size(),mid;
    r--;

    while(l<=r){
        mid = (l+r)/2;
        if(staff[mid].a> x.b)
            r=mid-1;
        else
            if(staff[mid].a<x.b)
                l=mid+1;
            else{
                if(staff[mid].b!=x.a)
                    return staff[mid];
                if( mid<((int)staff.size())-1 && staff[mid+1].a == x.b && staff[mid+1].b != x.a)
                    return staff[mid+1];
                return staff[mid-1];
            }
    }
    return  staff[mid];
}

int main(){
#ifdef  xDx
    freopen("input.txt","r",stdin);
    freopen("output.txt","w",stdout);
#endif
    scanf("%d",&n);
    for(int i=0; i < n ; i++){
        int k;
        para tmp;
        scanf("%d%d",&tmp.a,&tmp.b);
        staff.push_back(tmp);
        k=tmp.a;
        tmp.a = tmp.b;
        tmp.b = k;
        staff.push_back(tmp);
    }
    sort(staff.begin(),staff.end());
    int i;
    for(i=0 ; staff[i].a == staff[i+1].a; i+=2);
    para temp_city = staff[i];
    printf("%d ",temp_city.a);
    for(i =0 ;i< n; i++){
        if(i==n-1)
            printf("%d",temp_city.b);
        else{
            temp_city = find(temp_city);
            printf("%d ",temp_city.a);
        }
    }
    return 0;
}