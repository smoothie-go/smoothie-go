import os
# pyrefly: ignore [missing-import]
import vapoursynth as vs
from vapoursynth import core
import havsfunc

def apply_mask(clip: vs.VideoNode, original_clip: vs.VideoNode, recipe: dict) -> vs.VideoNode:
    if not recipe["artifact_masking"]["enabled"]:
        return clip

    folder_path = recipe["artifact_masking"]["folder_path"]
    file_name = recipe["artifact_masking"]["file_name"]
    mask_path = os.path.join(folder_path, file_name)

    if not os.path.exists(mask_path):
        raise FileNotFoundError(f"Mask file not found at: {mask_path}")

    mask = core.bs.VideoSource(source=mask_path, cachemode=0)
    if hasattr(core, "query_video_format"):
        mask_format = core.query_video_format(vs.GRAY, clip.format.sample_type, clip.format.bits_per_sample, 0, 0)
    else:
        mask_format = core.register_format(vs.GRAY, clip.format.sample_type, clip.format.bits_per_sample, 0, 0)
    if mask.format.color_family == vs.RGB:
        mask = core.resize.Bicubic(mask, width=clip.width, height=clip.height, format=mask_format.id, matrix_s="709")
    else:
        mask = core.resize.Bicubic(mask, width=clip.width, height=clip.height, format=mask_format.id)

    if recipe["artifact_masking"]["feathering"]:
        mask = mask.std.Minimum().std.BoxBlur(vradius=6, hradius=6, vpasses=2, hpasses=2)

    mask = mask * len(clip)

    if len(original_clip) != len(clip):
        original_clip = havsfunc.ChangeFPS(original_clip, fpsnum=clip.fps.numerator, fpsden=clip.fps.denominator)
        if len(original_clip) < len(clip):
            original_clip = original_clip + original_clip[-1] * (len(clip) - len(original_clip))
        elif len(original_clip) > len(clip):
            original_clip = original_clip[:len(clip)]

    return core.std.MaskedMerge(clipa=original_clip, clipb=clip, mask=mask, first_plane=True)
