package enst

import "reflect"

func indirect(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

// func isNill(v reflect.Value) bool {
// 	// if vv is not ptr, return false(v is not nil)
// 	// if vv is ptr, return v.IsNil()
// 	return v.Kind() == reflect.Pointer && v.IsNil()
// }

// func pkgPath(t reflect.Type) string {
// 	pkg := t.PkgPath()
// 	if pkg != "" {
// 		return pkg
// 	}
// 	switch t.Kind() {
// 	case reflect.Slice, reflect.Array, reflect.Ptr, reflect.Map:
// 		return pkgPath(t.Elem())
// 	}
// 	return pkg
// }
