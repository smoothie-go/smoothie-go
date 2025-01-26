import vapoursynth as vs
from vapoursynth import core

import havsfunc
import utils
import json
import sys

def interp(clip: vs.VideoNode, args: dict, recipe: dict) -> vs.VideoNode:
    # SVP

    # multiple based interpolation, is to be handled in go within recipe/validate.go

    if clip.format.id in [vs.YUV444P8]:
        scaled = True
        try:
            clip = utils.ScaleLuminance(True, clip)
        except Exception as e:
            raise
    else:
        scaled = False

    if recipe["interpolation"]["type"] == "svp":
        try:
            clip = havsfunc.InterFrame(
                    Input=clip,
                    GPU=recipe["interpolation"]["use_gpu"],
                    Preset=recipe["interpolation"]["speed"],
                    Tuning=recipe["interpolation"]["tuning"],
                    NewNum=recipe["interpolation"]["fps"],
                    NewDen=1,
                    OverrideAlgo=recipe["interpolation"]["algorithm"],
                    )
        except Exception:
            raise
    elif recipe["interpolation"]["type"] == "of":
        # Tekno's code
        smooth_options = {
            "rate": {"num": recipe["interpolation"]["fps"], "abs": True},
            "algo": recipe["interpolation"]["algorithm"],
            "mask": {
                "area": 0,
                "area_sharp": 1.2,
            },
            "scene": {
                "blend": False,
                "mode": 0,
                "limits": {"blocks": recipe["interpolation"]["block_size"]},
            },
        }
        try:
            clip = core.svp2.SmoothFps_NVOF(
                clip, json.dumps(smooth_options), vec_src=clip, src=clip, fps=clip.fps
            )
        except Exception as e:
            raise

    if scaled:
        try:
            clip = utils.ScaleLuminance(False, clip)
        except Exception as e:
            raise

    return clip
