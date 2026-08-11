/**
 * shellApi.ts — Bridge between Wails v3 runtime bindings and the pre-built
 * HttpToolkit web UI.
 *
 * The pre-built UI (from the separate httptoolkit-web project) expects two
 * global objects to be available:
 *   - window.go.main.ShellApp  — methods like GetServerAuthToken, GetServerPort, etc.
 *   - window.runtime           — EventsOn, EventsOff, EventsEmit for Wails events
 *
 * In Wails v2 these were injected automatically by the runtime. In Wails v3,
 * bindings are imported as ES modules from @wailsio/runtime. This shim
 * re-exports them onto the window globals so the pre-built UI works unchanged.
 */

import { Events, Call } from "@wailsio/runtime";

// Import the generated ShellService bindings. The path is relative to the
// bindings output directory configured in vite.config.ts (./bindings).
// After running `wails3 generate bindings`, this file will exist at
// frontend/bindings/httptoolkitpro/shellservice.js
import * as shellService from "../bindings/github.com/Arslan10227/HttpToolkit-Go-Pro/shellservice";

// Expose ShellService methods on window.go.main.ShellApp for the pre-built UI
interface GoGlobal {
  [pkg: string]: any;
}
interface WindowGo {
  main?: {
    ShellApp?: any;
  };
  [key: string]: any;
}

declare global {
  interface Window {
    go?: WindowGo;
    runtime?: {
      EventsOn: (name: string, cb: (data: any) => void) => void;
      EventsOff: (name: string, cb?: (data: any) => void) => void;
      EventsEmit: (name: string, data?: any) => void;
    };
    httpToolkitAuthToken?: string;
    httpToolkitDesktopVersion?: string;
  }
}

// Initialize window.go.main.ShellApp with the service bindings
window.go = window.go || {};
window.go.main = window.go.main || {};
window.go.main.ShellApp = shellService;

// Initialize window.runtime with the Wails v3 Events API (compatibility shim)
window.runtime = {
  EventsOn: (name: string, cb: (data: any) => void) => {
    Events.On(name, (event: any) => {
      // Wails v3 passes data wrapped in an event object; v2 passed it directly.
      // Extract the data field if present, otherwise pass the whole event.
      const data = event?.data ?? event;
      cb(data);
    });
  },
  EventsOff: (name: string, cb?: (data: any) => void) => {
    Events.Off(name);
  },
  EventsEmit: (name: string, data?: any) => {
    Events.Emit(name, data);
  },
};

// Mark that we're running in a Wails desktop environment
try {
  localStorage.setItem("htk-is-desktop-wails", "true");
} catch {
  // localStorage may not be available in some contexts
}

// Log successful initialization
console.log("[shellApi] Wails v3 bindings bridge initialized");
