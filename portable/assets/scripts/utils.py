import vapoursynth as vs
from vapoursynth import core

def yuv_heuristic(width: int, height: int) :
    result = {}

    if width >= 3840:
        result["matrix_in_s"] = "2020ncl"
    elif width >= 1280:
        result["matrix_in_s"] = "709"
    elif height == 576:
        result["matrix_in_s"] = "470bg"
    else:
        result["matrix_in_s"] = "170m"

    if width >= 3840:
        result["transfer_in_s"] = "st2084"
    elif width >= 1280:
        result["transfer_in_s"] = "709"
    elif height == 576:
        result["transfer_in_s"] = "470bg"
    else:
        result["transfer_in_s"] = "601"

    if width >= 3840:
        result["primaries_in_s"] = "2020"
    elif width >= 1280:
        result["primaries_in_s"] = "709"
    elif height == 576:
        result["primaries_in_s"] = "470bg"
    else:
        result["primaries_in_s"] = "170m"

    result["range_in_s"] = "limited"

    # ITU-T H.273 (07/2021), Note at the bottom of pg. 20
    if width >= 3840:
        result["chromaloc_in_s"] = "top_left"
    else:
        result["chromaloc_in_s"] = "left"

    return result

def ScaleLuminance (scale: bool, clip: vs.VideoNode):
    try:
        y = core.std.ShufflePlanes(clip, planes=0, colorfamily=vs.GRAY)
        u = core.std.ShufflePlanes(clip, planes=1, colorfamily=vs.GRAY)
        v = core.std.ShufflePlanes(clip, planes=2, colorfamily=vs.GRAY)
    except Exception:
        raise
    try:
        if scale: # up
            y = core.resize.Point(y, width=y.width * 2, height=y.height * 2)
        else: # down
            y = core.resize.Point(y, width=y.width / 2, height=y.height / 2)
    except Exception:
        raise

    try:
        clip = core.std.ShufflePlanes(clips=[y, u, v], planes=[0, 0, 0], colorfamily=vs.YUV)
    except Exception:
        raise

    return clip
