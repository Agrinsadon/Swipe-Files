"use client";

import { CSSProperties } from "react";
import { FileInfo } from "@/lib/api";
import { FilePreview } from "@/components/FilePreview";
import { formatBytes, formatDate } from "@/utils/format";

type SwipeDir = "left" | "right" | null;

type Props = {
  file: FileInfo;
  style: CSSProperties;
  dragging: boolean;
  trashOpacity: number;
  keepOpacity: number;
  onPointerDown: (e: React.PointerEvent) => void;
  onPointerMove: (e: React.PointerEvent) => void;
  onPointerUp: (e: React.PointerEvent) => void;
  onPointerCancel: (e: React.PointerEvent) => void;
  anim: SwipeDir;
  previewUrl: string | null;
};

export function Card({
  file,
  style,
  dragging,
  trashOpacity,
  keepOpacity,
  onPointerDown,
  onPointerMove,
  onPointerUp,
  onPointerCancel,
  anim,
  previewUrl,
}: Props) {
  return (
    <div
      className={`card ${dragging ? "dragging" : ""}`}
      style={style}
      onPointerDown={onPointerDown}
      onPointerMove={onPointerMove}
      onPointerUp={onPointerUp}
      onPointerCancel={onPointerCancel}
      role="group"
      aria-roledescription="swipeable card"
      data-anim={anim ?? undefined}
    >
      <div className="badge badge-left" style={{ opacity: trashOpacity }}>
        TRASH
      </div>
      <div className="badge badge-right" style={{ opacity: keepOpacity }}>
        KEEP
      </div>

      <FilePreview file={file} previewUrl={previewUrl} />

      <div className="meta">
        <div className="name" title={file.name}>
          {file.name}
        </div>
        <div className="sub">
          {formatBytes(file.size)} • Muokattu {formatDate(file.modTime)}
        </div>
      </div>
    </div>
  );
}

