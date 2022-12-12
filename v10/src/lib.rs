extern crate alloc;
use std::collections::BTreeMap;
use value::{Secrets, Value};
use vrl::diagnostic::Formatter;
use vrl::Program;
use vrl::TimeZone;
use vrl::{state, Runtime, TargetValueRef};
use vrl_diagnostic::DiagnosticList;

use alloc::vec::Vec;
use std::mem::MaybeUninit;
use std::slice;

// Compiler

#[cfg_attr(all(target_arch = "wasm32"), export_name = "compile_vrl")]
#[no_mangle]
/// compile_vrl takes in a string representing a VRL program
/// Returns a 0 (null) if there was an error, otherwise
/// a pointer to the compiled program in linear memory
pub extern "C" fn compile_vrl(ptr: u32, len: u32) -> *mut Program {
    let program_string = unsafe { ptr_to_string(ptr, len) };

    match vrl::compile(&program_string, &vrl_stdlib::all()) {
        Ok(res) => return Box::into_raw(Box::new(res.program)),
        Err(err) => {
            // TODO return formatted string
            let f = Formatter::new(&program_string, err);
            panic!("{}", f.to_string());
            return std::ptr::null_mut();
        }
    }
}

unsafe fn ptr_to_string(ptr: u32, len: u32) -> String {
    let slice = slice::from_raw_parts_mut(ptr as *mut u8, len as usize);
    let utf8 = std::str::from_utf8_unchecked_mut(slice);
    return String::from(utf8);
}

// Runtime
#[cfg_attr(all(target_arch = "wasm32"), export_name = "new_runtime")]
#[no_mangle]
pub extern "C" fn new_runtime() -> *mut Runtime {
    Box::into_raw(Box::new(Runtime::new(state::Runtime::default())))
}

#[no_mangle]
pub extern "C" fn runtime_resolve(
    runtime: u32,
    program: u32,
    input_ptr: u32,
    input_len: u32,
) -> u64 {
    let rt = unsafe { (runtime as *mut Runtime).as_mut().unwrap() };
    let prog = unsafe { (program as *const Program).as_ref().unwrap() };
    let inpt = unsafe { ptr_to_string(input_ptr, input_len) };

    let mut value: Value = Value::from(inpt.as_str());
    let mut metadata = Value::Object(BTreeMap::new());
    let mut secrets = Secrets::new();
    let mut target = TargetValueRef {
        value: &mut value,
        metadata: &mut metadata,
        secrets: &mut secrets,
    };

    match rt.resolve(&mut target, &prog, &TimeZone::Local) {
        Ok(res) => {
            let s = res.to_string();
            return ((s.as_ptr() as u64) << 32) | s.len() as u64;
        }
        Err(_err) => {
            return 0;
        }
    }
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "runtime_clear")]
#[no_mangle]
pub extern "C" fn runtime_clear(runtime: u32) {
    let rt = unsafe { (runtime as *mut Runtime).as_mut().unwrap() };
    rt.clear()
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "runtime_is_empty")]
#[no_mangle]
pub extern "C" fn runtime_is_empty(runtime: u32) -> bool {
    let rt = unsafe { (runtime as *mut Runtime).as_mut().unwrap() };
    return rt.is_empty();
}

// WASM Memory-related helper functinos

/// WebAssembly export that allocates a pointer (linear memory offset) that can
/// be used for a string.
///
/// This is an ownership transfer, which means the caller must call
/// [`deallocate`] when finished.
#[cfg_attr(all(target_arch = "wasm32"), export_name = "allocate")]
#[no_mangle]
pub extern "C" fn allocate(size: u32) -> *mut u8 {
    // Allocate the amount of bytes needed.
    let vec: Vec<MaybeUninit<u8>> = Vec::with_capacity(size.try_into().unwrap());

    // into_raw leaks the memory to the caller.
    Box::into_raw(vec.into_boxed_slice()) as *mut u8
}

/// WebAssembly export that deallocates a pointer of the given size (linear
/// memory offset, byteCount) allocated by [`allocate`].
#[cfg_attr(all(target_arch = "wasm32"), export_name = "deallocate")]
#[no_mangle]
pub unsafe extern "C" fn _deallocate(ptr: u32, size: u32) {
    deallocate(ptr as *mut u8, size as usize);
}

/// Retakes the pointer which allows its memory to be freed.
unsafe fn deallocate(ptr: *mut u8, size: usize) {
    // TODO - should this be Box::from_raw? (see Box::into_raw docs)
    let _ = Vec::from_raw_parts(ptr, 0, size);
}
