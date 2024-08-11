pub mod upload;
mod thumb;
mod enc;

use console_error_panic_hook;
use wasm_bindgen::prelude::wasm_bindgen;

#[wasm_bindgen]
extern "C" {
    #[wasm_bindgen(js_namespace = console)]
    fn log(s: &str);
}

#[wasm_bindgen(start)]
pub fn init() {
    log("Hook panic");
    std::panic::set_hook(Box::new(console_error_panic_hook::hook));
}