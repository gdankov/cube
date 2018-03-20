set -euo pipefail
GOOS=linux go build -a -o recipe
docker build --build-arg buildpacks="$(< "buildpacks.json")" -t "dvrs/recipe:build" .
docker run -it -e 'APP_ID=19791b75-7cad-4b34-83f8-db636427d673' -e 'STAGING_GUID=staging-guid' --rm -v "$(pwd)/out:/out" dvrs/recipe:build
