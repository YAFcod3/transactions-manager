$envFile = ".env"

function Show-Help {
    Write-Host "Usage: .\compose.ps1 [up|down]" -ForegroundColor Yellow
    Write-Host
    Write-Host "Options:"
    Write-Host "  up      Build and start the environment"
    Write-Host "  down    Stop and remove all containers"
}

if ($Args.Count -eq 0) {
    Show-Help
    exit 1
}

switch ($Args[0]) {
    "up" {
        if (-Not (Test-Path $envFile)) {
            Write-Host "Error: $envFile not found." -ForegroundColor Red
            exit 1
        }
        Write-Host "Building and starting the environment..." -ForegroundColor Green
        docker-compose --env-file $envFile up --build -d
    }
    "down" {
        Write-Host "Stopping and removing all containers..." -ForegroundColor Yellow
        docker-compose down
    }

    Default {
        Write-Host "Invalid option: $($Args[0])" -ForegroundColor Red
        Show-Help
        exit 1
    }
}
