package image

/*
#include <stdlib.h>
#include <stdio.h>
#include "jpeglib.h"
#include "jerror.h"

void calc_hist(unsigned char *jpegFile, size_t size, ulong r[], ulong g[], ulong b[]) {
	int row_stride;
	struct jpeg_error_mgr jerr;
	struct jpeg_decompress_struct cinfo;

	cinfo.err = jpeg_std_error(&jerr);
	jpeg_create_decompress(&cinfo);
	jpeg_mem_src(&cinfo, jpegFile, size);
	jpeg_read_header(&cinfo, TRUE);
	jpeg_start_decompress(&cinfo);
	row_stride = cinfo.output_width * cinfo.output_components;

	int i;
	JSAMPARRAY buffer;
	buffer = (cinfo.mem->alloc_sarray)((j_common_ptr)&cinfo, JPOOL_IMAGE, row_stride, 1);
	while (cinfo.output_scanline < cinfo.output_height) {
		jpeg_read_scanlines(&cinfo, buffer, 1);
		for (i = 0; i < row_stride; i += 3) {
			r[buffer[0][i]]++;
			g[buffer[0][i+1]]++;
			b[buffer[0][i+2]]++;
		}
	}

	jpeg_finish_decompress(&cinfo);
	jpeg_destroy_decompress(&cinfo);
}
*/
// #cgo LDFLAGS: -ljpeg
import "C"
import (
	"io"
	"io/ioutil"
	"unsafe"
)

func decodeJpeg(reader io.Reader) (*Histogram, error) {
	// TODO: use source reading the file as needed?
	// https://github.com/pixiv/go-libjpeg/blob/master/jpeg/sourceManager.go#L133
	// https://cs.stanford.edu/~acoates/decompressJpegFromMemory.txt
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var hist Histogram
	C.calc_hist((*C.uchar)(unsafe.Pointer(&bytes[0])), C.size_t(len(bytes)), (*C.ulong)(&hist.Red[0]), (*C.ulong)(&hist.Green[0]), (*C.ulong)(&hist.Blue[0]))
	return &hist, nil
}
