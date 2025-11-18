#!/bin/sh

echo "--- Starting Entrypoint Script ---"

# Создаем все необходимые директории, используя переменные окружения
mkdir -p ${CONFIG_FILEPATH} ${ADULT_PROMPT_FILEPATH} ${SCHOOLCHILD_PROMPT_FILEPATH} ${ADULT_HTML_TEMPLATE_FILEPATH} ${SCHOOLCHILD_HTML_TEMPLATE_FILEPATH} ${PDF_FILEPATH} ${AI_RESPONSE_FILEPATH}

echo "--- Dirs created. Listing contents of VOLUME_PATH: ---"
ls -la ${VOLUME_PATH}

# Запускаем ваше Go-приложение
exec /bin/backend-service-profreport "$@"
