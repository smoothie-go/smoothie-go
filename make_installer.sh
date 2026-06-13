#!/bin/bash
set -e

TMP_DIR=$(mktemp -d -t smoothie-go-build.XXXXXX)

cleanup() {
	if [ -n "$TMP_DIR" ] && [ -d "$TMP_DIR" ] && [[ "$TMP_DIR" == */smoothie-go-build.* ]]; then
		find "$TMP_DIR" -delete
	fi
}
trap cleanup EXIT

DL_DIR="$TMP_DIR/downloads"
LAYOUT_DIR="$TMP_DIR/layout"

mkdir -p "$DL_DIR" "$LAYOUT_DIR"

GOOS=linux GOARCH=amd64 go build -o "$LAYOUT_DIR/smoothie-go-linux-amd64" .

curl -L -o "$DL_DIR/python.tar.gz" "https://github.com/smoothie-go/pyenv-build/releases/download/py3104-ubuntu2204/python-3.10.4.tar.gz"
curl -L -o "$DL_DIR/vapoursynth.tar.gz" "https://github.com/Stefan-Olt/vs-plugin-build/releases/download/vapoursynth-fa1f579a63b4fdef5f436c086875fe9c4879b5d2/vapoursynth-build-linux-x86_64.tar.gz"
curl -L -o "$DL_DIR/fmtconv.zip" "https://github.com/Stefan-Olt/vs-plugin-build/releases/download/vsplugin%2Ffmtconv%2Fgit-18a9cecb%2Flinux-glibc-x86_64%2F2024-10-10T14.44.42%2B00.00Z/fmtconv-git-18a9cecb-linux-glibc-x86_64.zip"
curl -L -o "$DL_DIR/bestsource.zip" "https://github.com/Stefan-Olt/vs-plugin-build/releases/download/vsplugin%2Fcom.vapoursynth.bestsource%2FR8%2Flinux-glibc-x86_64%2F2024-12-12T18.37.59%2B00.00Z/BestSource-R8-linux-glibc-x86_64.zip"
curl -L -o "$DL_DIR/libsvpflow1.so" "https://github.com/smoothie-go/smoothie-go/raw/refs/heads/master/resources/vapoursynth/libsvpflow1.so"
curl -L -o "$DL_DIR/libsvpflow2.so" "https://github.com/smoothie-go/smoothie-go/raw/refs/heads/master/resources/vapoursynth/libsvpflow2.so"
curl -L -o "$DL_DIR/frameblender.so" "https://github.com/couleurm/vs-frameblender/releases/download/1.2/vs-frameblender-1.2.so"
curl -L -o "$DL_DIR/librife.so" "https://github.com/styler00dollar/VapourSynth-RIFE-ncnn-Vulkan/releases/download/r9_mod_v32/librife_linux_x86-64.so"
curl -L -o "$DL_DIR/ffmpeg.tar.xz" "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-n7.1-latest-linux64-gpl-7.1.tar.xz"

mkdir -p "$TMP_DIR/python" && tar -xzf "$DL_DIR/python.tar.gz" -C "$TMP_DIR/python"
mkdir -p "$TMP_DIR/vapoursynth" && tar -xzf "$DL_DIR/vapoursynth.tar.gz" -C "$TMP_DIR/vapoursynth"
mkdir -p "$TMP_DIR/fmtconv" && unzip -o "$DL_DIR/fmtconv.zip" -d "$TMP_DIR/fmtconv"
mkdir -p "$TMP_DIR/bestsource" && unzip -o "$DL_DIR/bestsource.zip" -d "$TMP_DIR/bestsource"
mkdir -p "$TMP_DIR/ffmpeg" && tar -xJf "$DL_DIR/ffmpeg.tar.xz" -C "$TMP_DIR/ffmpeg"

mkdir -p "$LAYOUT_DIR/lib/vapoursynth"

cp -r "$TMP_DIR/python/3.10.4/bin/"* "$LAYOUT_DIR/"
cp -r "$TMP_DIR/python/3.10.4/lib/"* "$LAYOUT_DIR/lib/"
cp -r "$TMP_DIR/vapoursynth/workspace/lib/"* "$LAYOUT_DIR/lib/"
cp "$TMP_DIR/vapoursynth/workspace/bin/vspipe" "$LAYOUT_DIR/"

cp "$TMP_DIR/ffmpeg/ffmpeg-n7.1-latest-linux64-gpl-7.1/bin/ffmpeg" "$LAYOUT_DIR/"
cp "$TMP_DIR/ffmpeg/ffmpeg-n7.1-latest-linux64-gpl-7.1/bin/ffplay" "$LAYOUT_DIR/"
cp "$TMP_DIR/ffmpeg/ffmpeg-n7.1-latest-linux64-gpl-7.1/bin/ffprobe" "$LAYOUT_DIR/"

cp "$TMP_DIR/bestsource/bestsource.so" "$LAYOUT_DIR/lib/vapoursynth/"
cp "$TMP_DIR/fmtconv/libfmtconv.so" "$LAYOUT_DIR/lib/vapoursynth/"
cp "$DL_DIR/frameblender.so" "$LAYOUT_DIR/lib/vapoursynth/frameblender.so"
cp "$DL_DIR/libsvpflow1.so" "$LAYOUT_DIR/lib/vapoursynth/"
cp "$DL_DIR/libsvpflow2.so" "$LAYOUT_DIR/lib/vapoursynth/"
cp "$DL_DIR/librife.so" "$LAYOUT_DIR/lib/vapoursynth/"

tar -czf "$TMP_DIR/payload.tar.gz" -C "$LAYOUT_DIR" .

cat << 'EOF' > "$TMP_DIR/installer_header.sh"
#!/bin/bash
set -e
echo "Installing smoothie-go..."
INSTALL_DIR="$HOME/.local/share/smoothie-go"
mkdir -p "$INSTALL_DIR"
PAYLOAD_LINE=$(awk '/^__PAYLOAD_BELOW__/ {print NR + 1; exit 0;}' "$0")
tail -n +$PAYLOAD_LINE "$0" | tar -xz -C "$INSTALL_DIR"
echo "SystemPluginDir=$INSTALL_DIR/lib/vapoursynth" > "$INSTALL_DIR/vapoursynth.conf"
echo "Creating launcher in /usr/bin/smoothie-go (requires sudo)..."
sudo tee /usr/bin/smoothie-go > /dev/null << INNER_EOF
#!/bin/bash
DIR="$INSTALL_DIR"
exec env LD_LIBRARY_PATH="\$DIR/lib" PATH="\$DIR:\$PATH" PYTHONPATH="\$DIR/lib/python3.10/site-packages" PYTHONHOME="\$DIR" VAPOURSYNTH_CONF_PATH="\$DIR/vapoursynth.conf" "\$DIR/smoothie-go-linux-amd64" "\$@"
INNER_EOF
sudo chmod +x /usr/bin/smoothie-go
echo "smoothie-go has been successfully installed to $INSTALL_DIR"
echo "Launcher created at /usr/bin/smoothie-go"
exit 0
__PAYLOAD_BELOW__
EOF

cat "$TMP_DIR/installer_header.sh" "$TMP_DIR/payload.tar.gz" > smoothie-go-installer.run
chmod +x smoothie-go-installer.run
