$envFile = ".env"

function Show-Help {
    Write-Host "Usage: .\compose.ps1 [up|down|clean]" -ForegroundColor Yellow
    Write-Host
    Write-Host "Options:"
    Write-Host "  up      Build and start the environment"
    Write-Host "  down    Stop and remove all containers"
    Write-Host "  clean   Remove all unused Docker resources"
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
    "clean" {
        Write-Host "Cleaning up unused Docker resources..." -ForegroundColor Yellow
        docker system prune -f
    }
    Default {
        Write-Host "Invalid option: $($Args[0])" -ForegroundColor Red
        Show-Help
        exit 1
    }
}
