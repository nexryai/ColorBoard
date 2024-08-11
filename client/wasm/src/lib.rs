use std::io::Cursor;
use uuid::Uuid;
use image::*;
use wasm_bindgen::prelude::wasm_bindgen;
use ehttp::{Request, Response};
use ehttp::multipart::MultipartBuilder;

use mime;


#[wasm_bindgen]
pub fn generate_thumbnail(data: Vec<u8>) -> Vec<u8> {
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
    img.write_to(&mut Cursor::new(&mut buf), ImageFormat::WebP).unwrap();

    buf
}

#[wasm_bindgen]
pub fn upload_file(gallery_id: String, mut data: Vec<u8>) -> u16 {
    let tmp_id = Uuid::new_v4();
    let filename = format!("{}.webp", tmp_id);

    let request = ehttp::Request::multipart(
        format!("/api/gallery/{}/upload", gallery_id),
        MultipartBuilder::new()
            .add_text("test", "dummy")
            .add_stream(
                &mut Cursor::new(&mut data),
                &filename,
                Some(&filename),
                // FIXME: https://github.com/hyperium/mime/pull/129
                Some(mime::IMAGE_PNG)
            )
            .unwrap(),
    );

    let (sender, receiver) = std::sync::mpsc::channel();
    ehttp::fetch(request, move |response| {
        match response {
            Ok(response) => sender.send(response.status).unwrap(),
            Err(_) => sender.send(0).unwrap(), // エラーが発生した場合は0を返す
        }
    });

    receiver.recv().unwrap()
}