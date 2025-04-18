# vi: set ft=python:
from pathlib import Path
import logging
import vapoursynth as vs
from vapoursynth import core
import json # to parse json recipe and args
import sys # logger.error
from sys import exit
import os 

script_dir = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, script_dir)

logger = logging.getLogger("smoothie")
console_handler = logging.StreamHandler()

formatter = logging.Formatter(
    "[{asctime}] [{levelname:<8}] {name}: {message}", 
    datefmt="%Y-%m-%d %H:%M:%S", 
    style="{"
)

console_handler.setFormatter(formatter)

logger.addHandler(console_handler)

for arg in ["args", "rec"]:
    if arg not in vars():
        logger.error(f"{arg} not defined")

try:
    arguments: dict = json.loads(str(args))
except json.JSONDecodeError as e:
    logger.error(f"Failed to parse args: {e}")
    exit(1)

file_handler = logging.FileHandler(arguments["log_file"], encoding="utf-8", mode="w")
file_handler.setFormatter(formatter)
logger.addHandler(file_handler)

try:
    recipe: dict = json.loads(str(rec))
except json.JSONDecodeError as e:
    logger.error(f"Failed to parse recipe: {e}")
    exit(1)

verbose: bool = arguments["verbose"]

if verbose:
    file_handler.setLevel(logging.DEBUG)
else:
    file_handler.setLevel(logging.WARN)
# load input video

arguments["out_dir"] = arguments["out_dir"].replace('\\','/')
arguments["input_file"] = arguments["input_file"].replace('\\','/')

cachepath: str = arguments["out_dir"] +"/"+ arguments["input_file"].split("/")[-1] + "-bsindex"

logger.debug(f"Cache path: {cachepath}")

global clip

try:
    clip = core.bs.VideoSource(
            source=arguments["input_file"],
            cachemode=3,
            cachepath=cachepath,
            showprogress=False
            )
    
    logger.debug(f"Loaded input video: {arguments['input_file']}")
except Exception as e:
    logger.error(f"Failed to load input video: {e}")
    exit(1)

if recipe["miscellaneous"]["dedup_threshold"] > 0.0:
    import filldrops
    
    clip = filldrops.FillDrops(
            clip,
            recipe["miscellaneous"]["dedup_threshold"]
            )
    
    logger.debug(f"Dedup threshold: {recipe['miscellaneous']['dedup_threshold']}")
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
    logger.debug(f"Color grading: {recipe['color_grading']}")
if recipe["pre_interp"]["enabled"]:
    import pre_interp # pre_interp (RIFE)
    logger.debug(f"RIFE model: {recipe['pre_interp']['model']}")
    try:
        clip = pre_interp.pre_interp(clip, arguments, recipe)
    except Exception as e:
        logger.error(f"Failed to apply pre-interpolation: {e}")
        exit(1)

if recipe["interpolation"]["enabled"]:
    import interp  # interp
    logger.debug(f"Interpolation: {recipe['interpolation']}")
    if clip.format.id in [vs.YUV444P8]:
        logger.warning("Slow interpolation due to YUV444P8 format")
    try:
        clip = interp.interp(clip, arguments, recipe)
    except Exception as e:
        logger.error(f"Failed to apply interpolation: {e}")
        exit(1)

    if abs(clip.fps.numerator / clip.fps.denominator - recipe["interpolation"]["fps"]) > 1e-6:
        logger.error(f"Interpolation Failed with unknown reason. Interpolation FPS: {clip.fps.numerator/clip.fps.denominator}, Recipe FPS: {recipe['interpolation']['fps']}")
        exit(1)
    logger.debug(f"Interpolation: successful")

if recipe["frame_blending"]["enabled"]:
    import frameblend
    logger.debug(f"Frame blending: {recipe['frame_blending']}")
    try:
        clip = frameblend.Frameblend(clip, arguments, recipe)
    except Exception as e:
        logger.error(f"Failed to apply frame blending: {e}")
        exit(1)
    logger.debug(f"Frame blending: successful")

if clip is None:
    logger.error("Clip is None. Ensure all processing steps succeeded.")
elif not isinstance(clip, vs.VideoNode):
    logger.error("Clip is not a valid VideoNode.")

clip.set_output()
