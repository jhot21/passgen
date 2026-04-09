const KEY = 'passgen_settings';

const DEFAULTS = {
  salt: '',
  version: 3,
  length: 40,
  noUppers: false,
  noNumbers: false,
  noSpecials: false,
  customSpecials: null,  // null = use library default (!@#$%^&*)
  noAmbiguous: false,
  customAmbiguous: null, // null = use library default (Il1O05S)
};

/**
 * Load settings from localStorage, merging with defaults.
 * @returns {typeof DEFAULTS}
 */
export function loadSettings() {
  try {
    const raw = localStorage.getItem(KEY);
    if (!raw) return { ...DEFAULTS };
    return { ...DEFAULTS, ...JSON.parse(raw) };
  } catch {
    return { ...DEFAULTS };
  }
}

/**
 * Persist settings to localStorage.
 * @param {typeof DEFAULTS} settings
 */
export function saveSettings(settings) {
  try {
    localStorage.setItem(KEY, JSON.stringify(settings));
  } catch {
    // Ignore quota or privacy errors
  }
}
