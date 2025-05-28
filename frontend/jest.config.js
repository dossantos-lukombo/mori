// jest.config.js
module.exports = {
  // 1. Simuler un navigateur ("window", DOM, etc.)
  testEnvironment: 'jsdom',
  testEnvironmentOptions: {
    customExportConditions: ['node', 'node-addons']
  },

  // 2. Où Jest va chercher vos tests et fichiers sources
  roots: ['<rootDir>/src'],

  // 3. Extensions à reconnaître
  moduleFileExtensions: ['js', 'json', 'vue'],

  // 4. Transformeurs : .vue via Vue3-Jest, .js via Babel
  transform: {
    '^.+\.vue$': '@vue/vue3-jest',
    '^.+\.js$': 'babel-jest'
  },

  // 5. Mapping des alias (pour @/…)
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/src/$1',
    '\.(css|less|sass|scss)$': 'identity-obj-proxy',
    '\.(png|jpg|jpeg|gif|svg|webp)$': '<rootDir>/src/__mocks__/fileMock.js'
  },

  // 6. Options supplémentaires
  globals: {
    'vue-jest': {
      transformAssetUrls: false,
      compilerOptions: {
        isCustomElement: tag => tag.startsWith('router-')
      }
    }
  },

  // 7. Rapport de couverture (optionnel)
  collectCoverage: true,
  coverageDirectory: 'coverage',
  coverageReporters: ['text', 'lcov']
}