import vapoursynth as vs
from vapoursynth import core
import math

import havsfunc

def Frameblend(clip, args, recipe) -> vs.VideoNode:
    bright_blend = recipe["frame_blending"]["bright_blend"]
    
    if bright_blend:
        og_format = clip.format
        try:
            og_matrix = clip.get_frame(0).props['_Matrix']
        except KeyError:
            og_matrix = 1
        clip = core.resize.Bicubic(clip=clip, format=vs.RGB48, transfer_in_s="709", transfer_s="linear", matrix_in_s="709")

    try:
        clip = core.frameblender.FrameBlend(clip, args["weighting"], True)
    except Exception as e:
        raise

    if bright_blend:
        clip = core.resize.Bicubic(clip=clip, format=og_format, matrix=og_matrix, transfer_s="709", transfer_in_s="linear", matrix_s="709")

    try:
        clip = havsfunc.ChangeFPS(
                clip,
                recipe["frame_blending"]["fps"]
                )
    except Exception as e:
        raise
    return clip
