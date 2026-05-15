@echo off
echo ============================================
echo   Test API - Qwen3 LLM + Embedding
echo ============================================
echo.

echo [TEST 1] LLM - qwen3:0.6b ...
echo.
curl -s -X POST http://localhost:11434/api/chat ^
  -H "Content-Type: application/json" ^
  -d "{\"model\": \"qwen3.5:0.8b\", \"messages\": [{\"role\": \"user\", \"content\": \"Say hello in one sentence\"}], \"stream\": true, \"options\": {\"think\": false}}"
echo.
echo.

echo [TEST 2] Embedding - qwen3-embedding:0.6b ...
echo.
curl -s -X POST http://localhost:11434/api/embed ^
  -H "Content-Type: application/json" ^
  -d "{\"model\": \"qwen3-embedding:0.6b\", \"input\": \"Hello world\"}"
echo.
echo.

echo ============================================
echo   Neu thay JSON response la thanh cong!
echo ============================================
pause