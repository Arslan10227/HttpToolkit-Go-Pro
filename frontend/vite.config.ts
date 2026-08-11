import { defineConfig } from "vite";
import wails from "@wailsio/runtime/plugins/vite";
import { resolve, join } from "path";
import { copyFileSync, mkdirSync, readdirSync, existsSync, readFileSync, writeFileSync } from "fs";
import { fileURLToPath } from "url";

const __dirname = fileURLToPath(new URL(".", import.meta.url));

// Custom plugin to copy the pre-built UI assets and inject the shellApi script
function copyPrebuiltUI() {
  return {
    name: "copy-prebuilt-ui",
    closeBundle() {
      const assetsDir = resolve(__dirname, "../assets");
      const distDir = resolve(__dirname, "dist");

      // Copy index.html with shellApi.js injection
      const indexHtml = join(assetsDir, "index.html");
      if (existsSync(indexHtml)) {
        const content = readFileSync(indexHtml, "utf-8");
        const scriptTag = '<script type="module" src="/shellApi.js"></script>';
        const modified = content.replace(
          '<script type="module" crossorigin src="/assets/index-',
          scriptTag + '\n    <script type="module" crossorigin src="/assets/index-'
        );
        writeFileSync(join(distDir, "index.html"), modified);
      }

      // Copy the assets/assets/ directory (pre-built JS/CSS bundles)
      const builtAssetsDir = join(assetsDir, "assets");
      const destAssetsDir = join(distDir, "assets");
      if (existsSync(builtAssetsDir)) {
        mkdirSync(destAssetsDir, { recursive: true });
        for (const file of readdirSync(builtAssetsDir)) {
          copyFileSync(join(builtAssetsDir, file), join(destAssetsDir, file));
        }
      }

      // Copy root-level static files (favicon, logos, manifest)
      for (const file of ["favicon.png", "logo.png", "logo-icon.png", "manifest.json"]) {
        const src = join(assetsDir, file);
        if (existsSync(src)) {
          copyFileSync(src, join(distDir, file));
        }
      }
    },
  };
}

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    host: "127.0.0.1",
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
  plugins: [
    wails("./bindings"),
    copyPrebuiltUI(),
  ],
  build: {
    outDir: "dist",
    emptyOutDir: true,
    rollupOptions: {
      input: {
        shellApi: resolve(__dirname, "src/shellApi.ts"),
      },
      output: {
        entryFileNames: "shellApi.js",
        chunkFileNames: "[name].js",
        assetFileNames: "[name][extname]",
      },
    },
  },
});
