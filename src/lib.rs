extern crate alloc;
use std::collections::BTreeMap;
use std::ffi::CStr;
use std::ffi::CString;
use value::{Secrets, Value};
use vrl::diagnostic::Formatter;
use vrl::Program;
use vrl::TimeZone;
use vrl::{state, Runtime, TargetValueRef};

#[repr(C)]
pub struct CResult<T> {
    value: *mut T,
    error: *mut libc::c_char
}

// Compiler

#[no_mangle]
pub extern "C" fn compile_vrl(input: *const libc::c_char) -> CResult<Program> {
    let program_string = unsafe { CStr::from_ptr(input) }.to_str().unwrap();
    match vrl::compile(&program_string, &vrl_stdlib::all()) {
        Ok(res) => {
            return CResult { 
                value: Box::into_raw(Box::new(res.program)),
                error: std::ptr::null_mut()
            }
        }
        Err(err) => {
            let f = Formatter::new(program_string, err);
            return CResult { 
                value: std::ptr::null_mut(),
                error: CString::new(f.to_string()).unwrap().into_raw()
            }
        }
    }
}

// Runtime 

#[no_mangle]
pub extern "C" fn new_runtime() -> *mut Runtime {
    Box::into_raw(Box::new(Runtime::new(state::Runtime::default())))
}

#[no_mangle]
pub extern "C" fn runtime_resolve(runtime: *mut Runtime, program: *mut Program, input: *const libc::c_char) -> CResult<libc::c_char> {
    let rt = unsafe { runtime.as_mut().unwrap() };
    let prog = unsafe { program.as_ref().unwrap() };
    let inpt: &CStr = unsafe { CStr::from_ptr(input) };

    let mut value: Value = Value::from(inpt.to_str().unwrap());
    let mut metadata = Value::Object(BTreeMap::new());
    let mut secrets = Secrets::new();
    let mut target = TargetValueRef {
        value: &mut value,
        metadata: &mut metadata,
        secrets: &mut secrets,
    };

    match rt.resolve(&mut target, &prog, &TimeZone::Local) {
        Ok(res) => {
            return CResult {
                value: CString::new(res.to_string().as_bytes()).unwrap().into_raw(),
                error: std::ptr::null_mut()
            }
        }
        Err(err) => {
            return CResult { 
                value: std::ptr::null_mut(),
                error: CString::new(err.to_string()).unwrap().into_raw()
            }
        }
    }
}

#[no_mangle]
pub extern "C" fn runtime_clear(runtime: *mut Runtime) {
    let rt = unsafe { runtime.as_mut().unwrap() };
    rt.clear()
}

#[no_mangle]
pub extern "C" fn runtime_is_empty(runtime: *mut Runtime) -> bool {
    let rt = unsafe { runtime.as_mut().unwrap() };
    return rt.is_empty()
}