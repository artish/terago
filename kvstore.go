// +build !mock

package terago

/*
#cgo LDFLAGS: -ltera_c
#include "c/kvstore.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"sync"
	"unsafe"
)

type KvStore struct {
	Name   string
	CTable *C.tera_table_t
}

func (p KvStore) Close() {
	fmt.Println("close table: " + p.Name)
	if p.CTable != nil {
		C.tera_table_close(p.CTable)
	}
}

// ttl(time-to-live)
// Key-value will expired after <ttl> seconds. -1 means never expired.
func (p KvStore) Put(key, value string, ttl int) (err error) {
	if p.CTable == nil {
		return errors.New("table not open: " + p.Name)
	}
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	ckeylen := C.uint64_t(len(key))
	cvallen := C.uint64_t(len(value))
	ret := C.tera_table_put_kv(p.CTable, ckey, ckeylen, cvalue, cvallen, C.int32_t(ttl), nil)
	if !ret {
		err = errors.New("put kv error")
	}
	return
}

// Async put key-value into tera. Return success immediately and run put operation at background.
// Caution: If put failed, specify kv would be dump to error log.
func (p KvStore) PutAsync(key, value string, ttl int) (err error) {
	if p.CTable == nil {
		return errors.New("table not open: " + p.Name)
	}
	C.table_put_kv_async(p.CTable, C.CString(key), C.int(len(key)),
		C.CString(value), C.int(len(value)), C.int(ttl))
	return
}

func (p KvStore) Get(key string) (value string, err error) {
	if p.CTable == nil {
		err = errors.New("table not open: " + p.Name)
		return
	}
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	var vallen C.int
	cvalue := C.table_get_kv_sync(p.CTable, ckey, C.int(len(key)), (*C.int)(&vallen))
	if vallen >= 0 {
		value = C.GoStringN(cvalue, vallen)
		C.free(unsafe.Pointer(cvalue))
	} else {
		err = errors.New("key not found")
		value = ""
	}
	return
}

func (p KvStore) BatchPut(kvs []KeyValue) (err error) {
	wg := sync.WaitGroup{}
	wg.Add(len(kvs))
	succ := true
	for _, kvt := range kvs {
		kv := kvt
		go func() {
			if kv.TTL == 0 {
				kv.TTL = -1
			}
			kv.Err = p.Put(kv.Key, kv.Value, kv.TTL)
			if kv.Err != nil {
				succ = false
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if succ {
		return nil
	} else {
		return errors.New("error")
	}
}

func (p KvStore) BatchGet(keys []string) (result []KeyValue, err error) {
	wg := sync.WaitGroup{}
	wg.Add(len(keys))
	succ := true
	c := make(chan *KeyValue, len(keys))
	for _, kt := range keys {
		k := kt
		go func() {
			value, e := p.Get(k)
			if err != nil {
				c <- &KeyValue{Key: k, Err: e}
				succ = false
			} else {
				c <- &KeyValue{Key: k, Value: value}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(c)
	m := make(map[string]*KeyValue)
	for kv := range c {
		m[kv.Key] = kv
	}
	if len(m) != len(keys) {
		panic(m)
	}
	for _, k := range keys {
		result = append(result, *m[k])
	}
	if succ {
		return result, nil
	} else {
		return result, errors.New("error")
	}
	return
}

func (p KvStore) RangeGet(start, end string, maxNum int) (result []KeyValue, err error) {
	cstart := C.CString(start)
	defer C.free(unsafe.Pointer(cstart))
	cend := C.CString(end)
	defer C.free(unsafe.Pointer(cend))
	desc := C.tera_scan_descriptor(cstart, C.uint64_t(len(start)))
	C.tera_scan_descriptor_set_end(desc, cend, C.uint64_t(len(end)))
	defer C.tera_scan_descriptor_destroy(desc)
	scanner := C.tera_table_scan(p.CTable, desc, nil)
	defer C.tera_result_stream_destroy(scanner)
	for i := 0; i < maxNum; i++ {
		if C.tera_result_stream_done(scanner, nil) {
			break
		}
		var keylen, vallen C.int
		keyPtr := C.scanner_key(scanner, (*C.int)(&keylen))
		key := C.GoStringN(keyPtr, keylen)
		C.free(unsafe.Pointer(keyPtr))
		valPtr := C.scanner_value(scanner, (*C.int)(&vallen))
		value := C.GoStringN(valPtr, vallen)
		C.free(unsafe.Pointer(valPtr))
		C.tera_result_stream_next(scanner)

		result = append(result, KeyValue{Key: key, Value: value})
	}
	return
}

func (p KvStore) Delete(key string) (err error) {
	if p.CTable == nil {
		return errors.New("table not open: " + p.Name)
	}
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	ret := C.table_delete_kv_sync(p.CTable, ckey, C.int(len(key)))
	if !ret {
		err = errors.New("put kv error")
	}
	return
}
