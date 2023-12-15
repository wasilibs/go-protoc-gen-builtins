package wasix_32v1

import (
	"context"
	"os"

	wazero "github.com/wasilibs/wazerox"
	"github.com/wasilibs/wazerox/api"
	"github.com/wasilibs/wazerox/experimental/sys"
)

const ModuleName = "wasix_32v1"

const i32, i64 = api.ValueTypeI32, api.ValueTypeI64

// MustInstantiate calls Instantiate or panics on error.
//
// This is a simpler function for those who know the module ModuleName is not
// already instantiated, and don't need to unload it.
func MustInstantiate(ctx context.Context, r wazero.Runtime) {
	if _, err := Instantiate(ctx, r); err != nil {
		panic(err)
	}
}

// Instantiate instantiates the ModuleName module into the runtime.
//
// # Notes
//
//   - Failure cases are documented on wazero.Runtime InstantiateModule.
//   - Closing the wazero.Runtime has the same effect as closing the result.
func Instantiate(ctx context.Context, r wazero.Runtime) (api.Closer, error) {
	return NewBuilder(r).Instantiate(ctx)
}

// Builder configures the ModuleName module for later use via Compile or Instantiate.
//
// # Notes
//
//   - This is an interface for decoupling, not third-party implementations.
//     All implementations are in wazero.
type Builder interface {
	// Compile compiles the ModuleName module. Call this before Instantiate.
	//
	// Note: This has the same effect as the same function on wazero.HostModuleBuilder.
	Compile(context.Context) (wazero.CompiledModule, error)

	// Instantiate instantiates the ModuleName module and returns a function to close it.
	//
	// Note: This has the same effect as the same function on wazero.HostModuleBuilder.
	Instantiate(context.Context) (api.Closer, error)
}

// NewBuilder returns a new Builder.
func NewBuilder(r wazero.Runtime) Builder {
	return &builder{r}
}

type builder struct{ r wazero.Runtime }

// hostModuleBuilder returns a new wazero.HostModuleBuilder for ModuleName
func (b *builder) hostModuleBuilder() wazero.HostModuleBuilder {
	ret := b.r.NewHostModuleBuilder(ModuleName)
	exportFunctions(ret)
	return ret
}

// Compile implements Builder.Compile
func (b *builder) Compile(ctx context.Context) (wazero.CompiledModule, error) {
	return b.hostModuleBuilder().Compile(ctx)
}

// Instantiate implements Builder.Instantiate
func (b *builder) Instantiate(ctx context.Context) (api.Closer, error) {
	return b.hostModuleBuilder().Instantiate(ctx)
}

func exportFunctions(builder wazero.HostModuleBuilder) {
	builder.NewFunctionBuilder().
		WithGoModuleFunction(callbackSignalFn, []api.ValueType{i32, i32}, []api.ValueType{}).
		Export("callback_signal")

	builder.NewFunctionBuilder().
		WithGoModuleFunction(chdirFn, []api.ValueType{i32, i32}, []api.ValueType{i32}).
		Export("chdir")

	builder.NewFunctionBuilder().
		WithGoModuleFunction(fdDupFn, []api.ValueType{i32, i32}, []api.ValueType{i32}).
		Export("fd_dup")

	builder.NewFunctionBuilder().
		WithGoModuleFunction(fdPipeFn, []api.ValueType{i32, i32}, []api.ValueType{i32}).
		Export("fd_pipe")

	builder.NewFunctionBuilder().
		WithGoModuleFunction(futexWaitFn, []api.ValueType{i32, i32, i32, i32}, []api.ValueType{i32}).
		Export("futex_wait")

	builder.NewFunctionBuilder().
		WithGoModuleFunction(futexWakeFn, []api.ValueType{i32, i32}, []api.ValueType{i32}).
		Export("futex_wake")

	builder.NewFunctionBuilder().
		WithGoModuleFunction(futexWakeAllFn, []api.ValueType{i32, i32}, []api.ValueType{i32}).
		Export("futex_wake_all")

	builder.NewFunctionBuilder().
		WithGoModuleFunction(getCWDFn, []api.ValueType{i32, i32}, []api.ValueType{i32}).
		Export("getcwd")

	builder.NewFunctionBuilder().
		WithGoModuleFunction(procExecFn, []api.ValueType{i32, i32, i32, i32}, []api.ValueType{}).
		Export("proc_exec")

	builder.NewFunctionBuilder().
		WithGoModuleFunction(procForkFn, []api.ValueType{i32, i32}, []api.ValueType{i32}).
		Export("proc_fork")

	builder.NewFunctionBuilder().
		WithGoModuleFunction(procJoinFn, []api.ValueType{i32, i32, i32}, []api.ValueType{i32}).
		Export("proc_join")

	builder.NewFunctionBuilder().
		WithGoModuleFunction(threadExitFn, []api.ValueType{i32}, []api.ValueType{}).
		Export("thread_exit")

	builder.NewFunctionBuilder().
		WithGoModuleFunction(threadIdFn, []api.ValueType{i32}, []api.ValueType{i32}).
		Export("thread_id")

	builder.NewFunctionBuilder().
		WithGoModuleFunction(threadParallelismFn, []api.ValueType{i32}, []api.ValueType{i32}).
		Export("thread_parallelism")

	builder.NewFunctionBuilder().
		WithGoModuleFunction(threadSignalFn, []api.ValueType{i32, i32}, []api.ValueType{i32}).
		Export("thread_signal")
}

var callbackSignalFn = api.GoModuleFunc(func(_ context.Context, _ api.Module, stack []uint64) {
	// We do not execute the wasm module concurrently so only have a single thread, we
	// can ignore signals.
	stack[0] = 0
})

var chdirFn = api.GoModuleFunc(func(_ context.Context, m api.Module, stack []uint64) {
	pathPtr := stack[0]
	pathLen := stack[1]
	path, ok := m.Memory().Read(uint32(pathPtr), uint32(pathLen))
	if !ok {
		panic("failed to read path")
	}
	if err := os.Chdir(string(path)); err != nil {
		panic(err)
	}
	stack[0] = 0
})

var fdDupFn = api.GoModuleFunc(func(_ context.Context, _ api.Module, _ []uint64) {
	// We do not support child plugins so never call this.
	panic("fd_dup")
})

var fdPipeFn = api.GoModuleFunc(func(_ context.Context, _ api.Module, _ []uint64) {
	// We do not support child plugins so never call this.
	panic("fd_pipe")
})

var futexWaitFn = api.GoModuleFunc(func(_ context.Context, _ api.Module, _ []uint64) {
	// We do not execute the wasm module concurrently so know this is never called.
	panic("futex_wait")
})

var futexWakeFn = api.GoModuleFunc(func(_ context.Context, _ api.Module, _ []uint64) {
	// We do not execute the wasm module concurrently so know this is never called.
	panic("futex_wake")
})

var futexWakeAllFn = api.GoModuleFunc(func(_ context.Context, _ api.Module, _ []uint64) {
	// We do not execute the wasm module concurrently so know this is never called.
	panic("futex_wake_all")
})

var getCWDFn = api.GoModuleFunc(func(_ context.Context, m api.Module, stack []uint64) {
	pathPtr := stack[0]
	pathLen := stack[1]
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if len(path) > int(pathLen)+1 {
		stack[0] = uint64(sys.ERANGE)
		return
	}

	m.Memory().Write(uint32(pathPtr), []byte(path))
	m.Memory().WriteByte(uint32(pathPtr)+uint32(len(path)), 0)
	stack[0] = 0
})

var procExecFn = api.GoModuleFunc(func(_ context.Context, _ api.Module, _ []uint64) {
	// We do not support child plugins so never call this.
	panic("proc_exec")
})

var procForkFn = api.GoModuleFunc(func(_ context.Context, _ api.Module, _ []uint64) {
	// We do not support child plugins so never call this.
	panic("proc_fork")
})

var procJoinFn = api.GoModuleFunc(func(_ context.Context, _ api.Module, _ []uint64) {
	// We do not support child plugins so never call this.
	panic("proc_join")
})

var threadExitFn = api.GoModuleFunc(func(_ context.Context, _ api.Module, _ []uint64) {
	// We do not execute the wasm module concurrently so know this is never called.
	panic("thread_exit")
})

var threadIdFn = api.GoModuleFunc(func(_ context.Context, m api.Module, stack []uint64) {
	// We do not execute the wasm module concurrently so return an arbitrary value.
	resPtr := uint32(stack[0])
	m.Memory().WriteUint32Le(resPtr, 1)
	stack[0] = 0
})

var threadParallelismFn = api.GoModuleFunc(func(_ context.Context, m api.Module, stack []uint64) {
	// We do not execute the wasm module concurrently so force this to 1, as if 1 CPU.
	resPtr := uint32(stack[0])
	m.Memory().WriteUint32Le(resPtr, 1)
	stack[0] = 0
})

var threadSignalFn = api.GoModuleFunc(func(_ context.Context, _ api.Module, stack []uint64) {
	// We do not execute the wasm module concurrently so only have a single thread, we
	// can ignore signals.
	stack[0] = 0
})
