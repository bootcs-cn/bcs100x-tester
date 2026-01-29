#!/bin/bash
# Local Docker build script for bcs100x-tester
# This script handles the build context correctly

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
COURSES_ROOT="$(dirname "$PROJECT_ROOT")"

IMAGE_NAME="bootcs/bcs100x-tester"
TAG="${1:-local}"

echo "Building bcs100x-tester Docker image..."
echo "Image: ${IMAGE_NAME}:${TAG}"
echo ""

# Build from the courses directory to include both projects
cd "$COURSES_ROOT"

echo "Build context: $PWD"
echo "Dockerfile: bcs100x-tester/Dockerfile.multi-context"
echo ""

# Build the image
docker build \
  -f bcs100x-tester/Dockerfile.multi-context \
  -t "${IMAGE_NAME}:${TAG}" \
  .

echo ""
echo "✓ Build complete!"
echo ""
echo "Test the image:"
echo "  docker run --rm ${IMAGE_NAME}:${TAG} --help"
echo ""
echo "Run tests:"
echo "  docker run --rm -v \"\$(pwd)/your-project:/workspace:ro\" \\"
echo "    ${IMAGE_NAME}:${TAG} -s <stage> -d /workspace"
