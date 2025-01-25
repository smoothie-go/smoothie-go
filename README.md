# Smoothie-GO
Why? for the lols. (Rewrite, of a rewrite, of a rewrite, of a rewrite)

The rewrites being, smoothie-go, smoothie-rs, smoothie, teres

![doing this](assets/doingthis.png)

## Checklist

### Internal

- [x] **config**
  - [x] Parsing INI
  - [x] Getting configs
  - [x] Presistent configs (Pretty much just, delete the file, it reappears on the next run)

- [x] **Interpolation**
  - [x] SVP
  - [x] RIFE
  - [x] Nvidia Optical Flow (Will be limited to linux users, unless you buy SVPFlow unfortunately)

- [x] **Args**
  - [x] Parsing args
  - [x] Validating args

- [ ] **VPYs**
  - [ ] Built in
    - [x] Interp
    - [x] Best source loading
    - [x] Pre-Interp
    - [x] Blending
    - [ ] Flowblur
    - [ ] Masking
  - [ ] ~~User provided vpy~~ *Not going to be implemented*

- [ ] **Server**
  - [ ] oauth
    - [ ] Discord
  - [ ] Auth
    - [ ] Basic access token
  - [ ] `smoothie-go server --authless`
  - [ ] /render
  - [ ] /config
  - [ ] /status

- [ ] **Web-UI**
  - [ ] Status page
  - [ ] Login page
  - [ ] Config page

### Recipe
- [x] **interpolation**
  - [x] interpolation::enabled
  - [x] interpolation::type (SVP/OF)
  - [ ] interpolation::masking
  - [x] interpolation::fps
  - [x] interpolation::speed
  - [x] interpolation::tuning
  - [x] interpolation::algorithm
  - [ ] interpolation::block_size
  - [x] interpolation::use_gpu
  - [ ] interpolation::area

- [ ] **frame_blending**
  - [x] frame_blending::enabled
  - [x] frame_blending::fps
  - [x] frame_blending::intensity
  - [x] frame_blending::weighting
  - [ ] frame_blending::bright_blend

- [ ] **flowblur**
  - [ ] flowblur::enabled
  - [ ] flowblur::masking
  - [ ] flowblur::amount
  - [ ] flowblur::do_blending

- [ ] **output**
  - [ ] output::process
  - [x] output::enc_args
  - [ ] output::file_format
  - [x] output::container

- [ ] **preview_window**
  - [x] preview_window::enabled
  - [ ] preview_window::process
  - [x] preview_window::output_args

- [ ] **artifact_masking**
  - [ ] artifact_masking::enabled
  - [ ] artifact_masking::feathering
  - [ ] artifact_masking::folder_path
  - [ ] artifact_masking::file_name

- [ ] **miscellaneous**
  - [ ] miscellaneous::play_ding
  - [x] miscellaneous::always_verbose
  - [x] miscellaneous::dedup_threshold
  - [ ] miscellaneous::global_output_folder
  - [ ] miscellaneous::source_indexing
  - [x] miscellaneous::ffmpeg_options
  - [x] miscellaneous::ffplay_options

- [ ] **console**
  - [ ] console::stay_on_top
  - [ ] console::borderless
  - [ ] console::position
  - [ ] console::width
  - [ ] console::height

- [ ] **timescale**
  - [ ] timescale::in
  - [ ] timescale::out

- [x] **color_grading**
  - [x] color_grading::enabled
  - [x] color_grading::brightness
  - [x] color_grading::saturation
  - [x] color_grading::contrast
  - [x] color_grading::hue
  - [x] color_grading::coring

- [ ] **lut**
  - [ ] lut::enabled
  - [ ] lut::path
  - [ ] lut::opacity

- [x] **pre_interp**
  - [x] pre_interp::enabled
  - [x] pre_interp::scene_change
  - [x] pre_interp::tta
  - [x] pre_interp::uhd
  - [x] pre_interp::masking
  - [x] pre_interp::factor
  - [x] pre_interp::model

## Priorities

* Portablity & stability over speed - I don't want it to be the fastest possible, I want it to work on most platforms and be as stable as possible, while being fast enough.

## Features that WONT be implemented (by me)

If you want any features from here, **implement it and PR**, you will most likely be accepted.

GUI, purely CLI for now.

`last_args.txt`, never saw the use in that.

`--rerun, -!!`, depends on last_args, just hit the up arrow, or use shell history

`--json`

`--tui`, GUI file picker

Frameserver, just pre-render bro


## Thanks
[couleur-tweak-tips/smoothie-rs](https://github.com/couleur-tweak-tips/smoothie-rs) - For the og implementation

