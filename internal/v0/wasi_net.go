package v0

import (
	"fmt"

	"github.com/bytecodealliance/wasmtime-go/v13"
	"github.com/gaukas/water/internal/wasm"
)

// WASIConnectFunc is a function that creates a connection.
type WASIConnectFunc = func(caller *wasmtime.Caller) (fd int32, err error)

// WASIConnectFuncType is the signature of WASIConnectFunc.
var WASIConnectFuncType *wasmtime.FuncType = wasmtime.NewFuncType(
	[]*wasmtime.ValType{},
	[]*wasmtime.ValType{
		wasmtime.NewValType(wasmtime.KindI32), // return: connectionFd
	},
)

// WrapConnectFunc wraps a WASIConnectFunc into a WASM function.
func WrapConnectFunc(f WASIConnectFunc) wasm.WASMTIMEStoreIndependentFunction {
	return func(caller *wasmtime.Caller, vals []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
		if len(vals) != 0 {
			return []wasmtime.Val{wasmtime.ValI32(wasm.INVALID_ARGUMENT)}, wasmtime.NewTrap(fmt.Sprintf("v0.WASIConnectFunc expects 0 argument, got %d", len(vals)))
		}

		fd, err := f(caller)
		if err != nil { // here fd is expected to be an error code (negative)
			return []wasmtime.Val{wasmtime.ValI32(fd)}, wasmtime.NewTrap(fmt.Sprintf("v0.WASIConnectFunc: %v", err))
		}

		return []wasmtime.Val{wasmtime.ValI32(fd)}, nil
	}
}

// UnimplementedWASIConnectFunc wraps unimplementedWASIConnectFunc into a
// wasmtime-compliant function.
func WrappedUnimplementedWASIConnectFunc() wasm.WASMTIMEStoreIndependentFunction {
	return WrapConnectFunc(unimplementedWASIConnectFunc)
}

// unimplementedWASIConnectFunc is a WASIConnectFunc that does nothing.
func unimplementedWASIConnectFunc(_ *wasmtime.Caller) (fd int32, err error) {
	return wasm.INVALID_FUNCTION, fmt.Errorf("NOP WASIConnectFunc is called")
}
