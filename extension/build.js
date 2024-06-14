const esbuild = require('esbuild');

esbuild.build({
  entryPoints: [
    { in: 'src/script.js', out: 'script' },
    { in: 'src/background.js', out: 'background' },
  ],
  outdir: 'dist',
  bundle: true,
});
