# vi: set ft=python:
import vapoursynth as vs
from vapoursynth import core


import json # to parse json recipe and args
import sys # print_error
import os 

script_dir = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, script_dir)

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
    recipe = json.loads(str(rec))
except json.JSONDecodeError as e:
    print_error(f"Failed to parse recipe: {e}")

verbose: bool = arguments["verbose"]

def verbose_print(msg: str):
    if verbose:
        print("VERB: " + msg, file=sys.stderr)

# load input video
cachepath: str = arguments["out_dir"] +"/"+ arguments["input_file"].split("/")[-1] + "-bsindex"

global clip

try:
    clip = core.bs.VideoSource(
            source=arguments["input_file"],
            cachemode=3,
            cachepath=cachepath,
            showprogress=False
            )
except Exception as e:
    print_error(f"Failed to load input video: {e}")

if recipe["miscellaneous"]["dedup_threshold"] > 0.0:
    import filldrops
    
    clip = filldrops.FillDrops(
            clip,
            recipe["miscellaneous"]["dedup_threshold"]
            )

if recipe["color_grading"]["enabled"]:
    import adjust # colour grading
    clip = adjust.Tweak(
            clip,
            hue=recipe["color_grading"]["hue"],
            sat=recipe["color_grading"]["saturation"],
            bright=recipe["color_grading"]["brightness"],
            cont=recipe["color_grading"]["contrast"],
            coring=recipe["color_grading"]["coring"]
            )
if recipe["pre_interp"]["enabled"]:
    import pre_interp # pre_interp (RIFE)
    verbose_print(f"RIFE model: {recipe['pre_interp']['model']}")
    clip = pre_interp.pre_interp(clip, arguments, recipe)
if recipe["interpolation"]["enabled"]:
    import interp # interp

if recipe["frame_blending"]["enabled"]:
    import frameblend
    clip = frameblend.Frameblend(clip, arguments, recipe)


print(clip, file=sys.stderr)

if clip is None:
    print_error("Clip is None. Ensure all processing steps succeeded.")
elif not isinstance(clip, vs.VideoNode):
    print_error("Clip is not a valid VideoNode.")


clip.set_output()
