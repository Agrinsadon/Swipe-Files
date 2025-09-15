"use client";

import { FileInfo } from "@/lib/api";
import { isImageByNameOrExt } from "@/utils/format";

export function FilePreview({ file, previewUrl }: { file: FileInfo; previewUrl: string | null }) {
  const isImage = isImageByNameOrExt(file.name, file.ext);

  if (isImage && previewUrl) {
    return (
      <div className="media">
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img src={previewUrl} alt={file.name} className="img" />
      </div>
    );
  }

  return (
    <div className="media file-fallback" aria-label="Tiedosto">
      <div className="file-icon" aria-hidden>
        📄
      </div>
    </div>
  );
}

