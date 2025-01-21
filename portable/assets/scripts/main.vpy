# vi: set ft=python:
from tabnanny import verbose
from portable.assets.scripts.filldrops import FillDrops
import vapoursynth as vs
from vapoursynth import core


import json # to parse json recipe and args
import sys # print_error

def print_error(msg: str):
    print(msg, file=sys.stderr)
    exit(1)

for arg in ["args", "rec"]:
    if arg not in vars():
        print_error(f"{arg} not defined")

try:
    arguments = json.loads(str(args))
except json.JSONDecodeError as e:
    print_error(f"Failed to parse args: {e}")

try:
    recipe = json.loads(str(recipe))
except json.JSONDecodeError as e:
    print_error(f"Failed to parse recipe: {e}")

verbose: bool = arguments["verbose"]

def verbose_print(msg: str):
    if verbose:
        print(msg)

# load input video
cachepath: str = arguments["out_dir"] + arguments["input_video"].split("/")[-1] + "-bsindex"

try:
    clip: vs.VideoNode = core.bs.VideoSource(
            source=arguments["input_video"],
            cachemode=3,
            cachepath=cachepath,
            showprogress=False
            )
except Exception as e:
    print_error(f"Failed to load input video: {e}")

if recipe["miscellaneous"]["dedup_threshold"] > 0.0:
    import filldrops
    
    filldrops.FillDrops(
            clip,
            recipe["miscellaneous"]["dedup_threshold"]
            )
if recipe["frame_blending"]["enabled"]:
    import blend # frameblending using vs-frameblender

if recipe["interpolation"]["enabled"]:
    import interp # interp

if recipe["pre_interpolation"]["enabled"]:
    import pre_interp # pre_interp (RIFE)

if recipe["color_grading"]["enabled"]:
    import adjust # colour grading

