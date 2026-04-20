@echo off
echo ============================================
echo   Ollama Docker - Qwen3 LLM + Embedding
echo ============================================
echo.

:: Check Docker is running
docker info >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [!] Docker chua chay. Hay mo Docker Desktop truoc.
    pause
    exit /b 1
)
echo [OK] Docker dang chay.
echo.

:: Start containers
echo [*] Khoi dong Ollama container...
docker compose up -d ollama

echo.
echo [*] Cho Ollama san sang (healthcheck)...
:WAIT_LOOP
timeout /t 3 /nobreak >nul
docker inspect --format="{{.State.Health.Status}}" ollama 2>nul | findstr "healthy" >nul
if %ERRORLEVEL% NEQ 0 goto WAIT_LOOP
echo [OK] Ollama container da san sang.
echo.

:: Pull models (ollama-init container)
echo [*] Bat dau pull models (co the mat vai phut lan dau)...
docker compose up ollama-init
echo.

echo ============================================
echo   DONE! API san sang tai:
echo   http://localhost:11434
echo.
echo   LLM   : POST /api/chat
echo           { "model": "qwen3.5:0.8b", ... }
echo.
echo   Embed : POST /api/embed
echo           { "model": "qwen3-embedding:0.6b", ... }
echo.
echo   De test: chay test.bat
echo   De dung: docker compose down
echo ============================================
pause