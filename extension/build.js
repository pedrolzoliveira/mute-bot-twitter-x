const esbuild = require('esbuild');
const { copy } = require('esbuild-plugin-copy');

esbuild.build({
  entryPoints: [
    { in: 'src/script.js', out: 'script' },
    { in: 'src/background.js', out: 'background' },
  ],
  outdir: 'dist',
  bundle: true,
  plugins: [
    copy({
      assets: [
        {
          from: ['./manifest.json'],
          to: ['./manifest.json'],
        },
        {
          from: ['./icon.png'],
          to: ['./icon.png'],
        },
      ],
    })
  ],
});
