def yuv_heuristics(width: int, height: int):
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
