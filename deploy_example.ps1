# === CONFIGURATION ===
$RemoteUser = "ubuntu"
$RemoteHost = "your-host-here"
$SSHKeyPath = "path\to\ssh\key"
$RemotePath = "/tmp/backend"

$BinaryName = "main"
$GoMainFile = "main.go"
$BuildDir = "build"
$BuildPath = "$BuildDir\backend"

function Write-Color {
    param (
        [string]$Text,
        [ConsoleColor]$Color = "White"
    )
    $origColor = $Host.UI.RawUI.ForegroundColor
    $Host.UI.RawUI.ForegroundColor = $Color
    Write-Host $Text
    $Host.UI.RawUI.ForegroundColor = $origColor
}

# Go to script directory (assumed to be inside the Go project)
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
Set-Location $ScriptDir

# Clean and prepare build directory
Write-Color "🧹  Cleaning build directory..." DarkYellow
Remove-Item -Recurse -Force $BuildPath

# Build Go binary for Linux
Write-Color "⚙️  Building Linux AMD64 Go binary..." Cyan
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o "$BuildPath\$BinaryName" $GoMainFile

if (!(Test-Path "$BuildPath\$BinaryName")) {
    Write-Color "❌  Build failed. Aborting." Red
    exit 1
}

# Copy necessary files into the build path
Copy-Item "config.json" "$BuildPath\" -Force
Copy-Item "words.txt" "$BuildPath\" -Force

# Upload folder to remote server
Write-Color "📤  Uploading backend folder to remote server..." Yellow
$scpCmd = "scp -i `"$SSHKeyPath`" -r `"$BuildPath\*`" $RemoteUser@${RemoteHost}:`"$RemotePath`""
Invoke-Expression $scpCmd

# Run remote restart script
Write-Color "🚀  Running remote restart script..." Yellow
$sshCmd = "ssh -i `"$SSHKeyPath`" $RemoteUser@$RemoteHost 'bash ~/restart_backend.sh'"
Invoke-Expression $sshCmd

Write-Color "✅  Backend deployment complete!" Green