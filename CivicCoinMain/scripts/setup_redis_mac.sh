#!/bin/bash

# Script para instalar y levantar Redis en macOS

# Verificar si Homebrew está instalado
if ! command -v brew &> /dev/null; then
    echo "Homebrew no está instalado. Instalando Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    echo "Homebrew instalado correctamente."
fi

# Actualizar Homebrew
echo "Actualizando Homebrew..."
brew update

# Instalar Redis usando Homebrew
echo "Instalando Redis..."
brew install redis

# Verificar si Redis se instaló correctamente
if ! command -v redis-server &> /dev/null; then
    echo "Error: Redis no se instaló correctamente."
    exit 1
fi

echo "Redis instalado correctamente."

# Iniciar Redis como un servicio
echo "Iniciando Redis como un servicio..."
brew services start redis

sleep 5

# Verificar que Redis está en ejecución
if redis-cli ping | grep -q "PONG"; then
    echo "Redis está en ejecución y respondiendo correctamente."
else
    echo "Error: Redis no está respondiendo. Verifica la configuración."
    exit 1
fi

echo "Redis ha sido instalado y está listo para ser usado."