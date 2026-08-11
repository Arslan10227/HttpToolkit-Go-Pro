# Server override binaries

The HTTP Toolkit server expects bundled tools in this directory:

| File | Purpose |
|------|---------|
| `java-agent.jar` | JVM attach + `-javaagent` interception |

If `java-agent.jar` is missing, the **Attach JVM** interceptor will not be activable.

Obtain it from an [httptoolkit-server](https://github.com/httptoolkit/httptoolkit-server) release or a full HTTP Toolkit install, then copy the JAR here.

Run `npm run verify:server-overrides` from the repo root to check required files.
