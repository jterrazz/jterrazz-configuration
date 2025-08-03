"use strict";

module.exports = {
  config: {
    // Updates
    updateChannel: "stable",

    // Typography
    fontSize: 12,
    fontFamily:
      'Menlo, "DejaVu Sans Mono", Consolas, "Lucida Console", monospace',
    fontWeight: "normal",
    fontWeightBold: "bold",
    lineHeight: 1,
    letterSpacing: 0,

    // Cursor
    cursorColor: "rgba(248,28,229,0.8)",
    cursorAccentColor: "#000",
    cursorShape: "BLOCK",
    cursorBlink: false,

    // Colors
    foregroundColor: "#fff",
    backgroundColor: "#000",
    selectionColor: "rgba(248,28,229,0.3)",
    borderColor: "#333",

    // Styling
    css: "",
    termCSS: "",
    padding: "12px 14px",

    // Window
    workingDirectory: "",
    showHamburgerMenu: "",
    showWindowControls: "",

    // Color palette
    colors: {
      black: "#000000",
      red: "#C51E14",
      green: "#1DC121",
      yellow: "#C7C329",
      blue: "#0A2FC4",
      magenta: "#C839C5",
      cyan: "#20C5C6",
      white: "#C7C7C7",
      lightBlack: "#686868",
      lightRed: "#FD6F6B",
      lightGreen: "#67F86F",
      lightYellow: "#FFFA72",
      lightBlue: "#6A76FB",
      lightMagenta: "#FD7CFC",
      lightCyan: "#68FDFE",
      lightWhite: "#FFFFFF",
      limeGreen: "#32CD32",
      lightCoral: "#F08080",
    },

    // Shell
    shell: "",
    shellArgs: ["--login"],
    env: {},

    // Behavior
    bell: "SOUND",
    copyOnSelect: false,
    defaultSSHApp: true,
    quickEdit: false,
    preserveCWD: true,

    // Rendering
    webGLRenderer: true,
    disableLigatures: true,
    macOptionSelectionMode: "vertical",
    webLinksActivationKey: "",

    // Features
    disableAutoUpdates: false,
    screenReaderMode: false,
  },

  // Extensions
  plugins: [],
  localPlugins: [],
  keymaps: {},
};
//# sourceMappingURL=config-default.js.map
