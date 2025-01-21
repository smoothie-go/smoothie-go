import vapoursynth as vs
from vapoursynth import core
import math

import havsfunc

def Frameblend(clip, args, recipe) -> vs.VideoNode:
    clip = core.frameblender.FrameBlend(clip, args["weighting"], True)
    clip = havsfunc.ChangeFPS(
            clip,
            recipe["frame_blending"]["fps"]
            )
    return clip
