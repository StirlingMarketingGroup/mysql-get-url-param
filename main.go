/*
 * USAGE INSTRUCTIONS:
 *
 * make sure libmysqlclient-dev is installed:
 * apt install libmysqlclient-dev
 *
 * Replace "/usr/lib/mysql/plugin" with your MySQL plugins directory (can be found by running "select @@plugin_dir;")
 * go build -buildmode=c-shared -o get_url_param.so && cp get_url_param.so /usr/lib/mysql/plugin/get_url_param.so
 *
 * Then, on the server:
 * create function`get_url_param`returns string soname'get_url_param.so';
 *
 * And use/test like:
 * select`get_url_param`('https://stackoverflow.com/questions/51446087/how-to-debug-dump-go-variable-while-building-with-cgo?noredirect=1#comment89863750_51446087', 'noredirect'); -- outputs '1'
 *
 * Yeet!
 * Brian Leishman
 *
 */

package main

// #include <stdio.h>
// #include <sys/types.h>
// #include <sys/stat.h>
// #include <stdlib.h>
// #include <string.h>
// #include <mysql.h>
// #cgo CFLAGS: -I/usr/include/mysql -fno-omit-frame-pointer
import (
	"C"
)
import (
	"net/url"
	"unsafe"
)

//export get_url_param_deinit
func get_url_param_deinit(initid *C.UDF_INIT) {}

//export get_url_param_init
func get_url_param_init(initid *C.UDF_INIT, args *C.UDF_ARGS, message *C.char) C.my_bool {
	if args.arg_count != 2 {
		message = C.CString("`get_url_param` require 2 parameters: the URL string and the param name")
		return 1
	}

	argsTypes := *(*[2]uint32)(unsafe.Pointer(args.arg_type))

	argsTypes[0] = C.STRING_RESULT
	argsTypes[1] = C.STRING_RESULT

	initid.maybe_null = 1

	return 0
}

//export get_url_param
func get_url_param(initid *C.UDF_INIT, args *C.UDF_ARGS, result *C.char, length *uint64, isNull *C.char, message *C.char) *C.char {
	argsArgs := (*[1 << 30]*C.char)(unsafe.Pointer(args.args))[:2:2]

	a := make([]string, 2, 2)
	for i, argsArg := range argsArgs {
		a[i] = C.GoString(argsArg)
	}

	u, err := url.Parse(a[0])
	if err != nil {
		message = C.CString(err.Error())
		return nil
	}
	q := u.Query()

	v, ok := q[a[1]]
	if !ok || v == nil || len(v) == 0 {
		*length = 0
		*isNull = 1
		return nil
	}

	*length = uint64(len(v[0]))
	return C.CString(v[0])
}

func main() {}
