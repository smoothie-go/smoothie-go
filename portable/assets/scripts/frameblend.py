import vapoursynth as vs
from vapoursynth import core
import math

import havsfunc

def Frameblend(clip, args, recipe) -> vs.VideoNode:
    try:
        clip = core.frameblender.FrameBlend(clip, args["weighting"], True)
    except Exception as e:
        raise

    try:
        clip = havsfunc.ChangeFPS(
                clip,
                recipe["frame_blending"]["fps"]
                )
    except Exception as e:
        raise
    return clip
