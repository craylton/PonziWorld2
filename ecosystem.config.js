module.exports = {
  apps: [
    {
      name: 'backend',
      script: 'main.go',
      cwd: './backend',
      interpreter: 'go',
      watch: true
    }
  ]
};
