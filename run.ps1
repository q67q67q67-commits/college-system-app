# Запуск API колледжа — добавляет Go в PATH и запускает сервер
$goPath = "C:\Program Files\Go\bin"
if (Test-Path $goPath) {
    $env:Path = "$goPath;" + $env:Path
}
Set-Location $PSScriptRoot
Write-Host "Запуск API..."
go run ./cmd/api
