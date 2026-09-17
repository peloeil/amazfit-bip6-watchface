# /// script
# requires-python = ">=3.10"
# dependencies = [
#     "Pillow==12.3.0",
#     "zopfli==0.4.3",
# ]
# ///

"""Exact source RGBA encoding and lossless standard ZIP compression search.
Only the exact L8/ARGB8888 format emitted for this target is accepted.
Native headers, dimensions and padding stay intact. Visible pixels are encoded
directly from the source PNG, avoiding the compiler's lossy palette conversion.
Configuration, expiry, signature and bytecode are never modified here.
"""
from pathlib import Path
from functools import lru_cache
from collections import Counter
import argparse
import copy
import hashlib
import io
import json
import struct
import zipfile
import zlib
from PIL import Image

import zopfli.zlib


def image_parts(data):
    if (len(data) < 1088 or
        data[:12] != bytes.fromhex('2e0101000000012000000000') or
        data[16:22] != bytes.fromhex('0820534f4d48') or
        data[24:28] != bytes.fromhex('00000100')):
        raise ValueError('Unexpected native image format; refusing to rewrite')
    stride, height = struct.unpack_from('<HH', data, 12)
    width = struct.unpack_from('<H', data, 22)[0]
    if not (0 < width <= stride and stride % 16 == 0 and height > 0 and
            len(data) == 1088 + stride*height):
        raise ValueError('Unexpected native image dimensions/length')
    palette = [data[p:p+4] for p in range(64, 1088, 4)]
    return palette, data[1088:]


def decoded_image(data):
    palette, indices = image_parts(data)
    return data[:64] + b''.join(palette[index] for index in indices)


def normalize_image(data, order):
    palette, indices = image_parts(data)
    counts = Counter(palette[index] for index in indices)
    colors = (sorted(counts, key=lambda c: (-counts[c], c))
              if order == 'frequency' else sorted(counts))
    positions = {color: index for index, color in enumerate(colors)}
    remap = bytes(positions.get(color, 0) for color in palette)
    result = data[:64] + b''.join(colors).ljust(1024, b'\0') + indices.translate(remap)
    if decoded_image(data) != decoded_image(result):
        raise ValueError('Native pixel or header changed')
    return result


def encode_source_image(data, source_path):
    image_parts(data)  # Strict format and length validation before any rewrite.
    stride, height = struct.unpack_from('<HH', data, 12)
    width = struct.unpack_from('<H', data, 22)[0]
    with Image.open(source_path) as source:
        if source.size != (width, height):
            raise ValueError('Source and native image dimensions differ')
        source_rgba = source.convert('RGBA').tobytes()
    rgba = bytearray(decoded_image(data)[64:])
    for y in range(height):
        rgba[y*stride*4:(y*stride+width)*4] = source_rgba[y*width*4:(y+1)*width*4]
    pixels = [bytes(rgba[i:i+4]) for i in range(0, len(rgba), 4)]
    colors = sorted(set(pixels))
    if len(colors) > 256:
        raise ValueError('Exact source colors exceed native palette capacity')
    positions = {color: index for index, color in enumerate(colors)}
    result = (data[:64] + b''.join(colors).ljust(1024, b'\0') +
              bytes(positions[color] for color in pixels))
    if decoded_image(result) != data[:64] + rgba:
        raise ValueError('Source RGBA encoding mismatch')
    return result


def encode_source_assets(original, assets_dir):
    outer = []
    for info, inner in read_zip(original):
        entries = []
        for item, data in read_zip(inner):
            if is_image(item.filename):
                source = (assets_dir/item.filename.removeprefix('assets/')).resolve()
                if not source.is_relative_to(assets_dir.resolve()):
                    raise ValueError('Image path leaves the supplied asset directory')
                data = encode_source_image(data, source)
            entries.append((item, data))
        outer.append((info, write_zip(entries, zipfile.ZIP_STORED)))
    return write_zip(outer, zipfile.ZIP_DEFLATED)


@lru_cache(maxsize=256)
def best_deflate(data, iterations=30):
    a = zopfli.zlib.compress(data, numiterations=iterations)[2:-4]
    b = zlib.compress(data, 9)[2:-4]
    return a if len(a) < len(b) else b


class BestDeflate:
    def __init__(self, iterations=30):
        self.parts = []
        self.iterations = iterations

    def compress(self, data):
        self.parts.append(data)
        return b''

    def flush(self):
        return best_deflate(b''.join(self.parts), self.iterations)


def read_zip(data):
    with zipfile.ZipFile(io.BytesIO(data)) as archive:
        names = archive.namelist()
        if len(names) != len(set(names)):
            raise ValueError('Duplicate ZIP entry')
        return [(info, archive.read(info)) for info in archive.infolist()]


def write_zip(entries, method):
    out = io.BytesIO()
    with zipfile.ZipFile(out, 'w') as archive:
        for info, data in entries:
            clean = copy.copy(info)
            clean.compress_type = method
            archive.writestr(clean, data)
    return out.getvalue()


def is_image(name):
    return name.startswith('assets/') and name.endswith('.png')


def fingerprints(data, decoded=False):
    results = {}
    for info, inner in read_zip(data):
        files = {}
        for item, payload in read_zip(inner):
            if decoded and is_image(item.filename):
                payload = decoded_image(payload)
            files[item.filename] = hashlib.sha256(payload).hexdigest()
        results[info.filename] = files
    return results


def optimize(original, baseline_label='official'):
    outer = read_zip(original)
    if {info.filename for info, _ in outer} != {'device.zip', 'app-side.zip'}:
        raise ValueError('Unexpected ZPK structure')
    reference = fingerprints(original, decoded=True)
    candidates = [(baseline_label, original)]
    compressor = zipfile._get_compressor
    zipfile._get_compressor = lambda method, level=None: (
        BestDeflate() if method == zipfile.ZIP_DEFLATED else compressor(method, level))
    try:
        for order in (None, 'sorted', 'frequency'):
            for method in (zipfile.ZIP_STORED, zipfile.ZIP_DEFLATED):
                entries = []
                for info, inner in outer:
                    contents = []
                    for item, payload in read_zip(inner):
                        if order and is_image(item.filename):
                            payload = normalize_image(payload, order)
                        contents.append((item, payload))
                    entries.append((info, write_zip(contents, method)))
                result = write_zip(entries, zipfile.ZIP_DEFLATED)
                if reference != fingerprints(result, decoded=True):
                    raise ValueError('Payload changed beyond palette encoding')
                label = f'{order or "original"}-inner-{method}'
                candidates.append((label, result))
                print(f'{label}: {len(result)} B', flush=True)
        # ZIP member order is independent of rendering order. Group repeated
        # image data within DEFLATE's dictionary window without runtime work.
        orderings = {
            'lexical': lambda entry: entry[0].filename,
            'size-ascending': lambda entry: (len(entry[1]), entry[0].filename),
            'preview-first': lambda entry: (entry[0].filename != 'assets/preview.png', entry[0].filename),
            'preview-last': lambda entry: (entry[0].filename == 'assets/preview.png', entry[0].filename),
        }
        for label, key in orderings.items():
            entries = []
            for info, inner in outer:
                contents = [(item, normalize_image(data, 'sorted') if is_image(item.filename) else data)
                            for item, data in read_zip(inner)]
                if info.filename == 'device.zip':
                    contents.sort(key=key)
                entries.append((info, write_zip(contents, zipfile.ZIP_STORED)))
            result = write_zip(entries, zipfile.ZIP_DEFLATED)
            if reference != fingerprints(result, decoded=True):
                raise ValueError('ZIP reordering changed a payload')
            candidates.append((label, result))
            print(f'{label}: {len(result)} B', flush=True)
        label, result = min(candidates, key=lambda candidate: len(candidate[1]))
        zipfile._get_compressor = lambda method, level=None: (
            BestDeflate(120) if method == zipfile.ZIP_DEFLATED else compressor(method, level))
        deeper = write_zip(read_zip(result), zipfile.ZIP_DEFLATED)
        if reference != fingerprints(deeper, decoded=True):
            raise ValueError('Deeper compression changed a payload')
        candidates.append((label+'-deep', deeper))
        print(f'{label}-deep: {len(deeper)} B', flush=True)
    finally:
        zipfile._get_compressor = compressor
    label, result = min(candidates, key=lambda candidate: len(candidate[1]))
    return result, label, {name: len(data) for name, data in candidates}



def decode_native(data):
    """Independently decode SOMH indices with Pillow using Zepp RGBA tuples."""
    image_parts(data)
    image = Image.open(io.BytesIO(data))
    if image.format != 'TGA' or image.mode != 'P' or image.palette.rawmode != 'BGRA':
        raise ValueError('Unexpected Pillow image interpretation')
    image.palette.mode = 'RGBA'
    image.palette.rawmode = 'RGBA'
    return image.convert('RGBA')


def verify_release(original, result, assets_dir):
    before = {info.filename: dict((i.filename, b) for i, b in read_zip(inner))
              for info, inner in read_zip(original)}
    after = {info.filename: dict((i.filename, b) for i, b in read_zip(inner))
             for info, inner in read_zip(result)}
    if before.keys() != after.keys():
        raise ValueError('Archive members changed')
    checked = 0
    for archive, files in before.items():
        if files.keys() != after[archive].keys():
            raise ValueError('Inner ZIP members changed')
        for name, data in files.items():
            actual = after[archive][name]
            if not is_image(name):
                if actual != data:
                    raise ValueError('Non-image payload changed: '+name)
                continue
            if data[:64] != actual[:64] or len(data) != len(actual):
                raise ValueError('Native image header or allocation changed: '+name)
            original_image, final_image = decode_native(data), decode_native(actual)
            width = struct.unpack_from('<H', actual, 22)[0]
            source_path = (assets_dir/name.removeprefix('assets/')).resolve()
            if not source_path.is_relative_to(assets_dir.resolve()):
                raise ValueError('Image path leaves the supplied asset directory')
            with Image.open(source_path) as png:
                source = png.convert('RGBA')
                if source.size != (width, final_image.height):
                    raise ValueError('Source dimensions differ: '+name)
                if final_image.crop((0,0,width,final_image.height)).tobytes() != source.tobytes():
                    raise ValueError('Source RGBA differs: '+name)
            padding = (width,0,final_image.width,final_image.height)
            if original_image.crop(padding).tobytes() != final_image.crop(padding).tobytes():
                raise ValueError('Native row padding changed: '+name)
            checked += 1
    return checked, after


def read_input(path):
    if path.is_dir():
        bundles = sorted(path.glob('*.zab'))
        if len(bundles) != 1:
            raise ValueError('Specify one .zab when the input directory has zero or multiple bundles')
        path = bundles[0]
    payload = path.read_bytes()
    entries = dict((i.filename, b) for i, b in read_zip(payload))
    if set(entries) == {'device.zip','app-side.zip'}:
        return payload
    manifest = json.loads(entries['manifest.json'])['zpks']
    if len(manifest) != 1:
        raise ValueError('Select a single-target ZAB')
    item = manifest[0]
    original = entries[item['name']]
    if hashlib.md5(original).hexdigest() != item['md5']:
        raise ValueError('ZAB manifest MD5 mismatch')
    return original


def fixed_optimize(source_exact):
    compressor = zipfile._get_compressor
    zipfile._get_compressor = lambda method, level=None: (
        BestDeflate(120) if method == zipfile.ZIP_DEFLATED else compressor(method,level))
    try:
        result = write_zip(read_zip(source_exact),zipfile.ZIP_DEFLATED)
    finally:
        zipfile._get_compressor = compressor
    if fingerprints(source_exact,decoded=True) != fingerprints(result,decoded=True):
        raise ValueError('Compression changed a payload')
    return result


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('input',type=Path,help='Standard .zab/.zpk, or a directory containing one .zab')
    parser.add_argument('output',type=Path,help='Optimized .zpk path')
    parser.add_argument('--source',type=Path,required=True,help='Matching rebuild source directory')
    parser.add_argument('--fixed',action='store_true',help='Use the selected original-order/STORE/Zopfli-120 preset')
    args = parser.parse_args()
    if args.input.resolve() == args.output.resolve():
        parser.error('Use distinct input and output paths')
    source_cfg = json.loads((args.source/'app.json').read_text())
    assets_dir = args.source/'assets/bip6'
    original = read_input(args.input)
    outer = dict((i.filename,b) for i,b in read_zip(original))
    if set(outer) != {'device.zip','app-side.zip'}:
        raise ValueError('Unexpected ZPK structure')
    device = dict((i.filename,b) for i,b in read_zip(outer['device.zip']))
    built_cfg = json.loads(device['app.json'])
    for field in ('appId','appName','version'):
        if built_cfg['app'][field] != source_cfg['app'][field]:
            raise ValueError('Build/source mismatch: '+field)
    source_exact = encode_source_assets(original,assets_dir)
    if args.fixed:
        result = fixed_optimize(source_exact)
        strategy = 'source-RGBA/original-order/inner-STORE/outer-best-Zopfli120-or-zlib9'
    else:
        result,strategy,_ = optimize(source_exact,baseline_label='source-exact-input')
    checked,after = verify_release(original,result,assets_dir)
    args.output.parent.mkdir(parents=True,exist_ok=True)
    args.output.write_bytes(result)
    print(json.dumps({
        'output':str(args.output),'bytes':len(result),'strategy':strategy,
        'images_matching_source_RGBA':checked,
        'headers_and_padding_unchanged':True,'non_image_payloads_unchanged':True,
        'device_expanded_bytes':sum(map(len,after['device.zip'].values())),
        'sha256':hashlib.sha256(result).hexdigest()
    },ensure_ascii=False),flush=True)


if __name__ == '__main__':
    main()
