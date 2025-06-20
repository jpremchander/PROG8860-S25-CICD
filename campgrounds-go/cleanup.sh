#!/bin/bash
echo "🧹 Cleaning up Next.js/React files..."

# Remove Next.js/React directories
rm -rf app/
rm -rf components/
rm -rf hooks/
rm -rf lib/
rm -rf styles/

# Remove Next.js/React config files
rm -f components.json
rm -f next.config.mjs
rm -f pnpm-lock.yaml
rm -f postcss.config.mjs
rm -f tailwind.config.ts
rm -f tsconfig.json

echo "✅ Cleanup complete! Now you have a pure Go project."
