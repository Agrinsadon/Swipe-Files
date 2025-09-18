// formatBytes converts a byte count to a human-readable string.
export function formatBytes(n: number) {
  if (n < 1024) return `${n} B`;
  const units = ["KB", "MB", "GB", "TB"];
  let i = -1;
  do {
    n = n / 1024;
    i++;
  } while (n >= 1024 && i < units.length - 1);
  return `${n.toFixed(1)} ${units[i]}`;
}

// formatDate renders a user-friendly date from an ISO string.
export function formatDate(iso: string) {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return iso;
  return d.toLocaleDateString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

// isImageByNameOrExt returns true when a filename or extension matches common image types.
export function isImageByNameOrExt(name?: string, ext?: string) {
  const target = (ext || name || "").toLowerCase();
  return /\.(png|jpe?g|gif|webp|bmp|svg)$/.test(target);
}

// Additional helpers for media/text/office detection by extension.
export function isVideoByNameOrExt(name?: string, ext?: string) {
  const target = (ext || name || "").toLowerCase();
  return /\.(mp4|webm|ogv|ogg|mov|m4v)$/i.test(target);
}

export function isAudioByNameOrExt(name?: string, ext?: string) {
  const target = (ext || name || "").toLowerCase();
  return /\.(mp3|wav|ogg|oga|m4a|aac)$/i.test(target);
}

export function isTextByNameOrExt(name?: string, ext?: string) {
  const target = (ext || name || "").toLowerCase();
  return /\.(txt|md|markdown|json|log|csv|tsv|js|ts|jsx|tsx|py|go|rs|java|c|cpp|cs|sh|yml|yaml|ini|cfg|toml)$/i.test(target);
}

export function isOfficeByNameOrExt(name?: string, ext?: string) {
  const target = (ext || name || "").toLowerCase();
  return /\.(docx|doc|dotx|xlsx|xls|pptx|ppt|odt|odp|ods)$/i.test(target);
}
