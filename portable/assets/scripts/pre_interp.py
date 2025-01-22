import vapoursynth as vs
from vapoursynth import core

import utils
import sys

def pre_interp(clip, args, recipe) -> vs.VideoNode:
    # mostly taken straight out of https://github.com/couleur-tweak-tips/smoothie-rs/blob/main/target/jamba.vpy
    # beacuse i have absolutely no idea about all this heuristic shit
    model_path: str = recipe["pre_interp"]["model"].strip('"')
    
    og_format = clip.format
    
    heuristic = utils.yuv_heuristic(clip.width, clip.height)

    not_in_heuristic = {}
    
    for key, value in heuristic.items():
        not_in_heuristic[key.replace("_in", "")] = value

    try:
        clip = core.resize.Bicubic(
                clip=clip,
                format=vs.RGBS,
                **heuristic
                )
    except Exception as e:
        raise
    
    factor = recipe["pre_interp"]["factor"].strip('x')

    # masking (later)
    try:
        clip = core.rife.RIFE(
                clip=clip,
                factor_num=factor,
                model_path=model_path,
                gpu_id=0,
                gpu_thread=1,
                tta=recipe["pre_interp"]["tta"],
                uhd=recipe["pre_interp"]["uhd"],
                sc=recipe["pre_interp"]["scene_change"],
                )
    except Exception as e:
        raise
        
    try:
        clip = core.resize.Bicubic(
               clip=clip,
               format=og_format,
                **not_in_heuristic
                )
    except Exception as e:
        raise

    return clip
