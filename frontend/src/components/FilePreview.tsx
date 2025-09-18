"use client";
// FilePreview: keskittää esikatselun (kuva/PDF/video/audio/teksti/toimisto) yhteen komponenttiin.
// Miksi: logiikka yhdessä paikassa; saa FileInfo-olion ja backendin tuottaman preview-URL:n.

import { FileInfo } from "@/lib/api";
import { useEffect, useMemo, useState } from "react";
import { isAudioByNameOrExt, isImageByNameOrExt, isOfficeByNameOrExt, isTextByNameOrExt, isVideoByNameOrExt } from "@/utils/format";

export function FilePreview({ file, previewUrl }: { file: FileInfo; previewUrl: string | null }) {
  const isImage = isImageByNameOrExt(file.name, file.ext);
  const isPdf = /\.pdf$/i.test((file.ext || file.name || "").toLowerCase());
  const isVideo = isVideoByNameOrExt(file.name, file.ext);
  const isAudio = isAudioByNameOrExt(file.name, file.ext);
  const isText = isTextByNameOrExt(file.name, file.ext);
  const isOffice = isOfficeByNameOrExt(file.name, file.ext);
  const [failed, setFailed] = useState(false);
  const wantsPreview = isImage || isPdf || isVideo || isAudio || isText || isOffice;

  const mime = useMemo(() => {
    const t = (file.ext || file.name || "").toLowerCase();
    if (/\.mp4$/.test(t) || /\.m4v$/.test(t)) return "video/mp4";
    if (/\.webm$/.test(t)) return "video/webm";
    if (/\.(ogv|ogg)$/.test(t)) return "video/ogg";
    if (/\.mov$/.test(t)) return "video/quicktime";
    if (/\.mp3$/.test(t)) return "audio/mpeg";
    if (/\.wav$/.test(t)) return "audio/wav";
    if (/\.(ogg|oga)$/.test(t)) return "audio/ogg";
    if (/\.m4a$/.test(t) || /\.aac$/.test(t)) return "audio/mp4";
    return undefined;
  }, [file.ext, file.name]);

  if (!failed && isImage && previewUrl) {
    return (
      <div className="media">
        {/* Kuva suoraan; ei Next Imagea tässä */}
        <img
          src={previewUrl}
          alt={file.name}
          className="img"
          onError={() => {
            console.warn("Image preview failed", { url: previewUrl, file });
            setFailed(true);
          }}
        />
      </div>
    );
  }

  if (!failed && isPdf && previewUrl) {
    return (
      <div className="media">
        {/* PDF upotettuna; selaimen oma renderöinti */}
        <embed
          src={previewUrl}
          type="application/pdf"
          className="doc"
          onError={() => {
            console.warn("PDF preview failed", { url: previewUrl, file });
            setFailed(true);
          }}
        />
      </div>
    );
  }

  if (!failed && isVideo && previewUrl) {
    return (
      <div className="media">
        <video className="media-box" controls preload="metadata" onError={() => setFailed(true)}>
          <source src={previewUrl} {...(mime ? { type: mime } : {})} />
        </video>
      </div>
    );
  }

  if (!failed && isAudio && previewUrl) {
    return (
      <div className="media">
        <audio className="audio" controls onError={() => setFailed(true)}>
          <source src={previewUrl} {...(mime ? { type: mime } : {})} />
        </audio>
      </div>
    );
  }

  if (!failed && isText && previewUrl) {
    return <TextPreview url={previewUrl} />;
  }

  // Näyttäisi esikatseltavalta, mutta URL puuttuu tai lataus kaatui.
  if (wantsPreview && (!previewUrl || failed)) {
    return (
      <div className="media file-fallback" aria-label="Ei esikatselua">
        <div className="file-icon" aria-hidden>📄</div>
        <div className="preview-note">
          {failed ? "Esikatselu epäonnistui (404?)" : "Esikatselu ei käytössä"}
        </div>
        {previewUrl ? (
          <a className="preview-link" href={previewUrl} target="_blank" rel="noreferrer">
            Avaa suoraan
          </a>
        ) : null}
      </div>
    );
  }

  // Toimistoasiakirjat: kokeile PDF-muunnosta (/api/convert)
  if (!failed && isOffice && previewUrl) {
    return (
      <div className="media">
        <embed
          src={previewUrl}
          type="application/pdf"
          className="doc"
          onError={() => {
            console.warn("Office->PDF preview failed", { url: previewUrl, file });
            setFailed(true);
          }}
        />
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

// TextPreview: hakee ~64KB tekstiä ja näyttää monospace-lohkona.
function TextPreview({ url }: { url: string }) {
  const [state, setState] = useState<{ text?: string; err?: string }>({});
  useEffect(() => {
    let alive = true;
    (async () => {
      try {
        const res = await fetch(url, { cache: "no-store" });
        if (!res.ok) throw new Error(`${res.status}`);
        const blob = await res.blob();
        // Rajoita kevyeksi (~64KB)
        const buf = await blob.slice(0, 64 * 1024).text();
        if (alive) setState({ text: buf });
      } catch (e: any) {
        if (alive) setState({ err: String(e) });
      }
    })();
    return () => {
      alive = false;
    };
  }, [url]);
  if (state.err) {
    return (
      <div className="media file-fallback">
        <div className="preview-note">Tekstiesikatselu epäonnistui</div>
        <a className="preview-link" href={url} target="_blank" rel="noreferrer">
          Avaa suoraan
        </a>
      </div>
    );
  }
  return (
    <div className="media">
      <pre className="pre">{state.text ?? "Ladataan…"}</pre>
    </div>
  );
}
