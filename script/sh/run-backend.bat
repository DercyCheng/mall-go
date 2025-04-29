@echo off
echo ===== Mall-Go Backend Services =====

REM Check if infrastructure services are running
docker ps | findstr "mall-mysql" > nul
if %errorlevel% neq 0 (
    echo Starting infrastructure services...
    cd %~dp0script\docker
    docker-compose -f docker-compose-env.yml up -d
    cd %~dp0
    timeout /t 10
) else (
    echo Infrastructure services are already running.
)

REM Set working directory to the project root
cd %~dp0

REM Create log directories if they don't exist
if not exist .\logs mkdir .\logs

echo.
echo Choose which service to run:
echo 1. User Service
echo 2. Gateway Service
echo 3. All Services (in separate terminals)
echo 4. Stop infrastructure services
echo 5. Exit
echo.

set /p choice="Enter your choice (1-5): "

if "%choice%"=="1" (
    echo Starting User Service...
    cd services\user-service
    start "User Service" cmd /k "go run cmd\server\main.go"
) else if "%choice%"=="2" (
    echo Starting Gateway Service...
    cd services\gateway-service
    start "Gateway Service" cmd /k "go run cmd\server\main.go"
) else if "%choice%"=="3" (
    echo Starting All Services...
    cd services\user-service
    start "User Service" cmd /k "go run cmd\server\main.go"
    cd ..\..\services\gateway-service
    start "Gateway Service" cmd /k "go run cmd\server\main.go"
) else if "%choice%"=="4" (
    echo Stopping infrastructure services...
    cd %~dp0script\docker
    docker-compose -f docker-compose-env.yml down
) else if "%choice%"=="5" (
    echo Exiting...
    exit
) else (
    echo Invalid choice!
)

echo.
echo Done!