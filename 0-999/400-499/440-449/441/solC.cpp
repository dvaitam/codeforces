#include<bits/stdc++.h>

using namespace std;







int main() {

    int n,m,k;

    scanf("%d%d%d",&n,&m,&k);

    int noOfCells=n*m/k;

    int CellsLeft=(n*m)%k;

    int i=0,j=0;

    int co=0;

    if(k!=1) {

        printf("%d ",noOfCells);

        for(i=0;i<n;i++) {

            if(!(i&1)) {

                for(j=0;j<m;j++) {

                    printf("%d %d ",i+1,j+1);

                    co++;

                    if(co%noOfCells==0) {

                        k--;

                        puts("");

                        if(k==1) {

                            j++;

                            break;

                        }

                        printf("%d ",noOfCells);

                    }

                }

            }else {

                for(j=m-1;j>=0;j--) {

                    printf("%d %d ",i+1,j+1);

                    co++;

                    if(co%noOfCells==0) {

                        k--;

                        puts("");

                        if(k==1) {

                            j--;

                            break;

                        }

                        printf("%d ",noOfCells);

                    }

                }

            }

            if(k==1) {

                //i++;

                break;

            }

        }



    }

    if(j<0||j>=m) {

        i++;

        if(j<0) {

            j=0;

        }else {

            j=m-1;

        }

    }

    printf("%d ",noOfCells+CellsLeft);

    for(;i<n;i++) {

        if(j<0) {

            j=0;

        }else if(j>=m) {

            j=m-1;

        }

        if(i&1) {

            for(;j>=0;j--) {

                printf("%d %d ",i+1,j+1);

            }

        }else {

            for(;j<m;j++) {

                printf("%d %d ",i+1,j+1);

            }

        }

    }



}