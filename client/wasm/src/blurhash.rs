use blurhash::encode;
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