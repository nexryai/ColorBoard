use blurhash::{encode, decode};
use image::*;

use crate::log;

pub fn get_blurhash(data: &Vec<u8>) -> String {
    let image = image::load_from_memory(data).unwrap();
    let (width, height) = image.dimensions();

    let result = encode(4, 3, width, height, &image.to_rgba8().into_vec());
    let blurhash = match result {
        Ok(hash) => hash,
        Err(e) => {
            log(&format!("[Fatal] failed to generate blurhash: {}", e));
            panic!()
        },
    };

    blurhash
}

pub fn decode_blurhash(hash: &String) -> Vec<u8> {
    let result = decode(hash, 150, 150, 1.0);
    let pixels = match result {
        Ok(pix) => pix,
        Err(e) => {
            log(&format!("[Fatal] failed to decode blurhash: {}", e));
            panic!()
        },
    };

    pixels
}