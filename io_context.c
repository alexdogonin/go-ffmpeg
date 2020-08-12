#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>
#include "io_context.h"
// #include <_cgo_export.h>
#include <libavformat/avio.h>


// int io_copy(int (*r)(void* opaque, uint8_t *buf, uint8_t buf_size), writer wr) {
    
//     uint8_t buf[8];
    
//     int count;
//     for(;;) {
//         int ret = r(NULL, buf, 8);
//         if (ret <= 0)
//             return ret;

//         fprintf(stdout, "read %d bytes", ret);
        
//         ret = wr(NULL, buf, ret);
//         if (ret <= 0)
//             return ret;

//         fprintf(stdout, "write %d bytes", ret);

//         count += ret;
//     }

//     return count;
// }

