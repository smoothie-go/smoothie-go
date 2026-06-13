param (
    [string]$Tag = "latest",
    [string]$LocalBinary = "smoothie-go.exe",
    [string]$OutFile = "smoothie-go-setup.msi",
    [string]$Version = "1.0.0"
)

$ErrorActionPreference = "Stop"

$OutFile = [System.IO.Path]::GetFullPath($OutFile)
if (Test-Path $LocalBinary) {
    $LocalBinary = [System.IO.Path]::GetFullPath($LocalBinary)
}

if ($Version -match '^(\d+\.\d+\.\d+)') {
    $Version = $Matches[1]
} else {
    $Version = "1.0.0"
}

$DlDir = Join-Path $PSScriptRoot "build_downloads"
$TmpDir = Join-Path $PSScriptRoot "build_temp"
$LayoutDir = Join-Path $TmpDir "layout"

# Clean up any leftover directories from previous runs
if ($TmpDir -and (Test-Path $TmpDir) -and ($TmpDir -like "*build_temp*")) {
    Remove-Item -Path $TmpDir -Recurse -Force
}
if ($DlDir -and (Test-Path $DlDir) -and ($DlDir -like "*build_downloads*")) {
    Remove-Item -Path $DlDir -Recurse -Force
}

New-Item -ItemType Directory -Force -Path $DlDir, $TmpDir, $LayoutDir | Out-Null

try {
    if (-not (Test-Path $LocalBinary)) {
    $ReleaseUrl = "https://api.github.com/repos/smoothie-go/smoothie-go/releases/tags/$Tag"
    if ($Tag -eq "latest") {
        $ReleaseUrl = "https://api.github.com/repos/smoothie-go/smoothie-go/releases/latest"
    }
    $ReleaseJson = Invoke-RestMethod -Uri $ReleaseUrl -UseBasicParsing
    $Asset = $ReleaseJson.assets | Where-Object { $_.name -like "*windows*" -or $_.name -like "*win64*" } | Select-Object -First 1
    $SmgoUrl = $Asset.browser_download_url
    $SmgoZip = Join-Path $DlDir "smoothie-go.zip"
    Invoke-WebRequest -Uri $SmgoUrl -OutFile $SmgoZip -UseBasicParsing
    Expand-Archive -Path $SmgoZip -DestinationPath (Join-Path $TmpDir "smoothie-go") -Force
    $SmgoExe = Get-ChildItem -Path (Join-Path $TmpDir "smoothie-go") -File | Where-Object { $_.Name -like "smoothie-go*" } | Select-Object -First 1
    Copy-Item -Path $SmgoExe.FullName -Destination (Join-Path $LayoutDir "smoothie-go.exe") -Force
} else {
    Copy-Item -Path $LocalBinary -Destination (Join-Path $LayoutDir "smoothie-go.exe") -Force
}

$FfmpegUrl = "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-n7.1-latest-win64-gpl-7.1.zip"
$FfmpegZip = Join-Path $DlDir "ffmpeg.zip"
if (-not (Test-Path $FfmpegZip) -or (Get-Item $FfmpegZip).Length -lt 10MB) {
    Invoke-WebRequest -Uri $FfmpegUrl -OutFile $FfmpegZip -UseBasicParsing
}
Expand-Archive -Path $FfmpegZip -DestinationPath (Join-Path $TmpDir "ffmpeg") -Force
$Ffbin = Join-Path $TmpDir "ffmpeg\ffmpeg-n7.1-latest-win64-gpl-7.1\bin"
Copy-Item -Path (Join-Path $Ffbin "ffmpeg.exe") -Destination $LayoutDir -Force
Copy-Item -Path (Join-Path $Ffbin "ffplay.exe") -Destination $LayoutDir -Force
Copy-Item -Path (Join-Path $Ffbin "ffprobe.exe") -Destination $LayoutDir -Force

$VsUrl = "https://github.com/smoothie-go/VSBundler/releases/download/Nightly_2025.08.14_02-55/VapourSynth.zip"
$VsZip = Join-Path $DlDir "VapourSynth.zip"
if (-not (Test-Path $VsZip) -or (Get-Item $VsZip).Length -lt 10MB) {
    Invoke-WebRequest -Uri $VsUrl -OutFile $VsZip -UseBasicParsing
}
Expand-Archive -Path $VsZip -DestinationPath (Join-Path $TmpDir "vsbundle") -Force
Copy-Item -Path (Join-Path $TmpDir "vsbundle\VapourSynth\*") -Destination $LayoutDir -Recurse -Force

$WixDir = Join-Path $TmpDir "wix"
$Candle = "candle.exe"
$Light = "light.exe"
if (-not (Get-Command $Candle -ErrorAction SilentlyContinue)) {
    $WixZip = Join-Path $DlDir "wix.zip"
    if (-not (Test-Path $WixZip) -or (Get-Item $WixZip).Length -lt 10MB) {
        Invoke-WebRequest -Uri "https://github.com/wixtoolset/wix3/releases/download/wix3112rtm/wix311-binaries.zip" -OutFile $WixZip -UseBasicParsing
    }
    if (-not (Test-Path $WixDir)) {
        New-Item -ItemType Directory -Force -Path $WixDir | Out-Null
        Expand-Archive -Path $WixZip -DestinationPath $WixDir -Force
    }
    $Candle = Join-Path $WixDir "candle.exe"
    $Light = Join-Path $WixDir "light.exe"
}

$script:ComponentIds = [System.Collections.Generic.List[string]]::new()

function Get-WixDirectoryTree {
    param (
        [string]$Path,
        [string]$ParentId,
        [string]$Indent = "          "
    )

    $Files = Get-ChildItem -Path $Path -File
    $Dirs = Get-ChildItem -Path $Path -Directory
    $xml = ""

    foreach ($file in $Files) {
        $fileId = "f_" + [Guid]::NewGuid().Guid.Replace("-", "")
        $compId = "c_" + [Guid]::NewGuid().Guid.Replace("-", "")
        $filePath = $file.FullName
        $xml += "$Indent<Component Id=""$compId"" Guid=""$([Guid]::NewGuid())"">`n"
        $xml += "$Indent  <File Id=""$fileId"" Source=""$filePath"" KeyPath=""yes"" />`n"
        $xml += "$Indent</Component>`n"
        $script:ComponentIds.Add($compId)
    }

    foreach ($dir in $Dirs) {
        $dirId = "d_" + [Guid]::NewGuid().Guid.Replace("-", "")
        $xml += "$Indent<Directory Id=""$dirId"" Name=""$($dir.Name)"">`n"
        $xml += Get-WixDirectoryTree -Path $dir.FullName -ParentId $dirId -Indent ($Indent + "  ")
        $xml += "$Indent</Directory>`n"
    }

    return $xml
}

$EnvPathCompId = "comp_envpath"
$EnvPathCompGuid = [Guid]::NewGuid().Guid
$script:ComponentIds.Add($EnvPathCompId)

$WxsPath = Join-Path $TmpDir "smoothie-go.wxs"
$DirectoryTreeXml = Get-WixDirectoryTree -Path $LayoutDir -ParentId "INSTALLFOLDER" -Indent "          "

$ComponentRefsXml = ""
foreach ($id in $script:ComponentIds) {
    $ComponentRefsXml += "      <ComponentRef Id=""$id"" />`n"
}

$LicensePath = Join-Path $TmpDir "license.rtf"
$LicenseTextPath = Join-Path $PSScriptRoot "LICENSE"
if (Test-Path $LicenseTextPath) {
    $LicenseContent = Get-Content -Raw -Path $LicenseTextPath
    $EscapedContent = $LicenseContent.Replace('\', '\\').Replace('{', '\{').Replace('}', '\}')
    $RtfBody = $EscapedContent -replace "\r?\n", "\par`r`n"
    $RtfText = "{\rtf1\ansi\deff0{\fonttbl{\f0\fnil\fcharset0 Arial;}}\viewkind4\uc1\pard\lang1033\fs18 $RtfBody}"
    $RtfText | Out-File -FilePath $LicensePath -Encoding Ascii
} else {
    "{\rtf1\ansi\deff0{\fonttbl{\f0\fnil\fcharset0 Arial;}}\viewkind4\uc1\pard\lang1033\fs24\b smoothie-go Setup\b0\par\parPlease click Install to continue and install smoothie-go on your system.\par}" | Out-File -FilePath $LicensePath -Encoding Ascii
}

$WxsContent = @"
<?xml version="1.0" encoding="UTF-8"?>
<Wix xmlns="http://schemas.microsoft.com/wix/2006/wi">
  <Product Id="*" Name="smoothie-go" Language="1033" Version="$Version" Manufacturer="smoothie-go Team" UpgradeCode="8C46EFFA-15B8-4395-B62D-0F29B8AF654E">
    <Package InstallerVersion="200" Compressed="yes" InstallScope="perUser" />
    <MajorUpgrade AllowSameVersionUpgrades="yes" DowngradeErrorMessage="A newer version of [ProductName] is already installed." />
    <MediaTemplate EmbedCab="yes" />
    <UIRef Id="WixUI_Minimal" />
    <WixVariable Id="WixUILicenseRtf" Value="$LicensePath" />

    <Feature Id="ProductFeature" Title="smoothie-go" Level="1">
$ComponentRefsXml    </Feature>
  </Product>

  <Fragment>
    <Directory Id="TARGETDIR" Name="SourceDir">
      <Directory Id="LocalAppDataFolder">
        <Directory Id="INSTALLFOLDER" Name="smoothie-go">
          <Component Id="$EnvPathCompId" Guid="$EnvPathCompGuid">
            <Environment Id="PATH" Name="PATH" Value="[INSTALLFOLDER]" Action="set" Part="last" System="no" />
            <CreateFolder />
          </Component>
$DirectoryTreeXml        </Directory>
      </Directory>
    </Directory>
  </Fragment>
</Wix>
"@

    $WxsContent | Out-File -FilePath $WxsPath -Encoding utf8

    $WixUIExt = Join-Path $WixDir "WixUIExtension.dll"
    if (-not (Test-Path $WixUIExt)) {
        $WixUIExt = "WixUIExtension"
    }

    & $Candle $WxsPath -ext $WixUIExt -out (Join-Path $TmpDir "smoothie-go.wixobj")
    if ($LASTEXITCODE -ne 0) { throw "candle failed with exit code $LASTEXITCODE" }
    & $Light (Join-Path $TmpDir "smoothie-go.wixobj") -ext $WixUIExt -out $OutFile -sval
    if ($LASTEXITCODE -ne 0) { throw "light failed with exit code $LASTEXITCODE" }
}
finally {
    if ($TmpDir -and (Test-Path $TmpDir) -and ($TmpDir -like "*build_temp*")) {
        Remove-Item -Path $TmpDir -Recurse -Force
    }
    if ($DlDir -and (Test-Path $DlDir) -and ($DlDir -like "*build_downloads*")) {
        Remove-Item -Path $DlDir -Recurse -Force
    }
}
