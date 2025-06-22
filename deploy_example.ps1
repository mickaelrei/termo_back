# === CONFIGURATION ===
$GoMainFile = "main.go"
$buildPath = "build"
$BinaryName = "backend"
$RemoteUser = "ubuntu"
$RemoteHost = "your-host-here"
$RemotePath = "/tmp/backend"
$SSHKeyPath = "path\to\ssh\key"

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

if (!(Test-Path $buildPath)) {
    New-Item -ItemType Directory -Path $buildPath | Out-Null
} else {
    Get-ChildItem -Path $buildPath -Recurse -Force |
        ForEach-Object {
            try {
                Remove-Item -Path $_.FullName -Force -Recurse -ErrorAction Stop
            } catch {}
        }
}

# Build Go binary for Linux
Write-Color "⚙️  Building Linux AMD64 Go binary..." Cyan
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o "$buildPath\$BinaryName" $GoMainFile

if (!(Test-Path "$buildPath\$BinaryName")) {
    Write-Color "❌  Build failed. Aborting." Red
    exit 1
}

# Upload binary to server
Write-Color "📤  Uploading binary to remote server..." Yellow
$scpCmd = "scp -i `"$SSHKeyPath`" `"$buildPath\$BinaryName`" $RemoteUser@${RemoteHost}:`"$RemotePath`""
Invoke-Expression $scpCmd

# Run remote restart script
Write-Color "🚀  Running remote restart script..." Yellow
$sshCmd = "ssh -i `"$SSHKeyPath`" $RemoteUser@$RemoteHost 'bash ~/restart_backend.sh'"
Invoke-Expression $sshCmd

Write-Color "✅  Backend deployment complete!" Green
