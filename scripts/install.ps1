#!/usr/bin/env pwsh
param(
  [string]$Version
)

$ErrorActionPreference = 'Stop'

$Owner = 'BasWilson'
$Repo  = 'supershell'

Write-Host 'Detecting OS/ARCH...'
$arch = (Get-CimInstance Win32_Processor).Architecture
switch ($arch) {
  9 { $GOARCH = 'amd64' }
  12 { $GOARCH = 'arm64' }
  default { throw "Unsupported ARCH: $arch" }
}

if (-not $Version) {
  Write-Host 'Querying latest release tag...'
  $resp = Invoke-RestMethod -Uri "https://api.github.com/repos/$Owner/$Repo/releases/latest"
  $Version = $resp.tag_name
}

$asset = "supershell_${Version}_windows_${GOARCH}.zip"
$url   = "https://github.com/$Owner/$Repo/releases/download/$Version/$asset"
$tmp   = New-Item -ItemType Directory -Path ([System.IO.Path]::GetTempPath() + [System.Guid]::NewGuid().ToString())

Write-Host "Downloading $asset..."
Invoke-WebRequest -Uri $url -OutFile (Join-Path $tmp.FullName $asset)

Write-Host 'Extracting...'
Add-Type -AssemblyName System.IO.Compression.FileSystem
[System.IO.Compression.ZipFile]::ExtractToDirectory((Join-Path $tmp.FullName $asset), $tmp.FullName)

$dest = "$env:ProgramFiles\\supershell"
New-Item -ItemType Directory -Force -Path $dest | Out-Null
Copy-Item (Join-Path $tmp.FullName 'supershell.exe') (Join-Path $dest 'supershell.exe') -Force

$bin = "$env:ProgramFiles\\supershell"
if ($env:Path -notlike "*${bin}*") {
  Write-Host 'Adding to PATH for current user...'
  $current = [Environment]::GetEnvironmentVariable('Path', 'User')
  if (-not $current) { $current = '' }
  [Environment]::SetEnvironmentVariable('Path', $current + ";$bin", 'User')
}

Write-Host "Installed supershell to $bin. Open a new terminal to use it."


