use console_error_panic_hook;

use ehttp::multipart::MultipartBuilder;
use image::*;
use std::{io::Cursor, thread};
use uuid::Uuid;
use wasm_bindgen::prelude::wasm_bindgen;

use mime;

#[wasm_bindgen]
extern "C" {
    #[wasm_bindgen(js_namespace = console)]
    fn log(s: &str);
}

#[wasm_bindgen]
pub fn generate_thumbnail(data: Vec<u8>) -> Vec<u8> {
    log("[ColorBoard WASM] generating thumbinal...");

    // Using `image` crate, open the included .jpg file
    let img = image::load_from_memory(&data).unwrap();
    let (w, h) = img.dimensions();

    // 720pのサイズにリサイズ
    // 縦と横、超過が大きい方に合わせる
    let size_factor = 360.0 / (w.max(h) as f64);

    let img: DynamicImage = image::DynamicImage::ImageRgba8(imageops::resize(
        &img,
        (w as f64 * size_factor) as u32,
        (h as f64 * size_factor) as u32,
        imageops::FilterType::Triangle,
    ));

    // Create the WebP encoder for the above image
    let mut buf = Vec::new();
    img.write_to(&mut Cursor::new(&mut buf), ImageFormat::WebP)
        .unwrap();

    buf
}

#[wasm_bindgen]
pub fn upload_file(gallery_id: String, data: Vec<u8>) -> u16 {
    log("[ColorBoard WASM] Encoding lossless WebP");

    let img = image::load_from_memory(&data).unwrap();
    let mut buf = Vec::new();
    img.write_to(&mut Cursor::new(&mut buf), ImageFormat::WebP)
        .unwrap();

    let tmp_id = Uuid::new_v4();
    let filename = format!("{}.webp", tmp_id);
    
    log("[ColorBoard WASM] Uploading lossless image");

    let request = ehttp::Request::multipart(
        format!("/api/gallery/{}/upload", gallery_id),
        MultipartBuilder::new()
            .add_text("test", "dummy")
            .add_stream(
                &mut Cursor::new(&mut buf),
                &filename,
                Some(&filename),
                None,
            )
            .unwrap(),
    );

    log("[ColorBoard WASM] Uploading...");

    let status_code = std::sync::Arc::new(std::sync::Mutex::new(0));  // 共有可能な変数を作成
    let status_code_clone = std::sync::Arc::clone(&status_code);  // クローンを作成

    ehttp::fetch(request, move |response| {
        let mut status = status_code_clone.lock().unwrap();
        match response {
            Ok(response) => *status = response.status,
            Err(e) => {
                log(&format!("[Error] Failed to create multipart request: {}", e));
                *status = 0;
            }// エラーが発生した場合は0を設定
        }
        log("[ColorBoard WASM] Done");
    });

    let status = *status_code.lock().unwrap();  // ロックして値を取得

    status
}

#[wasm_bindgen(start)]
pub fn init() {
    log("Hook panic");
    std::panic::set_hook(Box::new(console_error_panic_hook::hook));
}