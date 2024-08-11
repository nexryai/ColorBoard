use std::io::Cursor;
use uuid::Uuid;
use wasm_bindgen::prelude::wasm_bindgen;

use ehttp::multipart::MultipartBuilder;

use crate::{enc::encode_to_webp_lossless, thumb::generate_thumbnail};


#[wasm_bindgen]
extern "C" {
    #[wasm_bindgen(js_namespace = console)]
    fn log(s: &str);
}

#[wasm_bindgen]
pub fn upload_file(gallery_id: String, data: Vec<u8>) -> u16 {
    // Encode lossless image
    log("[ColorBoard WASM] Encoding lossless WebP");
    let mut lossless_data = encode_to_webp_lossless(data.clone());
    
    // Encode thumbnail
    log("[ColorBoard WASM] Encoding thumbnail WebP");
    let mut thumbnail_data = generate_thumbnail(data);

    let tmp_id = Uuid::new_v4();
    let filename = format!("{}.webp", tmp_id);
    let thumbnail_filename = format!("{}_thumb.webp", tmp_id);
    
    log("[ColorBoard WASM] Uploading lossless image");

    let request = ehttp::Request::multipart(
        format!("/api/gallery/{}/upload", gallery_id),
        MultipartBuilder::new()
            .add_text("test", "dummy")
            .add_stream(
                &mut Cursor::new(&mut lossless_data),
                &filename,
                Some(&filename),
                None,
            )
            .unwrap()
            .add_stream(
                &mut Cursor::new(&mut thumbnail_data),
                &thumbnail_filename,
                Some(&thumbnail_filename),
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