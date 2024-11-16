#!/bin/bash

APP_NAME="transactions-manager-app"
envFile=".env"

GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m'

set -e

function show_help {
    echo -e "${YELLOW}Usage: $0 [up|down|clean]${NC}"
    echo
    echo "Options:"
    echo "  up      Build and start the environment"
    echo "  down    Stop and remove all containers"
    echo "  clean   Remove unused Docker resources"
}

if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}Error: docker-compose is not installed.${NC}"
    exit 1
fi

if [ $# -eq 0 ]; then
    show_help
    exit 1
fi

case $1 in
    up)
        echo -e "${GREEN}Building and starting the environment...${NC}"
        if [ ! -f "$envFile" ]; then
            echo -e "${RED}Error: $envFile not found.${NC}"
            exit 1
        fi
        docker-compose --env-file "$envFile" up --build -d
        ;;
    down)
        echo -e "${YELLOW}Stopping and removing all containers...${NC}"
        docker-compose down
        ;;
    clean)
        echo -e "${YELLOW}Cleaning up unused Docker resources...${NC}"
        docker system prune -f
        ;;
    *)
        echo -e "${RED}Invalid option: $1${NC}"
        show_help
        exit 1
        ;;
esac
