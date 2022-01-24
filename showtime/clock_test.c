#include <stdio.h>
#include <stdint.h>
#include <time.h>
#include <errno.h>
#include <fcntl.h>
#include <linux/ptp_clock.h>
#include "all_clock.h"

int main()
{
    clockid_t clkid;
    clkid = utils_get_clockid("/dev/ptp1");

    time_ns sys_clk;
    time_ns ptp_clk;

    for(int i = 0; i < 10; i++)
    {
        sys_clk = utils_get_current_sys_ns_timestamp(); 
        ptp_clk = utils_get_current_ptp_ns_timestamp(clkid); 
        printf("%ld\n", ptp_clk - sys_clk - 37000000000);
    }
    return 0;
}
