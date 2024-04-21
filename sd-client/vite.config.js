import { defineConfig } from "vite";

export default defineConfig({
  base: "/static/",
  build: {
    target: "es2015",
    sourcemap: true,
  },
});
