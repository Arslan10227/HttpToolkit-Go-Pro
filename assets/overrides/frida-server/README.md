# Bundled Frida server binaries (Android)

Pre-built `frida-server` executables for offline Android interception. Filenames match the dynamic dependency cache format used at runtime:

```
frida-server-android-<arch>-<version>.bin
```

Populate this folder before packaging the server:

```bash
cd httptoolkit-node
npm run download:frida-servers
```

Or from the web UI project:

```bash
cd webui
npm run build:server
```

At runtime the server uses these bundled files first and only downloads if a matching binary is missing.
