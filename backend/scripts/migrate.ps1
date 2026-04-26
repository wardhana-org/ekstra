$envFile = Join-Path $PSScriptRoot "..\.env"

if (Test-Path $envFile) {
    Get-Content $envFile | ForEach-Object {
        if ($_ -match '^([^#][^=]+)=(.+)$') {
            Set-Item -Path "env:$($matches[1].Trim())" -Value $matches[2].Trim()
        }
    }
}

$databaseUrl = $env:DATABASE_URL
if (-not $databaseUrl) {
    $databaseUrl = $env:DB_URL
}

if (-not $databaseUrl) {
    Write-Error "DATABASE_URL or DB_URL is required"
    exit 1
}

$command = $args[0]
$value = $args[1]
$migrationPath = Join-Path $PSScriptRoot "..\migrations"

switch ($command) {
    "up" {
        migrate -path $migrationPath -database $databaseUrl up
    }
    "down" {
        $count = if ($value) { $value } else { "1" }
        Write-Host "Rolling back $count migration(s). Continue? [y/N]"
        $confirm = Read-Host
        if ($confirm -eq "y") {
            migrate -path $migrationPath -database $databaseUrl down $count
        }
    }
    "create" {
        if (-not $value) {
            Write-Error "Migration name is required"
            exit 1
        }
        migrate create -ext sql -dir $migrationPath -seq $value
    }
    "force" {
        if (-not $value) {
            Write-Error "Migration version is required"
            exit 1
        }
        migrate -path $migrationPath -database $databaseUrl force $value
    }
    default {
        Write-Host "Usage: .\scripts\migrate.ps1 up|down|create|force [value]"
    }
}
