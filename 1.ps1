# generate-structure.ps1
# 在项目根目录运行：pwsh ./generate-structure.ps1

$ErrorActionPreference = "Stop"

# 目录 -> 文件列表 + 每个文件的 package
$tree = @{
    "internal/domain/audio" = @{
        "signal.go" = "audio"
        "asset.go"  = "audio"
        "errors.go" = "audio"
    }
    "internal/domain/effects" = @{
        "effect.go"    = "effects"
        "registry.go"  = "effects"
        "gain.go"      = "effects"
        "reverse.go"   = "effects"
        "normalize.go" = "effects"
    }
    "internal/domain/analysis" = @{
        "spectrum.go" = "analysis"
        "waveform.go" = "analysis"
        "stats.go"    = "analysis"
    }
    "internal/application" = @{
        "audio_service.go"    = "application"
        "effect_service.go"   = "application"
        "analysis_service.go" = "application"
        "clone_service.go"    = "application"
    }
    "internal/application/ports" = @{
        "audio_store.go" = "ports"
        "recorder.go"    = "ports"
        "player.go"      = "ports"
        "clone.go"       = "ports"
    }
    "internal/application/commands" = @{
        "audio.go"  = "commands"
        "effect.go" = "commands"
    }
    "internal/adapters/memory" = @{
        "audio_store.go" = "memory"
    }
    "internal/adapters/audioio" = @{
        "recorder.go" = "audioio"
        "player.go"   = "audioio"
        "file.go"     = "audioio"
    }
    "internal/adapters/clone" = @{
        "http_provider.go" = "clone"
    }
    "internal/transport/http" = @{
        "mapper.go" = "http"
        "errors.go" = "http"
    }
    # 空目录：值为空哈希
    "internal/transport/http/generated" = @{}
    "internal/transport/http/handlers"  = @{}
}

function Ensure-Dir([string]$path) {
    if (-not (Test-Path $path)) {
        New-Item -ItemType Directory -Path $path -Force | Out-Null
        Write-Host "DIR  $path"
    }
}

function Ensure-Gitkeep([string]$dir) {
    # 只有当目录下没有任何文件（不包括子目录里的 .gitkeep）时才放 .gitkeep
    $entries = Get-ChildItem -Path $dir -Force -File -ErrorAction SilentlyContinue
    if (-not $entries -or $entries.Count -eq 0) {
        $keep = Join-Path $dir ".gitkeep"
        if (-not (Test-Path $keep)) {
            New-Item -ItemType File -Path $keep -Force | Out-Null
            Write-Host "KEEP $keep"
        }
    }
}

foreach ($dir in $tree.Keys) {
    Ensure-Dir $dir

    foreach ($file in $tree[$dir].Keys) {
        $pkg = $tree[$dir][$file]
        $path = Join-Path $dir $file
        if (-not (Test-Path $path)) {
            Set-Content -Path $path -Value "package $pkg`r`n" -Encoding utf8
            Write-Host "FILE $path"
        } else {
            Write-Host "SKIP $path (exists)"
        }
    }

    Ensure-Gitkeep $dir
}

# 额外目录（可选）
$extra = @(
    "cmd/server",
    "api/v1",
    "api/openapi",
    "docs",
    "scripts",
    "tests",
    "web/dist"
)
foreach ($d in $extra) {
    Ensure-Dir $d
    Ensure-Gitkeep $d
}

Write-Host "`nDone." -ForegroundColor Green