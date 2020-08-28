// +build darwin

package v8

// #include <stdlib.h>
// #include <string.h>
// #include "v8_c_bridge.h"
// #cgo CXXFLAGS: -I${SRCDIR} -I${SRCDIR}/darwin/include -fno-rtti -fpic -std=c++11
// #cgo LDFLAGS: -pthread -L${SRCDIR}/darwin/libv8 -lv8_base -lv8_init -lv8_initializers -lv8_libbase -lv8_libplatform -lv8_libsampler -lv8_nosnapshot
import "C"

const LibV8LinkInfo = "static link from ${SRCDIR}/darwin/libv8"
