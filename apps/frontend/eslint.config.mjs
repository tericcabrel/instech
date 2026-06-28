import tsParser from '@typescript-eslint/parser';
import perfectionist from 'eslint-plugin-perfectionist';
import sortDestructureKeys from 'eslint-plugin-sort-destructure-keys';
import globals from 'globals';

export default [
  {
    ignores: [
      '**/dist',
      '**/*.css',
      '**/yarn.lock',
      '**/.husky',
      '**/package.json',
      '**/.yarn',
      '**/biome.jsonc',
      '**/node_modules',
      '**/components.json',
      'src/api/generated/**',
      'eslint.config.mjs',
    ],
  },
  {
    files: ['**/*.js', '**/*.ts', '**/*.tsx'],
    languageOptions: {
      globals: {
        ...globals.commonjs,
        ...globals.node,
      },
      parser: tsParser,
    },
    plugins: {
      perfectionist,
      'sort-destructure-keys': sortDestructureKeys,
    },
    rules: {
      'func-style': [
        'error',
        'expression',
        {
          overrides: {
            namedExports: 'expression',
          },
        },
      ],
      'newline-after-var': 'warn',
      'newline-before-return': 'warn',
      'perfectionist/sort-array-includes': 'warn',
      'perfectionist/sort-exports': 'warn',
      'perfectionist/sort-imports': 'off',
      'perfectionist/sort-interfaces': 'warn',
      'perfectionist/sort-intersection-types': 'warn',
      'perfectionist/sort-jsx-props': 'warn',
      'perfectionist/sort-named-exports': 'warn',
      'perfectionist/sort-object-types': 'warn',
      'perfectionist/sort-objects': 'warn',
      'perfectionist/sort-switch-case': 'warn',
      'perfectionist/sort-union-types': 'off',
      'sort-destructure-keys/sort-destructure-keys': 'error',
    },
    settings: {
      perfectionist: {
        type: 'natural',
      },
    },
  },
];
