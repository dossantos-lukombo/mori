const { defineConfig } = require('@vue/cli-service')

module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    port: 3000, // Serveur frontend sur le port 3000
    proxy: {
      '/api': {
        target: 'http://localhost:8080', // Serveur backend sur le port 8080
        changeOrigin: true,
        pathRewrite: { '^/api': '' },
      },
    },
  },
});
