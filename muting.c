#include <libavformat/avio.h>
#include <libavformat/avformat.h>

typedef struct {
    FILE *data;
} io;

int read_packet(void *, uint8_t *, int);
int write_packet(void *, uint8_t *, int);

void foo() {
    io *inputStream, *outputStream;

    const bufSize = 1024;
    
    unsigned char buf[bufSize];
    
    AVIOContext *input = avio_alloc_context(buf, bufSize, AVIO_FLAG_READ, inputStream, read_packet, NULL, NULL),
        *output = avio_alloc_context(buf, bufSize, AVIO_FLAG_WRITE, outputStream, NULL, write_packet, NULL);

    AVFormatContext *inputFormatCtx = avformat_alloc_context();
    
    int ret = av_probe_input_buffer2(inputStream, &(inputFormatCtx->iformat), "", inputFormatCtx, 0, 0);
    if (ret <= 0) {
        goto end;
    }

    AVFormatContext *outputFormatCtx = avformat_alloc_context();
    outputFormatCtx->oformat = av_guess_format(inputFormatCtx->iformat->name, "", "");
    if (!outputFormatCtx->oformat)
        goto end;

    

end:
    avformat_free_context(outputFormatCtx);
    avformat_free_context(inputFormatCtx);
    avio_context_free(&output);
    avio_context_free(&input);
}


int read_packet(void *opaque, uint8_t *buf, int buf_size) {
    io *input = opaque;

    return fread(buf, 1, buf_size, input->data);
}

int write_packet(void *opaque, uint8_t *buf, int buf_size) {
    io *input = opaque;

    return fwrite(buf, 1, buf_size, input->data);
}
