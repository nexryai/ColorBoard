pub mod upload;
mod thumb;
mod enc;
mod blurhash;

use blurhash::decode_blurhash;
use console_error_panic_hook;
use wasm_bindgen::prelude::wasm_bindgen;
use wasm_bindgen::JsCast;

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

#[wasm_bindgen]
pub fn render_blurhash(element_id: String, hash: String) {
    //let hash: String = "LuNIK4?DI;aL~9o{NHwMt7Seofay".to_string();
    let pix = decode_blurhash(&hash);
    log(&format!("[DEBUG] decode blurhash: {:?}", pix));

    let document = web_sys::window().unwrap().document().unwrap();
    let canvas = document.get_element_by_id(&element_id).unwrap();
    let canvas: web_sys::HtmlCanvasElement = canvas
        .dyn_into::<web_sys::HtmlCanvasElement>()
        .map_err(|_| ())
        .unwrap();

    let context: web_sys::CanvasRenderingContext2d = canvas
        .get_context("2d")
        .unwrap()
        .unwrap()
        .dyn_into::<web_sys::CanvasRenderingContext2d>()
        .unwrap();

    
    // Vec<u8> を Uint8ClampedArray に変換
    let clamped_array = wasm_bindgen::Clamped(&pix[..]);
    let data = web_sys::ImageData::new_with_u8_clamped_array_and_sh(clamped_array, 150, 150).unwrap();

    let res = context.put_image_data(&data, 0.0, 0.0).unwrap();
}